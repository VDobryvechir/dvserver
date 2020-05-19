// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
)

func (info *CollectorInfo) CollectAuxTables() error {
	tables := dvparser.ConvertToNonEmptyList(info.AuxTables)
	for _, tableId := range tables {
		ids := info.IdCollector[tableId]
		if len(ids) == 0 {
			continue
		}
		res, err := dvdbdata.GetSqlTableByIds(info.db, tableId, ids)
		if err != nil {
			return err
		}
		info.AddDataBySpecificId(tableId, res)
	}
	return nil
}

func (info *CollectorInfo) AddDataBySpecificId(tableId string, data [][]string) {
	info.DataCollector[tableId] = append(info.DataCollector[tableId], data...)
}
