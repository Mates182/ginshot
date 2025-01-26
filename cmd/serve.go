package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Build and run the project using Docker Compose",
	Long: `The 'serve' command builds and starts the project containers in detached mode 
using Docker Compose. It ensures that all services defined in the docker-compose.yml file are up and running.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the project with Docker Compose...")

		// Check if docker-compose.yml exists
		if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
			fmt.Println("Error: docker-compose.yml not found in the current directory.")
			return
		}

		// Execute docker-compose command
		command := exec.Command("docker-compose", "up", "--build", "-d")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			fmt.Println("Error running Docker Compose:", err)
			return
		}

		fmt.Println("Project is running successfully.")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
