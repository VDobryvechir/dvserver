// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
)

func (info *CollectorInfo) CollectAuxTables() error {
	tables := dvparser.ConvertToNonEmptyList(info.AuxTables)
	for _, tableId := range tables {
		ids := dvparser.GetKeysFromStringIntMap(info.IdCollector[tableId])
		if len(ids) == 0 {
			if logDbImportLevel >= dvlog.LogTrace {
				log.Printf("Not using %s because of no ids", tableId)
			}
			continue
		}
		if logDbImportLevel >= dvlog.LogDebug {
			log.Printf("Getting sql table for table %s ids: %v", tableId, ids)
		}
		res, err := dvdbdata.GetSqlTableByIds(info.db, tableId, ids)
		if err != nil {
			if logDbImportLevel >= dvlog.LogInfo {
				log.Printf("Failed to get sql table %s for ids %v because of %v", tableId, ids, err)
			}
			return err
		}
		info.AddDataBySpecificId(tableId, res)
	}
	return nil
}

func (info *CollectorInfo) AddDataBySpecificId(tableId string, data [][]interface{}) {
	info.DataCollector[tableId] = append(info.DataCollector[tableId], data...)
}
