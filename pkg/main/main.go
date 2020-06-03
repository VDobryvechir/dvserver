// package main is the entry point
// DvServer copyright 2020 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)

package main

import (
	"github.com/Dobryvechir/microcore/pkg/dvconfig"
	_ "github.com/Dobryvechir/microcore/pkg/dvdbdata"
	_ "github.com/Dobryvechir/microcore/pkg/dvgeolocation"
	_ "github.com/Dobryvechir/microcore/pkg/dvlicense"
	_ "github.com/Dobryvechir/microcore/pkg/dvoc"
	_ "github.com/VDobryvechir/dvserver/pkg/dbloader"
	_ "github.com/VDobryvechir/dvserver/pkg/stock"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

func main() {
	dvconfig.SetApplicationName("DvServer")
	dvconfig.ServerStart()
}
