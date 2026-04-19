package main

import (
	"fmt"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	memberService "github.com/flipped-aurora/gin-vue-admin/server/service/member"
)

func main() {
	global.GVA_VP = core.Viper()
	global.GVA_DB = initialize.Gorm()
	if global.GVA_DB == nil {
		fmt.Fprintln(os.Stderr, "database is not configured")
		os.Exit(1)
	}

	sqlDB, err := global.GVA_DB.DB()
	if err == nil {
		defer func() {
			_ = sqlDB.Close()
		}()
	}

	summary, err := memberService.MigrateLegacyMemberData(global.GVA_DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate legacy member data failed: %v\n", err)
		os.Exit(1)
	}

	if err := memberService.Bootstrap(); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap member resources failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("legacy member migration completed: %s\n", memberService.LoadLegacyMigrationSummaryText(summary))
}
