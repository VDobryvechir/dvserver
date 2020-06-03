// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvcsv"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func (info *CollectorInfo) Start() error {
	db, err := dvdbdata.GetDBConnection(info.ConnectionName)
	if err != nil {
		return err
	}
	info.db = db
	return nil
}

func (info *CollectorInfo) Finish() error {
	if info != nil && info.db != nil {
		db := info.db
		info.db = nil
		return db.Close(false)
	}
	return nil
}

func (info *CollectorInfo) SaveCsvFile() error {
	data := dvevaluation.ConvertInterfaceListsMapToStringListsMap(info.DataCollector, dvevaluation.ConversionOptionSimpleLike)
	return dvcsv.WriteCsvToFile(info.CsvFileName, data, info.CsvWriteOptions)
}
