// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvoc"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strings"
)

var logDbImportLevel = dvlog.LogError

func ProvideDbImportCommand() {
	dvoc.AddProcessFunction("dbimport", dvoc.ProcessFunction{
		Init: processDbImportInit,
		Run:  processDbImportRun,
	})
}

func Init() bool {
	ProvideDbImportCommand()
	return true
}

var inited = Init()

func processDbImportInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	command = strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if command == "" || !dvparser.IsUpperAlphaDigital(command) {
		log.Printf("Invalid parameters of dbimport execution command: provide the import name in upper letters")
		return nil, false
	}
	if dvparser.GlobalProperties[dbImport+command] == "" {
		log.Printf("Property %s is not defined for import", dbImport+command)
		return nil, false
	}
	return []interface{}{command}, true
}

func processDbImportRun(data []interface{}) bool {
	InitByGlobalProperties()
	name := data[0].(string)
	if logDbImportLevel >= dvlog.LogInfo {
		log.Printf("Started dbimport %s", name)
	}
	err := LoadRelatedTablesByImportName(name)
	if err != nil {
		if logDbImportLevel >= dvlog.LogError {
			log.Printf("Failed to import %s: %v", name, err)
		}
		return false
	} else if logDbImportLevel >= dvlog.LogInfo {
		log.Printf("Finish import %s successfully", name)
	}
	return true
}

func InitByGlobalProperties() {
	dvdbdata.InitByGlobalProperties()
	logDbImportLevel = dvlog.GetLogLevelByDefinition(dvparser.GlobalProperties["DVLOG_DBIMPORT_LEVEL"], logDbImportLevel)
}
