package main

import (
	"context"
	"das-account-indexer/block_parser"
	"das-account-indexer/config"
	"das-account-indexer/dao"
	"das-account-indexer/http_server"
	"das-account-indexer/http_server/handle"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/txbuilder"
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
		common.DasContractNameReverseRecordCellType, common.DASContractNameSubAccountCellType)
	ops := []core.DasCoreOption{
		core.WithClient(ckbClient),
		core.WithDasContractArgs(env.ContractArgs),
		core.WithDasContractCodeHash(env.ContractCodeHash),
		core.WithDasNetType(config.Cfg.Server.Net),
		core.WithTHQCodeHash(env.THQCodeHash),
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

	log.Info("das contract ok")

	// tx builder
	txBuilderBase := txbuilder.NewDasTxBuilderBase(ctxServer, dasCore, nil, "")

	// block parser
	bp := block_parser.BlockParser{
		DasCore:            dasCore,
		CurrentBlockNumber: config.Cfg.Chain.CurrentBlockNumber,
		DbDao:              dbDao,
		ConcurrencyNum:     config.Cfg.Chain.ConcurrencyNum,
		ConfirmNum:         config.Cfg.Chain.ConfirmNum,
		Ctx:                ctxServer,
		Wg:                 &wgServer,
	}
	if err := bp.RunParser(); err != nil {
		return fmt.Errorf("RunParser err: %s", err.Error())
	}
	log.Info("block parser ok")

	// http server
	hs := &http_server.HttpServer{
		Ctx:            ctxServer,
		Address:        config.Cfg.Server.HttpServerAddr,
		AddressIndexer: config.Cfg.Server.HttpServerAddrIndexer,
		AddressReverse: config.Cfg.Server.HttpServerAddrReverse,
		H: &handle.HttpHandle{
			Ctx:           ctxServer,
			Red:           red,
			DbDao:         dbDao,
			DasCore:       dasCore,
			TxBuilderBase: txBuilderBase,
		},
	}
	hs.Run()
	log.Info("http server ok")

	// ============= end service =============
	// exit
	toolib.ExitMonitoring(func(sig os.Signal) {
		log.Warn("ExitMonitoring:", sig.String())
		if watcher != nil {
			log.Warn("close watcher ... ")
			_ = watcher.Close()
		}
		cancel()
		hs.Shutdown()
		wgServer.Wait()
		exit <- struct{}{}
	})

	<-exit
	log.Warn("success exit server. bye bye!")
	return nil
}
