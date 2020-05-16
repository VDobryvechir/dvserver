// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvcsv"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
)

func (info *CollectorInfo) Start() error {
	db, kind, err := dvdbdata.GetDB(info.ConnectionName)
	if err != nil {
		return err
	}
	info.db = db
	info.dbKind = kind
	return nil
}

func (info *CollectorInfo) Finish() error {
	if info != nil && info.db != nil {
		db := info.db
		info.db = nil
		return dvdbdata.CloseConnection(db)
	}
	return nil
}

func (info *CollectorInfo) SaveCsvFile() error {
	return dvcsv.WriteCsvToFile(info.CsvFileName, info.DataCollector, info.CsvWriteOptions)
}
