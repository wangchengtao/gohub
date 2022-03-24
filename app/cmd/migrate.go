package cmd

import (
	"github.com/spf13/cobra"
	"gohub/database/migrations"
	"gohub/pkg/migrate"
)

var CmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
	// 所有 migrate 下的子命令都会执行一下代码
}

var CmdMigrateUp = &cobra.Command{
	Use:   "up",
	Short: "Run unmigrated migrations",
	Run:   runUp,
}

var CmdMigrateDown = &cobra.Command{
	Use:   "down",
	Short: "Reverse the up command",
	Run:   runDown,
}

func runDown(cmd *cobra.Command, args []string) {
	migrator().Rollback()
}

func init() {
	CmdMigrate.AddCommand(
		CmdMigrateUp,
		CmdMigrateDown,
	)
}

func migrator() *migrate.Migrator {
	migrations.Initialize()

	return migrate.NewMigrator()
}

func runUp(cmd *cobra.Command, args []string) {
	migrator().Up()
}
