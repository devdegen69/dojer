package cmd

import (
	"dojer/downloader"
	"dojer/extractors"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [id]|[url]",
	Short: "download a doj and add it do the gallery",
	Long:  `download the doj by the url or the id from the official site`,
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			for {
				pr, _ := downloader.GetPipe()
				io.Copy(os.Stdout, pr)
			}
		}()

		extractors.Run(args, false)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("it must have at least one argument")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
