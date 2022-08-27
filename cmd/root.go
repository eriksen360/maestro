package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"maestro/storage"
)

var rootCmd = &cobra.Command{
	Use:   "maestro",
	Short: "Maestro is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load Configuration file into maestro",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ok, err := storage.SaveConfigFile(args[0])
		if !ok {
			fmt.Println(err)
			return err
		}
		return nil
	},
}

// maestro run <task> on <host>
// Spgs er: Hvordan aktiverer vi en task
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a maestro play through the executor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Do something
		return nil
	},
}

// maestro list hosts --all
// maestro list tasks --all
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list resources available in storage",
	Args:  cobra.MinimumNArgs(1), // Add [host | task] types
	RunE: func(cmd *cobra.Command, args []string) error {

		if args[0] == "hosts" {
			hosts := storage.GetAllHosts()
			fmt.Println("FQDN\t IP_ADDR\t GROUPS")
			for _, host := range hosts {
				fmt.Printf("%v\t%v\t%v\n", host.Fqdn, host.IpAddr, host.Groups)
			}

		} else if args[0] == "tasks" {
			tasks := storage.GetAllTasks()
			fmt.Println("NAME\tINSTRUCTIONS")
			for _, task := range tasks {
				fmt.Printf("%v\t%v\n", task.Name, task.Instructions)
			}

		} else {
			fmt.Println("error: list argument must be 'hosts' or 'tasks'")
			return nil
		}
		return nil
	},
}

// Init Initialize sub commands
func Init() {
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(listCmd)
}

// Execute Execute CLI interface
func Execute() error {
	Init()
	return rootCmd.Execute()
}
