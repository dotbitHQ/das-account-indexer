package dascache

import (
	"context"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"sync"
	"time"
)

var log = mylog.NewLogger("cache", mylog.LevelDebug)

type DasCache struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	rw          sync.RWMutex
	mapOutPoint map[string]int64
}

func NewDasCache(ctx context.Context, wg *sync.WaitGroup) *DasCache {
	return &DasCache{
		ctx:         ctx,
		wg:          wg,
		rw:          sync.RWMutex{},
		mapOutPoint: make(map[string]int64),
	}
}

func (d *DasCache) AddOutPoint(outPoint []string) {
	if len(outPoint) == 0 {
		return
	}
	d.rw.Lock()
	defer d.rw.Unlock()
	for _, v := range outPoint {
		d.mapOutPoint[v] = time.Now().Unix()
	}
}

func (d *DasCache) clearExpiredOutPoint(t time.Duration) {
	d.rw.Lock()
	defer d.rw.Unlock()
	timestamp := time.Now().Add(-t).Unix()
	log.Info("clearExpiredOutPoint before:", len(d.mapOutPoint))
	for k, v := range d.mapOutPoint {
		if v < timestamp {
			delete(d.mapOutPoint, k)
			//log.Info("clearExpiredOutPoint:", k, time.Now().String())
		}
	}
	log.Info("clearExpiredOutPoint after:", len(d.mapOutPoint))
}

func (d *DasCache) ExistOutPoint(outPoint string) bool {
	d.rw.RLock()
	defer d.rw.RUnlock()
	if _, ok := d.mapOutPoint[outPoint]; ok {
		return true
	}
	return false
}

func (d *DasCache) RunClearExpiredOutPoint(t time.Duration) {
	ticker := time.NewTicker(t)
	d.wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				d.clearExpiredOutPoint(t)
			case <-d.ctx.Done():
				d.wg.Done()
				return
			}
		}
	}()
}

func (d *DasCache) AddCellInputByAction(action common.DasAction, inputs []*types.CellInput) {
	var outPoints []string
	switch action {
	case common.DasActionStartAccountSale:
		for i := 1; i < len(inputs); i++ {
			outPoints = append(outPoints, common.OutPointStruct2String(inputs[i].PreviousOutput))
		}
	case common.DasActionBuyAccount:
		for i := 2; i < len(inputs); i++ {
			outPoints = append(outPoints, common.OutPointStruct2String(inputs[i].PreviousOutput))
		}
	case common.DasActionWithdrawFromWallet:
		for i := 0; i < len(inputs); i++ {
			outPoints = append(outPoints, common.OutPointStruct2String(inputs[i].PreviousOutput))
		}
	default:
		for i := 0; i < len(inputs); i++ {
			outPoints = append(outPoints, common.OutPointStruct2String(inputs[i].PreviousOutput))
		}
	}
	d.AddOutPoint(outPoints)
}
