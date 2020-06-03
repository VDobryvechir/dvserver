// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
)

func (info *CollectorInfo) CollectBaseTables() error {
	prefix := dbImport + "BASE_"
	tables := dvparser.ConvertToNonEmptyList(info.BaseTables)
	for _, tableId := range tables {
		ids := dvparser.GetKeysFromStringIntMap(info.IdCollector[tableId])
		var res [][]interface{}
		var err error
		query := dvparser.GlobalProperties[prefix+tableId]
		if query == "" {
			if len(ids) == 0 {
				if logDbImportLevel >= dvlog.LogInfo {
					log.Printf("Query for %s omitted because neither query nor ids", tableId)
				}
				continue
			}
			res, err = dvdbdata.GetSqlTableByIds(info.db, tableId, ids)
		} else {
			_, _, startPos, endPos := dvdbdata.FindIdsPlaceholder(query)
			if startPos < 0 || endPos-startPos < 2 {
				if logDbImportLevel >= dvlog.LogInfo {
					log.Printf("Getting data by query %s", query)
				}
				res, err = dvdbdata.GetSqlTableByQuery(info.db, ids, query)
			} else {
				idKind := query[startPos : endPos-1]
				ids = dvparser.GetKeysFromStringIntMap(info.IdCollector[idKind])
				if len(ids) == 0 {
					log.Printf("Warning: no ids for %s", idKind)
					continue
				}
				if logDbImportLevel >= dvlog.LogInfo {
					log.Printf("Getting data by query %s with %s ids: %v", query, idKind, ids)
				}
				res, err = dvdbdata.GetSqlTableByQuery(info.db, ids, query)
			}
		}
		if err != nil {
			if logDbImportLevel >= dvlog.LogInfo {
				log.Printf("Failed to get data for %s: %v", tableId, err)
			}
			return err
		}
		err = dvdbdata.CollectAllChildInfo(tableId, res, info.IdCollector)
		if err != nil {
			if logDbImportLevel >= dvlog.LogInfo {
				log.Printf("Failed to collect child info for %s: %v", tableId, err)
			}
			return err
		}
		info.AddDataBySpecificId(tableId, res)
	}
	return nil
}
