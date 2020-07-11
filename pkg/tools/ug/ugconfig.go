package ug

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strconv"
	"strings"
)

func ConfigInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	info := strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if strings.ToLower(info) != "enabled" {
		log.Printf("ug config is not enabled")
		return nil, false
	}
	ids := ctx.ExtraAsDvObject.GetString("G_IDS")
	idList := dvparser.ConvertToNonEmptyList(ids)
	if len(idList) == 0 {
		log.Printf("ids is empty")
		return nil, false
	}
	return []interface{}{idList}, true
}

func ConfigRun(info []interface{}) bool {
	ids := info[0].([]string)
	ctx := info[1].(*dvcontext.RequestContext)
	var sql string
	buf := make([]byte, 0, 120000)
	db, err := dvdbdata.GetDBConnection("TM")
	if err != nil {
		log.Printf("Error connecting %v", err)
		return false
	}
	objIds := strings.Join(ids, ",")
	sql = "select o.OBJECT_ID,o.NAME,o.OBJECT_TYPE_ID,t.NAME,o.PARENT_ID from NC_OBJECTS o JOIN NC_OBJECT_TYPES t ON o.OBJECT_TYPE_ID=t.OBJECT_TYPE_ID  START WITH o.OBJECT_ID in (" + objIds + ") CONNECT BY o.OBJECT_ID=prior o.PARENT_ID"
	_, buf, err = presentTableInfo(buf, sql, db, "PARENT", 5, 4)
	sql = "select o.OBJECT_ID,o.NAME,o.OBJECT_TYPE_ID,t.NAME,o.PARENT_ID from NC_OBJECTS o JOIN NC_OBJECT_TYPES t ON o.OBJECT_TYPE_ID=t.OBJECT_TYPE_ID  START WITH o.OBJECT_ID in (" + objIds + ") CONNECT BY prior o.OBJECT_ID=o.PARENT_ID"
	var fullIdsInfo []byte
	fullIdsInfo, buf, err = presentTableInfo(buf, sql, db, "O", 5, 5)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.HandleInternalServerError()
		return true
	}
	fullIds := string(fullIdsInfo)
	sql = "select o.OBJECT_ID,o.VALUE,o.ATTR_ID,t.NAME,o.LIST_VALUE_ID,o.DATE_VALUE,v.VALUE,v.ATTR_TYPE_DEF_ID from NC_PARAMS o JOIN NC_ATTRIBUTES t ON o.ATTR_ID=t.ATTR_ID LEFT OUTER JOIN NC_LIST_VALUES v ON o.LIST_VALUE_ID=v.LIST_VALUE_ID WHERE o.OBJECT_ID in (" + fullIds + ")"
	_, buf, err = presentTableInfo(buf, sql, db, "A", 8, 2)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.HandleInternalServerError()
		return true
	}
	sql = "select o.OBJECT_ID,v.NAME,o.ATTR_ID,t.NAME,o.REFERENCE from NC_REFERENCES o JOIN NC_ATTRIBUTES t ON o.ATTR_ID=t.ATTR_ID LEFT OUTER JOIN NC_OBJECTS v ON o.REFERENCE=v.OBJECT_ID WHERE o.OBJECT_ID in (" + fullIds + ")"
	_, buf, err = presentTableInfo(buf, sql, db, "R", 5, 4)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.HandleInternalServerError()
		return true
	}
	sql = "select o.OBJECT_ID,v.NAME,o.ATTR_ID,t.NAME,o.REFERENCE from NC_REFERENCES o JOIN NC_ATTRIBUTES t ON o.ATTR_ID=t.ATTR_ID LEFT OUTER JOIN NC_OBJECTS v ON o.OBJECT_ID=v.OBJECT_ID WHERE o.REFERENCE in (" + fullIds + ")"
	_, buf, err = presentTableInfo(buf, sql, db, "REF_BY", 5, 4)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.HandleInternalServerError()
		return true
	}
	ctx.Output = buf
	return true
}

func presentTableInfo(nbuf []byte, query string, db *dvdbdata.DBConnection, pref string, colAmount int, options int) (collector []byte, buf []byte, err error) {
	buf = nbuf
	var rs [][]interface{}
	rs, err = dvdbdata.GetSqlTableByQuery(db, nil, query)
	if err != nil {
		return
	}
	collector = make([]byte, 0, 10240)
	count := 0
	collect := (options & 1) != 0
	next := false
	startNextColumn := 5
	params := (options & 2) != 0
	nextIsObjId := (options & 4) != 0
	if params {
		startNextColumn = 8
	}
	n := len(rs)
	for i := 0; i < n; i++ {
		res := rs[i]
		id := dvevaluation.AnyToString(res[0])
		if collect {
			if next {
				collector = append(collector, ',')
			} else {
				next = true
			}
			collector = append(collector, id...)
		}
		name := dvevaluation.AnyToString(res[1])
		tpId := dvevaluation.AnyToString(res[2])
		tpName := dvevaluation.AnyToString(res[3])
		if count == 0 {
			buf = append(buf, []byte("<table id='packetInfo'>")...)
		}
		if params {
			s := dvevaluation.AnyToString(res[4])
			if len(s) > 0 {
				tpId += "_LIST"
				name = s + " (" + dvevaluation.AnyToString(res[6]) + ")"
			} else {
				s = dvevaluation.AnyToString(res[5])
				if len(s) > 0 {
					tpId += "_DATE"
					name = s
				}
			}
		}
		buf = append(buf, []byte("<tr><td><a target='_blank' href='/ncobject.jsp?id=")...)
		buf = append(buf, []byte(id)...)
		buf = append(buf, []byte("'>")...)
		buf = append(buf, []byte(id)...)
		buf = append(buf, []byte("</a></td><td class='extra'>")...)
		buf = append(buf, []byte(name)...)
		buf = append(buf, []byte("</td><td>")...)
		buf = append(buf, []byte(pref)...)
		buf = append(buf, '_')
		buf = append(buf, []byte(tpId)...)
		buf = append(buf, []byte("</td><td>")...)
		buf = append(buf, []byte(tpName)...)
		for i := startNextColumn; i <= colAmount; i++ {
			s := dvevaluation.AnyToString(res[i-1])
			buf = append(buf, []byte("</td><td>")...)
			if i == startNextColumn {
				if nextIsObjId {
					s = "<a target='_blank' href='/ncobject.jsp?id=" + s + "'>" + s + "</a>"
				}
			}
			buf = append(buf, []byte(s)...)
		}
		buf = append(buf, []byte("</td></tr>")...)
		count++
	}
	if count == 0 {
		buf = append(buf, []byte("<div style='color:red'>No "+pref+" for these ids</div>")...)
	} else {
		buf = append(buf, []byte("</table><div class='extra' style='color:green'>Total objects "+pref+" :")...)
		buf = append(buf, []byte(strconv.Itoa(count))...)
		buf = append(buf, []byte("</div>")...)
	}
	return
}
