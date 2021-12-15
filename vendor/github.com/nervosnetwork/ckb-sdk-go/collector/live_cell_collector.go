package collector

import (
	"context"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/pkg/errors"
)

type ChangeOutputIndex struct {
	Value int
}

type LiveCellCollectResult struct {
	LiveCells []*indexer.LiveCell
	Capacity  uint64
	Options   map[string]interface{}
}

type CellCollectionIterator interface {
	HasNext() bool
	Next() error
	CurrentItem() (*indexer.LiveCell, error)
	Iterator() (CellCollectionIterator, error)
}

type LiveCellCollector struct {
	Client      rpc.Client
	SearchKey   *indexer.SearchKey
	SearchOrder indexer.SearchOrder
	Limit       uint64
	LastCursor  string
	EmptyData   bool
	TypeScript  *types.Script
	result      []*indexer.LiveCell
	itemIndex   int
}

func (c *LiveCellCollector) HasNext() bool {
	return c.itemIndex < len(c.result)
}

func (c *LiveCellCollector) Next() error {
	c.itemIndex++
	if c.itemIndex >= len(c.result) && c.LastCursor != "" {
		c.itemIndex = 0
		var err error
		c.result, c.LastCursor, err = c.collect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *LiveCellCollector) CurrentItem() (*indexer.LiveCell, error) {
	if c.itemIndex >= len(c.result) {
		return nil, errors.New("no such element")
	}
	return c.result[c.itemIndex], nil
}

func (c *LiveCellCollector) Iterator() (CellCollectionIterator, error) {
	result, lastCursor, err := c.collect()
	if err != nil {
		return nil, err
	}
	c.result = result
	c.LastCursor = lastCursor

	return c, nil
}

func (c *LiveCellCollector) collect() ([]*indexer.LiveCell, string, error) {
	if c.SearchKey == nil {
		return nil, "", errors.New("missing SearchKey error")
	}
	if c.SearchOrder != indexer.SearchOrderAsc && c.SearchOrder != indexer.SearchOrderDesc {
		return nil, "", errors.New("missing SearchOrder error")
	}
	var result []*indexer.LiveCell
	liveCells, err := c.Client.GetCells(context.Background(), c.SearchKey, c.SearchOrder, c.Limit, c.LastCursor)
	if err != nil {
		return nil, "", err
	}
	for _, cell := range liveCells.Objects {
		if c.TypeScript != nil {
			if !c.TypeScript.Equals(cell.Output.Type) {
				continue
			}
		} else {
			if cell.Output.Type != nil {
				continue
			}
		}
		if c.EmptyData && len(cell.OutputData) > 0 {
			continue
		}
		result = append(result, cell)
	}
	return result, liveCells.LastCursor, nil
}

func NewLiveCellCollector(client rpc.Client, searchKey *indexer.SearchKey, searchOrder indexer.SearchOrder, limit uint64, afterCursor string) *LiveCellCollector {
	return &LiveCellCollector{
		Client:      client,
		SearchKey:   searchKey,
		SearchOrder: searchOrder,
		Limit:       limit,
		LastCursor:  afterCursor,
	}
}
