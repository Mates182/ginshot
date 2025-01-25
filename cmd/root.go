/*
Copyright © 2025 mates182
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Grey    = "\033[38;5;245m"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ginshot",
	Short: "A brief description of your application",
	Long: `Ginshot is a command-line tool designed to generate a structured 
boilerplate for building APIs and microservices with the Gin framework. 
It helps developers by scaffolding routes, layered architecture files, 
and essential configurations to accelerate project setup and ensure best 
practices. `,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
