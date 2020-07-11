package stock

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"log"
)

func processGetTopOfferingRun(data []interface{}) bool {
	inProp := data[0].(string)
	outProp := data[1].(string)
	ctx := data[2].(*dvcontext.RequestContext)
	topOfferingInfo, ok := ctx.ExtraAsDvObject.Get(inProp)
	if !ok {
		ctx.HandleFileNotFound()
		return false
	}
	topOfferings := topOfferingInfo.([][]string)
	oeOfferings, err := GetTopOfferingInfo(topOfferings, ctx)
	if err != nil {
		log.Printf("Failed to get catalog %v", err)
		return false
	}
	ctx.ExtraAsDvObject.Set(outProp, oeOfferings)
	return true
}
