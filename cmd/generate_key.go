package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/gorilla/securecookie"
	"github.com/spf13/cobra"
)

// generateKeysCmd represents the hashPassword command
var generateKeysCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "Generate Key for Session Signing and CSRF Token",
	Run: func(cmd *cobra.Command, args []string) {
		key := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
		fmt.Printf("'%s'\r\n", key)
	},
}

func init() {
	rootCmd.AddCommand(generateKeysCmd)
}
