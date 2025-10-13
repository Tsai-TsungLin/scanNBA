// ============ cmd/root.go ============
package cmd

import (
	"log"
	"nba-scanner/internal/logic"
	"nba-scanner/internal/server"

	"github.com/spf13/cobra"
)

var (
	startTime  string
	serverMode bool
	port       int
)

var rootCmd = &cobra.Command{
	Use:   "nba-scan",
	Short: "NBA 資訊掃描工具",
	Run: func(cmd *cobra.Command, args []string) {
		if serverMode {
			// 啟動 Web Server
			if err := server.Start(port); err != nil {
				log.Fatalf("啟動 server 失敗: %v", err)
			}
		} else {
			// CLI 模式
			if startTime != "" {
				logic.PKTeamOnStartTime(startTime)
			} else {
				logic.PKTeam()
			}
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&startTime, "time", "", "", "指定時間 (格式: 15:04)")
	rootCmd.PersistentFlags().BoolVarP(&serverMode, "server", "s", false, "啟動 Web Server 模式")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Web Server 埠號")
	rootCmd.Execute()
}
