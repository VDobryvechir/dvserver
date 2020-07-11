package stock

import (
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvoc"
	"log"
	"strings"
)

func processInOutInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	info:= strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if info == "" || info[0] != '{' || info[len(info)-1] != '}' {
		log.Printf("Invalid parameters of command %s, config expected {}", command)
		return nil, false
	}
	cf := &InOutConfig{}
	err := json.Unmarshal([]byte(info), cf)
	if err != nil {
		log.Printf("Error in config %s: %v", command, err)
		return nil, false
	}
	return []interface{}{cf.In, cf.Out, ctx}, true
}

func RegisterActions() bool {
	dvoc.AddProcessFunction("offering_catalog", dvoc.ProcessFunction{
		Init: processInOutInit,
		Run:  processGetCatalogRun,
	})
	dvoc.AddProcessFunction("top_offerings", dvoc.ProcessFunction{
		Init: processInOutInit,
		Run:  processGetTopOfferingRun,
	})
	return true
}

var inited = RegisterActions()
