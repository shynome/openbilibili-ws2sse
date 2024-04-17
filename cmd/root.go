/*
Copyright © 2024 shynome <shynome@gmail.com>
*/
package cmd

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/shynome/openbilibili-ws2sse/sse"
	"github.com/spf13/cobra"
)

var args struct {
	listen string
	key    string
	secret string
	appid  int64
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openbilibili-ws2sse",
	Short: "将B站开放平台的弹幕转为Event Stream",
	// Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, _args []string) {
		if args.appid == 0 || args.key == "" || args.secret == "" {
			slog.Error("key, secret, appid is required")
			return
		}
		l, err := net.Listen("tcp", args.listen)
		if err != nil {
			slog.Error("端口监听失败", "err", err)
			return
		}
		defer l.Close()
		ctx := context.Background()
		srv := sse.New(args.key, args.secret, args.appid)
		go srv.BatchKeepAlive(ctx)

		slog.Info("服务启动成功", "addr", l.Addr().String())
		if err := http.Serve(l, srv); err != nil {
			slog.Error("服务监听退出", "err", err)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.openbilibili-ws2sse.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&args.listen, "listen", ":7070", "监听的地址")
	rootCmd.Flags().StringVar(&args.key, "key", "", "access_key_id")
	rootCmd.Flags().StringVar(&args.secret, "secret", "", "access_key_secret")
	rootCmd.Flags().Int64Var(&args.appid, "appid", 0, "项目ID")
}
