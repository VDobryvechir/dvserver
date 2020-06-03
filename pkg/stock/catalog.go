package stock

import (
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvoc"
	"log"
	"sort"
	"strconv"
	"strings"
)

func GetTopOfferingInfo(topOfferings [][]string, ctx *dvcontext.RequestContext) (*TopOffering, error) {
	n := len(topOfferings)
	if n == 0 {
		return nil, errors.New("Catalog is empty")
	}
	parentRefs := make(map[string]string)
	childRefs := make(map[string]string)
	totalPool := make(map[string]*PrimaryCategory)
	countTop := 0
	for i := 0; i < n; i++ {
		row := topOfferings[i]
		id := row[0]
		orderNumber, err1 := strconv.Atoi(row[2])
		if err1 != nil {
			orderNumber = 0
		}
		p := &PrimaryCategory{Id: id, Name: row[1], OrderNumber: orderNumber}
		parentRef := row[3]
		childRef := row[4]
		if parentRef != "" {
			parentRefs[parentRef] = id
		}
		if childRef != "" {
			childRefs[id] = childRef
		} else if totalPool[id] == nil {
			countTop++
		}
		totalPool[id] = p
	}
	topCategories := make([]*PrimaryCategory, 0, countTop)
	for id, primaryCategory := range totalPool {
		parentId := ""
		if childRefs[id] != "" {
			parentId = parentRefs[childRefs[id]]
			parent := totalPool[parentId]
			if parent == nil {
				log.Printf("Error catalog record id=%s child=%s parent=%s\n", id, childRefs[id], parentId)
				parentId = ""
			} else {
				if isDescendantCategory(primaryCategory, parent.Id) {
					log.Printf("Cyclic dependency omitted: %s (%s) and %s (%s)", id, primaryCategory.Name, parentId, parent.Name)
				} else {
					parent.SubCategories = append(parent.SubCategories, primaryCategory)
				}
			}
		}
		if parentId == "" {
			topCategories = append(topCategories, primaryCategory)
		}
	}
	if len(topCategories) == 0 {
		return nil, errors.New("No top level of categories")
	}
	sortPrimaryCategories(topCategories)
	categoryId := topCategories[0].Id
	subCategoryId := ""
	if len(topCategories[0].SubCategories) > 0 {
		subCategoryId = topCategories[0].SubCategories[0].Id
	}
	return &TopOffering{DefaultCategoryId: categoryId, DefaultSubCategoryId: subCategoryId, TopCategories: topCategories}, nil
}

func sortPrimaryCategories(categories []*PrimaryCategory) {
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].OrderNumber < categories[j].OrderNumber
	})
	n := len(categories)
	for i := 0; i < n; i++ {
		subcategories := categories[i].SubCategories
		if len(subcategories) > 0 {
			sortPrimaryCategories(subcategories)
		}
	}
}

func isDescendantCategory(category *PrimaryCategory, id string) bool {
	if category.Id == id {
		return true
	}
	n := len(category.SubCategories)
	for i := 0; i < n; i++ {
		if isDescendantCategory(category.SubCategories[i], id) {
			return true
		}
	}
	return false
}

func processInOutInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	command = strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if command == "" || command[0] != '{' || command[len(command)-1] != '}' {
		log.Printf("Invalid execution of get catalog command, config expected {}")
		return nil, false
	}
	cf := &InOutConfig{}
	err := json.Unmarshal([]byte(command), cf)
	if err != nil {
		log.Printf("Error in config %s: %v", command, err)
		return nil, false
	}
	return []interface{}{cf.In, cf.Out, ctx}, true
}

func processGetCatalogRun(data []interface{}) bool {
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

func RegisterActions() bool {
	dvoc.AddProcessFunction("catalog_offering", dvoc.ProcessFunction{
		Init: processInOutInit,
		Run:  processGetCatalogRun,
	})
	return true
}

var inited = RegisterActions()
