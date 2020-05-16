// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
)

func (info *CollectorInfo) CollectBaseTables() error {
	prefix := dbImport + "BASE_"
	tables := dvparser.ConvertToNonEmptyList(info.BaseTables)
	for _, tableId := range tables {
		ids := info.IdCollector[tableId]
		var res [][]string
		var err error
		query := dvparser.GlobalProperties[prefix+tableId]
		if query != "" {
			res, err = dvdbdata.GetSqlTableByQuery(db, ids, query)
		} else {
			if len(ids) == 0 {
				continue
			}
			res, err = dvdbdata.GetSqlTableByIds(db, tableId, ids)
		}
		if err != nil {
			return err
		}
		childInfo, err := dvdbdata.CollectAllChildInfo(tableId, res)
		if err!=nil {
			return err
		}
		info.AddIds(childInfo)
		info.AddDataBySpecificId(tableId, res)
	}
	return nil
}

func (info *CollectorInfo) AddIds(idMap map[string][]string) {
	if ids != nil && info != nil {
		for tableId, ids := range idMap {
			info.IdCollector[tableId] = append(info.IdCollector[tableId], ids...)
		}
	}
}
