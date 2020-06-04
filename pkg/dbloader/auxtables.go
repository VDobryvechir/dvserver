// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
)

func filterIdsByOtherIds(ids []string, filter map[string]int) []string {
	n := len(ids)
	m := 0
	for i := 0; i < n; i++ {
		id := ids[i]
		if _, ok := filter[id]; !ok {
			ids[m] = id
			m++
		}
	}
	return ids[:m]
}

func (info *CollectorInfo) CollectAuxTables() error {
	tables := dvparser.ConvertToNonEmptyList(info.AuxTables)
	prefix := dbImport + info.BasePrefix + "_FILTER_"
	for _, tableId := range tables {
		ids := dvparser.GetKeysFromStringIntMap(info.IdCollector[tableId])
		filter := dvparser.GlobalProperties[prefix+tableId]
		if filter != "" {
			ids = filterIdsByOtherIds(ids, info.IdCollector[filter])
		}
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
