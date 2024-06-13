package cmd

import (
	"fmt"

	"github.com/igsafe/bwid"
	"github.com/spf13/cobra"
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

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate string tokens using bwid",
}

func init() {
	generateCmd.AddCommand(generateTokenCmd)
	generateCmd.AddCommand(generateOidCmd)
	rootCmd.PersistentFlags().IntVarP(&count, "count", "c", 1, "Number of tokens to generate")
	rootCmd.PersistentFlags().IntVarP(&length, "length", "l", 0, "Length of tokens to generate")
	rootCmd.AddCommand(generateCmd)
}
