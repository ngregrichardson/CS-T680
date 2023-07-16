package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [flags] item_id",
	Aliases: []string{"d"},
	Short:   "Delete a TODO item",
	Long:    "Delete a TODO item",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		todo := initDB()
		fmt.Println("Running DELETE_DB_ITEM...")

		id, err := todo.GetID(args[0])

		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		err = todo.DeleteItem(id)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Ok")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
