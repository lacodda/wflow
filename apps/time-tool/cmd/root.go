package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "time-tool",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(LoginCmd)
	rootCmd.AddCommand(LogoutCmd)
	rootCmd.AddCommand(TimestampCmd)
	rootCmd.AddCommand(TaskCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(ReportCmd)

	TimestampCmd.Flags().VarP(&FlagTimestampType, "type", "t", `Type of timestamp. Allowed: "Start", "End", "StartBreak", "EndBreak"`)
	TimestampCmd.Flags().BoolVarP(&FlagTimestampShow, "show", "s", false, "Show a list of timestamps for the selected day")
	TimestampCmd.Flags().BoolVarP(&FlagTimestampRaw, "raw", "r", false, "Show a raw list of timestamps for the selected day")
	TimestampCmd.Flags().StringVarP(&FlagTimestampDate, "date", "d", "", `Select a date to display: e.g. "2023-01-19"`)
	TaskCmd.Flags().BoolVarP(&FlagTaskShow, "show", "s", false, "Show a list of tasks for the selected day")
	TaskCmd.Flags().StringVarP(&FlagTaskDate, "date", "d", "", `Select a date to display: e.g. "2023-01-19"`)
	SyncCmd.Flags().BoolVarP(&FlagSyncShow, "show", "s", false, "Show records for synchronizing with the server")
	SyncCmd.Flags().IntSliceVarP(&FlagSyncDelete, "delete", "D", []int{}, "ID records for deleting")
	ReportCmd.Flags().StringVarP(&FlagReportDate, "date", "d", "", `Select a date to display: e.g. "2023-01-19"`)
	ReportCmd.Flags().BoolVar(&FlagReportSend, "send", false, "Send report")
	ReportCmd.Flags().BoolVar(&FlagReportTestSend, "test", false, "Send test report")
}
