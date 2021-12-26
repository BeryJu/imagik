package cmd

import (
	"fmt"
	"syscall"

	"beryju.org/imagik/pkg/drivers/auth"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// hashPasswordCmd represents the hashPassword command
var hashPasswordCmd = &cobra.Command{
	Use:   "hash-password",
	Short: "Hash password for usage with static Authentication",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter password to hash")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		hash, err := auth.HashPassword(string(bytePassword))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hash is '%s'\r\n", hash)
	},
}

func init() {
	rootCmd.AddCommand(hashPasswordCmd)
}
