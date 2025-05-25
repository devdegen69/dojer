package cmd

import (
	"dojer/store"
	"dojer/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [id1] [id2] ...",
	Short: "delete a doujinshi",
	Long:  `if you'll ever want to delete something completely from the database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := store.Delete(args)
		if err != nil {
			fmt.Println(err)
		}

		for _, arg := range args {
			err := store.RemoveFromIndex(arg)
			if err != nil {
				fmt.Println(err)
			}
			path := utils.GetDataPath("downloads", arg)
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("could not delete %s, please go into the downloads/%s folder and delete it manually", arg, arg)
			} else {
				fmt.Printf("folder downloads/%s deleted!", arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
