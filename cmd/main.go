package main

import (
	"context"
	"das-account-indexer/block_parser"
	"das-account-indexer/config"
	"das-account-indexer/dao"
	"das-account-indexer/http_server"
	"das-account-indexer/http_server/handle"
	"das-account-indexer/prometheus"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/go-redis/redis"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/scorpiotzh/mylog"
	"github.com/scorpiotzh/toolib"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
	"time"
)

var (
	log               = mylog.NewLogger("main", mylog.LevelDebug)
	exit              = make(chan struct{})
	ctxServer, cancel = context.WithCancel(context.Background())
	wgServer          = sync.WaitGroup{}
)

func main() {
	log.Debug("Start service: ")
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Usage:   "Server Type, ``(default): api and timer server, `api`: api server, `timer`: timer server",
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx *cli.Context) error {
	// init config file
	configFilePath := ctx.String("config")
	if err := config.InitCfg(configFilePath); err != nil {
		return err
	}

	// config file watcher
	watcher, err := config.AddCfgFileWatcher(configFilePath)
	if err != nil {
		return err
	}
	// ============= start service =============

	// prometheus
	prometheus.Init()
	prometheus.Tools.Run()

	// db
	dbDao, err := dao.NewGormDB(config.Cfg.DB.Mysql)
	if err != nil {
		return fmt.Errorf("dao.NewGormDB err: %s", err.Error())
	}
	log.Info("db ok")

	// cache
	red, err := toolib.NewRedisClient(config.Cfg.Cache.Redis.Addr, config.Cfg.Cache.Redis.Password, config.Cfg.Cache.Redis.DbNum)
	if err != nil {
		log.Error("NewRedisClient err:", err.Error())
	}

	// ckb node
	ckbClient, err := rpc.DialWithIndexer(config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
	if err != nil {
		return fmt.Errorf("rpc.DialWithIndexer err: %s", err.Error())
	}
	log.Info("ckb node ok")

	// das core
	env := core.InitEnvOpt(config.Cfg.Server.Net, common.DasContractNameConfigCellType, common.DasContractNameAccountCellType,
		common.DasContractNameBalanceCellType, common.DasContractNameDispatchCellType,
		common.DasContractNameReverseRecordCellType, common.DASContractNameSubAccountCellType, common.DasContractNameReverseRecordRootCellType,
		common.DasContractNameDidCellType, common.DasContractNameAlwaysSuccess)
	ops := []core.DasCoreOption{
		core.WithClient(ckbClient),
		core.WithDasContractArgs(env.ContractArgs),
		core.WithDasContractCodeHash(env.ContractCodeHash),
		core.WithDasNetType(config.Cfg.Server.Net),
		core.WithTHQCodeHash(env.THQCodeHash),
		core.WithDasRedis(red),
	}
	dasCore := core.NewDasCore(ctxServer, &wgServer, ops...)
	dasCore.InitDasContract(env.MapContract)
	if err := dasCore.InitDasConfigCell(); err != nil {
		return fmt.Errorf("InitDasConfigCell err: %s", err.Error())
	}
	if err := dasCore.InitDasSoScript(); err != nil {
		return fmt.Errorf("InitDasSoScript err: %s", err.Error())
	}
	dasCore.RunAsyncDasContract(time.Minute * 5)   // contract outpoint
	dasCore.RunAsyncDasConfigCell(time.Minute * 2) // config cell outpoint
	dasCore.RunAsyncDasSoScript(time.Minute * 5)   // so

	dasCore.RunSetConfigCellByCache([]core.CacheConfigCellKey{
		core.CacheConfigCellKeyCharSet,
		core.CacheConfigCellKeyReservedAccounts,
	})
	log.Info("das contract ok")

	// tx builder
	txBuilderBase := txbuilder.NewDasTxBuilderBase(ctxServer, dasCore, nil, "")

	//service mode
	mode := ctx.String("mode")

	if mode == "api" {
		if err := initApiServer(txBuilderBase, dasCore, dbDao, red); err != nil {
			return fmt.Errorf("initApiServer err : %s", err.Error())
		}
	} else if mode == "timer" {
		if err := initTimer(dasCore, dbDao); err != nil {
			return fmt.Errorf("initTimer err : %s", err.Error())
		}
	} else {
		if err := initTimer(dasCore, dbDao); err != nil {
			return fmt.Errorf("initTimer err : %s", err.Error())
		}
		if err := initApiServer(txBuilderBase, dasCore, dbDao, red); err != nil {
			return fmt.Errorf("initApiServer err : %s", err.Error())
		}
	}

	// ============= end service =============
	// exit
	toolib.ExitMonitoring(func(sig os.Signal) {
		log.Warn("ExitMonitoring:", sig.String())
		if watcher != nil {
			log.Warn("close watcher ... ")
			_ = watcher.Close()
		}
		cancel()
		//hs.Shutdown()
		wgServer.Wait()
		exit <- struct{}{}
	})

	<-exit
	log.Warn("success exit server. bye bye!")
	return nil
}

func initTimer(dasCore *core.DasCore, dbDao *dao.DbDao) error {

	// block parser
	bp := block_parser.BlockParser{
		DasCore:            dasCore,
		CurrentBlockNumber: config.Cfg.Chain.CurrentBlockNumber,
		DbDao:              dbDao,
		ConcurrencyNum:     config.Cfg.Chain.ConcurrencyNum,
		ConfirmNum:         config.Cfg.Chain.ConfirmNum,
		Ctx:                ctxServer,
		Cancel:             cancel,
		Wg:                 &wgServer,
	}
	if err := bp.RunParser(); err != nil {
		return fmt.Errorf("RunParser err: %s", err.Error())
	}
	log.Info("block parser ok")
	return nil
}

func initApiServer(txBuilderBase *txbuilder.DasTxBuilderBase, dasCore *core.DasCore, dbDao *dao.DbDao, red *redis.Client) error {
	builderConfigCell, err := dasCore.ConfigCellDataBuilderByTypeArgsList(
		common.ConfigCellTypeArgsPreservedAccount00,
		common.ConfigCellTypeArgsPreservedAccount01,
		common.ConfigCellTypeArgsPreservedAccount02,
		common.ConfigCellTypeArgsPreservedAccount03,
		common.ConfigCellTypeArgsPreservedAccount04,
		common.ConfigCellTypeArgsPreservedAccount05,
		common.ConfigCellTypeArgsPreservedAccount06,
		common.ConfigCellTypeArgsPreservedAccount07,
		common.ConfigCellTypeArgsPreservedAccount08,
		common.ConfigCellTypeArgsPreservedAccount09,
		common.ConfigCellTypeArgsPreservedAccount10,
		common.ConfigCellTypeArgsPreservedAccount11,
		common.ConfigCellTypeArgsPreservedAccount12,
		common.ConfigCellTypeArgsPreservedAccount13,
		common.ConfigCellTypeArgsPreservedAccount14,
		common.ConfigCellTypeArgsPreservedAccount15,
		common.ConfigCellTypeArgsPreservedAccount16,
		common.ConfigCellTypeArgsPreservedAccount17,
		common.ConfigCellTypeArgsPreservedAccount18,
		common.ConfigCellTypeArgsPreservedAccount19,
		common.ConfigCellTypeArgsUnavailable,
	)
	var mapReservedAccounts = make(map[string]struct{})
	var mapUnAvailableAccounts = make(map[string]struct{})
	if err != nil {
		var cacheBuilder core.CacheConfigCellReservedAccounts
		strCache, errCache := dasCore.GetConfigCellByCache(core.CacheConfigCellKeyReservedAccounts)
		if errCache != nil {
			log.Error("GetConfigCellByCache err: %s", errCache.Error())
			return fmt.Errorf("ConfigCellDataBuilderByTypeArgsList1 err: %s", err.Error())
		} else if strCache == "" {
			return fmt.Errorf("ConfigCellDataBuilderByTypeArgsList2 err: %s", err.Error())
		} else if errCache = json.Unmarshal([]byte(strCache), &cacheBuilder); errCache != nil {
			log.Error("json.Unmarshal err: %s", errCache.Error())
			return fmt.Errorf("ConfigCellDataBuilderByTypeArgsList3 err: %s", err.Error())
		}
		mapReservedAccounts = cacheBuilder.MapReservedAccounts
		mapUnAvailableAccounts = cacheBuilder.MapUnAvailableAccounts
	} else {
		mapReservedAccounts = builderConfigCell.ConfigCellPreservedAccountMap
		mapUnAvailableAccounts = builderConfigCell.ConfigCellUnavailableAccountMap
	}
	// http server
	hs := &http_server.HttpServer{
		Ctx: ctxServer,
		//Address:        config.Cfg.Server.HttpServerAddr,
		AddressIndexer: config.Cfg.Server.HttpServerAddrIndexer,
		//AddressReverse: config.Cfg.Server.HttpServerAddrReverse,
		H: &handle.HttpHandle{
			Ctx:                    ctxServer,
			Red:                    red,
			DbDao:                  dbDao,
			DasCore:                dasCore,
			TxBuilderBase:          txBuilderBase,
			MapReservedAccounts:    mapReservedAccounts,
			MapUnAvailableAccounts: mapUnAvailableAccounts,
		},
	}
	hs.Run()
	log.Info("http server ok")
	return nil
}
