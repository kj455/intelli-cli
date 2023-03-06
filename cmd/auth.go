package cmd

import (
	"github.com/kj455/intelli-cli/internal/secret"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticates you to use chatGPT API through intelli-cli",
	Long:  `Authenticates you to use chatGPT API through intelli-cli`,
}

var AuthLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to chatGPT API",
	Long:  `login to chatGPT API`,
	Run: func(cmd *cobra.Command, args []string) {
		secret.SetupSecretIfNeeded()
	},
}

var AuthLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout from chatGPT API",
	Long:  `logout from chatGPT API`,
	Run: func(cmd *cobra.Command, args []string) {
		secret.DeleteApiKey()
	},
}

var AuthRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "refresh your chatGPT API key",
	Long:  `refresh your chatGPT API key`,
	Run: func(cmd *cobra.Command, args []string) {
		secret.DeleteApiKey()
		secret.SetupSecretIfNeeded()
	},
}

func init() {
	authCmd.AddCommand(AuthLoginCmd)
	authCmd.AddCommand(AuthLogoutCmd)
	authCmd.AddCommand(AuthRefreshCmd)
}
