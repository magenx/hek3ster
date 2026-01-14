package commands

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for the specified shell.

Examples:
  # Bash
  source <(hek3ster completion bash)
  # To load completions for each session, execute once:
  hek3ster completion bash > /etc/bash_completion.d/hek3ster

  # Zsh
  source <(hek3ster completion zsh)
  # To load completions for each session, execute once:
  hek3ster completion zsh > "${fpath[1]}/_hek3ster"

  # Fish
  hek3ster completion fish | source
  # To load completions for each session, execute once:
  hek3ster completion fish > ~/.config/fish/completions/hek3ster.fish

  # PowerShell
  hek3ster completion powershell | Out-String | Invoke-Expression
  # To load completions for each session, execute once:
  hek3ster completion powershell > hek3ster.ps1
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
		}
		return nil
	},
}

func init() {
	// This command will be added to the root command in root.go
	//
}
