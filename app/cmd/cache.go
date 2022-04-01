package cmd

import (
	"errors"
	"fmt"
	"gohub/pkg/cache"
	"gohub/pkg/console"

	"github.com/spf13/cobra"
)

var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
}

var CmdCacheClear = &cobra.Command{
	Use:   "clear",
	Short: "clear cache",
	Run:   runCacheClear,
}

var CmdCacheForget = &cobra.Command{
	Use:   "forget",
	Short: "delete redis key",
	Run:   runCacheForget,
}

var cacheKey string

func runCacheForget(cmd *cobra.Command, args []string) {
	cache.Forget(cacheKey)
	console.Success(fmt.Sprintf("Cache key [%s] deleted.", cacheKey))
}

func init() {
	CmdCache.AddCommand(
		CmdCacheClear,
	)

	CmdCacheForget.Flags().StringVarP(&cacheKey, "key", "k", "", "KEY of the cache")
	CmdCacheForget.MarkFlagRequired("key")
}

func runCacheClear(cmd *cobra.Command, args []string) {
	cache.Flush()
	console.Success("cache cleared.")
}

func runCache(cmd *cobra.Command, args []string) {

	console.Success("这是一条成功的提示")
	console.Warning("这是一条提示")
	console.Error("这是一条错误信息")
	console.Warning("终端输出最好使用英文，这样兼容性会更好~")
	console.Exit("exit 方法可以用来打印消息并中断程序！")
	console.ExitIf(errors.New("在 err = nil 的时候打印并退出"))
}
