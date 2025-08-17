package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(opsbrew completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ opsbrew completion bash > /etc/bash_completion.d/opsbrew
  # macOS:
  $ opsbrew completion bash > /usr/local/etc/bash_completion.d/opsbrew

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ opsbrew completion zsh > "${fpath[1]}/_opsbrew"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ opsbrew completion fish | source

  # To load completions for each session, execute once:
  $ opsbrew completion fish > ~/.config/fish/completions/opsbrew.fish

PowerShell:
  PS> opsbrew completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> opsbrew completion powershell > opsbrew.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
				cmd.PrintErrf("Error generating bash completion: %v\n", err)
			}
		case "zsh":
			if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
				cmd.PrintErrf("Error generating zsh completion: %v\n", err)
			}
		case "fish":
			if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
				cmd.PrintErrf("Error generating fish completion: %v\n", err)
			}
		case "powershell":
			if err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout); err != nil {
				cmd.PrintErrf("Error generating powershell completion: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
