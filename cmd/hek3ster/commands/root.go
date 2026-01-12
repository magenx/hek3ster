package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string

// Execute executes the root command
func Execute(ver string) error {
	version = ver
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "hek3ster",
	Short: "A tool to create k3s clusters on Hetzner Cloud",
	Long: `hek3ster - The easiest and fastest way to create production-ready Kubernetes clusters on Hetzner Cloud

hek3ster is a CLI tool that creates production-ready Kubernetes clusters 
on Hetzner Cloud in minutes. No Terraform, Packer, or Ansible required.`,
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(releasesCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(versionCmd)

	// Disable auto-generated completion command (we have our own)
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
		fmt.Printf("Version: %s\n", version)
	},
}

func printBanner() {
	green := "\033[32m"
	blue := "\033[34m"
	reset := "\033[0m"

	if os.Getenv("NO_COLOR") != "" {
		green = ""
		blue = ""
		reset = ""
	}

	fmt.Printf("%s _   _      _    _____     _            %s\n", green, reset)
	fmt.Printf("%s| | | | ___| | _|___ / ___| |_ ___ _ __ %s\n", green, reset)
	fmt.Printf("%s| |_| |/ _ \\ |/ / |_ \\/ __| __/ _ \\ '__|%s\n", green, reset)
	fmt.Printf("%s|  _  |  __/   < ___) \\__ \\ ||  __/ |   %s\n", green, reset)
	fmt.Printf("%s|_| |_|\\___|_|\\_\\____/|___/\\__\\___|_|   %s\n", green, reset)
	fmt.Println()
	fmt.Printf("%sVersion: %s%s\n", blue, version, reset)
	fmt.Println()
}

func printSponsorMessage() {
	blue := "\033[34m"
	reset := "\033[0m"

	if os.Getenv("NO_COLOR") != "" {
		blue = ""
		reset = ""
	}

	fmt.Println()
	fmt.Printf("%s=======================================================%s\n", blue, reset)
	fmt.Printf("%s  Do you like hek3ster? Support its development:%s\n", blue, reset)
	fmt.Printf("%s  https://github.com/magenx/hek3ster%s\n", blue, reset)
	fmt.Printf("%s=======================================================%s\n", blue, reset)
	fmt.Println()
}
