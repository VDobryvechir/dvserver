package tools

import (
	"github.com/Dobryvechir/microcore/pkg/dvoc"
	"github.com/VDobryvechir/dvserver/pkg/tools/ug"
)

func RegisterActions() bool {
	dvoc.AddProcessFunction("ug_config", dvoc.ProcessFunction{
		Init: ug.ConfigInit,
		Run:  ug.ConfigRun,
	})
	return true
}

var inited = RegisterActions()

