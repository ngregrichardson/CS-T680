package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all TODO items",
	Long:    "List all TODO items",
	Run: func(cmd *cobra.Command, args []string) {
		todo := initDB()
		fmt.Println("Running LIST_DB_ITEM...")
		todoList, err := todo.GetAllItems()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		for _, item := range todoList {
			todo.PrintItem(item)
		}
		fmt.Println("THERE ARE", len(todoList), "ITEMS IN THE DB")
		fmt.Println("Ok")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
