// package dbloader loads database to csv format
// copyright Volodymyr Dobryvechir 2020

package dbloader

import (
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvoc"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strings"
)

func ProvideDbImportCommand() {
	dvoc.AddProcessFunction("db", dvoc.ProcessFunction{
		Init: processDbImportInit,
		Run:  processDbImportRun,
	})
}

func Init() bool {
	ProvideDbImportCommand()
	return true
}

var inited = Init()

func processDbImportInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
	command = strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if command == "" || command[0] != '{' || command[len(command)-1] != '}' {
		log.Printf("Invalid execution dbimport command, import name expected ")
		return nil, false
	}
	if dvparser.GlobalProperties[dbImport+command] == "" {
		log.Printf("Property %s is not defined for import", dbImport+command)
		return nil, false
	}
	return []interface{}{command}, true
}

func processDbImportRun(data []interface{}) bool {
	name := data[0].(string)
	err := LoadRelatedTablesByImportName(name)
	if err != nil {
		log.Printf("Failed to import %s: %v", name, err)
		return false
	}
	return true
}
