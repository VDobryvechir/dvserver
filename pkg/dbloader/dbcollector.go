// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strings"
)

type CollectorInfo struct {
	ConnectionName  string `json:"connection"`
	CsvFileName     string `json:"file_name"`
	CsvWriteOptions int    `json:"csv_options"`
	AuxTables       string `json:"aux_tables"`
	BaseTables      string `json:"base_tables"`
	BasePrefix      string `json:"base_prefix"`
	IdCollector     map[string]map[string]int
	DataCollector   map[string][][]interface{}
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
	collectorInfo.IdCollector = make(map[string]map[string]int)
	collectorInfo.DataCollector = make(map[string][][]interface{})
	if collectorInfo.BasePrefix == "" {
		collectorInfo.BasePrefix = "BASE"
	} else {
		collectorInfo.BasePrefix = strings.ToUpper(collectorInfo.BasePrefix)
	}
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
		if logDbImportLevel >= dvlog.LogInfo {
			log.Printf("Failed to start dbimport %v", err)
		}
		return err
	}
	if logDbImportLevel >= dvlog.LogDetail {
		log.Printf("Connection opened, starting to collect info by base tables")
	}
	err = collectorInfo.CollectBaseTables()
	if err != nil {
		collectorInfo.Finish()
		return err
	}
	if logDbImportLevel >= dvlog.LogDetail {
		log.Printf("Base tables finished, starting to collect info by aux tables")
	}
	err = collectorInfo.CollectAuxTables()
	if err != nil {
		collectorInfo.Finish()
		return err
	}
	if logDbImportLevel >= dvlog.LogDetail {
		log.Printf("Aux tables finished, starting to save to csv file")
	}
	err = collectorInfo.SaveCsvFile()
	if err != nil {
		collectorInfo.Finish()
		return err
	}
	if logDbImportLevel >= dvlog.LogDetail {
		log.Printf("Csv saving finished, starting to clean up")
	}
	return collectorInfo.Finish()
}
