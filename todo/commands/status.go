package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status [flags] item_id new_status",
	Aliases: []string{"s"},
	Short:   "Update the status of a TODO item",
	Long:    "Update the status of a TODO item",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		todo := initDB()
		fmt.Println("Running CHANGE_ITEM_STATUS...")

		id, err := todo.GetID(args[0])

		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		value, err := strconv.ParseBool(args[1])

		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		err = todo.ChangeItemDoneStatus(id, value)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Ok")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
