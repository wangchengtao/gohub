package cmd

import (
	"github.com/spf13/cobra"
	"gohub/database/seeders"
	"gohub/pkg/console"
	"gohub/pkg/seed"
)

var CmdDBSeed = &cobra.Command{
	Use:   "seed",
	Short: "Insert fake data to the database",
	Run:   runSeeders,
	Args:  cobra.MaximumNArgs(1),
}

func runSeeders(cmd *cobra.Command, args []string) {
	seeders.Initialize()
	if len(args) > 0 {
		// 有传参的情况
		seed.RunSeeder(args[0])
	} else {
		seed.RunAll()
		console.Success("Done seeding.")
	}
}
