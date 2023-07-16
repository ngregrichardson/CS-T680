package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update [flags] updated_item_json",
	Aliases: []string{"u"},
	Short:   "Update a TODO item",
	Long:    "Update a TODO item",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		todo := initDB()
		fmt.Println("Running UPDATE_DB_ITEM...")
		item, err := todo.JsonToItem(args[0])
		if err != nil {
			fmt.Println("Update option requires a valid JSON todo item string")
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		if err := todo.UpdateItem(item); err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Ok")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
