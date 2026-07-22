package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/igsafe/bwid"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var length int
var count int

var generateTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate random string token",
	Run: func(cmd *cobra.Command, args []string) {
		if length == 0 {
			length = 40
		}
		for i := 0; i < count; i++ {
			fmt.Println(bwid.GenerateToken(length))
		}
	},
}

var generateOidCmd = &cobra.Command{
	Use:   "oid",
	Short: "Generate timestamped string token",
	Run: func(cmd *cobra.Command, args []string) {
		if length == 0 {
			length = 24
		}
		for i := 0; i < count; i++ {
			fmt.Println(bwid.GenerateTimestampedToken(length))
		}
	},
}

func readPassword(prompt string) string {
	fmt.Fprint(os.Stderr, prompt)
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr)
	return strings.TrimRight(string(pw), "\r\n")
}

var generateHashCmd = &cobra.Command{
	Use:   "hash <password>",
	Short: "Generate a bcrypt hash for a password",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var password string
		if len(args) == 1 {
			password = args[0]
		} else {
			password = readPassword("Password: ")
			confirm := readPassword("Confirm:  ")
			if password != confirm {
				fmt.Fprintln(os.Stderr, "Passwords do not match")
				os.Exit(1)
			}
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(hash))
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate string tokens using bwid",
}

func init() {
	generateCmd.AddCommand(generateTokenCmd)
	generateCmd.AddCommand(generateOidCmd)
	generateCmd.AddCommand(generateHashCmd)
	generateCmd.PersistentFlags().IntVarP(&count, "count", "c", 1, "Number of tokens to generate")
	generateCmd.PersistentFlags().IntVarP(&length, "length", "l", 0, "Length of tokens to generate")
	rootCmd.AddCommand(generateCmd)
}
