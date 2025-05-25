package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dojer",
	Short: "the last nh client that you'll ever need",
	Long:  `nh client for downloading, serving and storing dojs :)`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	defaultDir := getDataDir()

	viper.SetDefault("data_dir", defaultDir)

	viper.SetDefault("backup.interval", 4)
	viper.SetDefault("nhentai.user_agent", "")
	viper.SetDefault("nhentai.cookies", []string{})
	viper.SetDefault("server.port", port)
	viper.SetConfigName("dojer")
	viper.SetConfigType("toml")
	viper.AddConfigPath(defaultDir)

	_ = viper.SafeWriteConfig()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func getDataDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "./"
	}

	var appName string
	if os.Getenv("DOJER_ENV") == "development" {
		appName = os.Getenv("DOJER_NAME")
	} else {
		appName = "dojer"
	}

	dataDir := filepath.Join(configDir, appName)

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Println(err.Error())
		log.Fatal("Could not create default folder")
		os.Exit(1)
	}
	return dataDir
}
