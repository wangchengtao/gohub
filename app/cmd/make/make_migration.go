package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"gohub/pkg/app"
	"gohub/pkg/console"
)

var CmdMakeMigration = &cobra.Command{
	Use:   "migration",
	Short: "",
	Run:   runMakeMigration,
	Args:  cobra.ExactArgs(1),
}

func runMakeMigration(cmd *cobra.Command, args []string) {
	timeStr := app.TimenowInTimezone().Format("2006-01-02_150405")

	model := makeModelFromString(args[0])
	fileName := timeStr + "_" + model.PackageName
	filePath := fmt.Sprintf("database/migrations/%s.go", fileName)

	createFileFromStub(filePath, "migration", model, map[string]string{"{{FileName}}": fileName})

	console.Success("Migration file created, after modify it, use `migrate up` to migrate database.")
}
