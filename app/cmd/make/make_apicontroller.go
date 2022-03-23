package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"gohub/pkg/console"
	"strings"
)

var CmdMakeAPIController = &cobra.Command{
	Use:   "apicontroller",
	Short: "Create api controller, example: make apicontroller v1/user",
	Run:   runMakeAPIControler,
	Args:  cobra.ExactArgs(1),
}

func runMakeAPIControler(cmd *cobra.Command, args []string) {
	// 处理参数, 要求附带 api 版本
	array := strings.Split(args[0], "/")
	if len(array) != 2 {
		console.Exit("api controller name format: v1/user")
	}

	apiVersion, name := array[0], array[1]
	model := makeModelFromString(name)

	// 组建目标目录
	filePath := fmt.Sprintf("app/http/controllers/api/%s/%s_controller.go", apiVersion, model.TableName)

	createFileFromStub(filePath, "apicontroller", model)
}
