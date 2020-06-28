package cmd

import (
	"github.com/Matt-Gleich/statuser"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "texsch",
	Long: `
Automation for techy students that write papers for high school in LaTeX

🐙 Repository: https://github.com/Matt-Gleich/texsch
📟 Authors:
	- Matthew Gleich (@Matt-Gleich)

████████╗███████╗██╗  ██╗███████╗ ██████╗██╗  ██╗
╚══██╔══╝██╔════╝╚██╗██╔╝██╔════╝██╔════╝██║  ██║
   ██║   █████╗   ╚███╔╝ ███████╗██║     ███████║
   ██║   ██╔══╝   ██╔██╗ ╚════██║██║     ██╔══██║
   ██║   ███████╗██╔╝ ██╗███████║╚██████╗██║  ██║
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝
	`,
}

// Execute ... Execute the main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		statuser.Error("Failed to execute root command", err, 1)
	}
}
