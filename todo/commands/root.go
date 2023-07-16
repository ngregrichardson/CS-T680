package commands

import (
	"fmt"
	"os"

	"drexel.edu/todo/db"
	"github.com/spf13/cobra"
)

var (
	dbFileNameFlag string
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Basic TODO CLI",
	Long:  `A simple TODO app CLI to keep track of the state of tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbFileNameFlag, "file", "f", "./data/todo.json", "Name of the database file")
}

func initDB() *db.ToDo {
	todo, err := db.New(dbFileNameFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return todo
}

func Execute() error {
	return rootCmd.Execute()
}
