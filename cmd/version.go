/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"io/ioutil"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version of the application",
	Long:  `Show the version of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := ioutil.ReadFile("VERSION")

		if err != nil {
			return
		}

		fmt.Println(string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
