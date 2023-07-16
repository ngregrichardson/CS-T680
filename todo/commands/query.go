package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:     "query [flags] item_id",
	Aliases: []string{"q"},
	Short:   "Query a TODO item",
	Long:    "Query a TODO item",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		todo := initDB()
		fmt.Println("Running QUERY_DB_ITEM...")

		id, err := todo.GetID(args[0])

		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		item, err := todo.GetItem(id)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		todo.PrintItem(item)
		fmt.Println("Ok")
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
