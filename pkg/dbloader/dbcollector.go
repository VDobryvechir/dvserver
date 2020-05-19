// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

type CollectorInfo struct {
	ConnectionName  string `json:"connection"`
	CsvFileName     string `json:"file_name"`
	CsvWriteOptions int    `json:"csv_options"`
	AuxTables       string `json:"aux_tables"`
	BaseTables      string `json:"base_tables"`
	IdCollector     map[string][]string
	DataCollector   map[string][][]string
	db              *dvdbdata.DBConnection
}

const (
	dbImport = "DB_IMPORT_"
)

func validateCollectorInfo(collectorInfo *CollectorInfo) error {
	if collectorInfo.ConnectionName == "" {
		return errors.New("empty connection for the import")
	}
	if collectorInfo.CsvFileName == "" {
		return errors.New("empty csvFileName for the import")
	}
	return nil
}

func initCollectorInfo(importName string) (*CollectorInfo, error) {
	s := dbImport + importName
	info := strings.TrimSpace(dvparser.GlobalProperties[s])
	if info == "" {
		return nil, errors.New("Empty property " + s)
	}
	collectorInfo := &CollectorInfo{}
	err := json.Unmarshal([]byte(info), collectorInfo)
	if err != nil {
		return nil, err
	}
	err = validateCollectorInfo(collectorInfo)
	if err != nil {
		return nil, err
	}
	collectorInfo.IdCollector = make(map[string][]string)
	collectorInfo.DataCollector = make(map[string][][]string)
	return collectorInfo, nil
}

func LoadRelatedTablesByImportName(importName string) error {
	collectorInfo, err := initCollectorInfo(importName)
	if err != nil {
		return err
	}
	return LoadRelatedTables(collectorInfo)
}

func LoadRelatedTables(collectorInfo *CollectorInfo) error {
	err := collectorInfo.Start()
	if err != nil {
		return err
	}
	err = collectorInfo.CollectBaseTables()
	if err != nil {
		collectorInfo.Finish()
		return err
	}
	err = collectorInfo.CollectAuxTables()
	if err != nil {
		collectorInfo.Finish()
		return err
	}
	err = collectorInfo.SaveCsvFile()
	if err != nil {
		collectorInfo.Finish()
		return err
	}
	return collectorInfo.Finish()
}
