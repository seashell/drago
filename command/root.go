package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/seashell/drago/command/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "drago",
	Short: "A flexible configuration manager for Wireguard networks",
	Long:  `Usage: drago [-version] [-help] [-autocomplete-(un)install] <command> [args]`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "/etc/drago.d", "config file (default is /etc/drago.d)")

	rootCmd.AddCommand(agent.Command)
}

func initConfig() {

	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	LoadConfig(cfgFile)

	setDefaultAndEnv("server.enabled", true)
	setDefaultAndEnv("server.ui.enabled", true)
	setDefaultAndEnv("server.bindaddr", ":8080")
	setDefaultAndEnv("server.acl.secret", "")
}

func setDefaultAndEnv(key string, defaultConfig interface{}) {
	viper.SetDefault(key, defaultConfig)
	viper.Set(key, viper.Get(key))
}

func LoadConfig(path string) error {

	// Check if dir or file
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return LoadConfigDir(path)
	}

	return LoadConfigFile(path)
}

func LoadConfigDir(dir string) error {

	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("configuration path must be a directory: %s", dir)
	}

	file := ""
	err = nil
	for err != io.EOF {
		var fis []os.FileInfo
		fis, err = f.Readdir(128)
		if err != nil && err != io.EOF {
			return err
		}

		for _, fi := range fis {
			// Ignore directories
			if fi.IsDir() {
				continue
			}

			// Only care about files that are valid to load.
			name := fi.Name()
			skip := true
			if strings.HasSuffix(name, ".yml") {
				skip = false
			}
			if skip || file != "" {
				continue
			}

			path := filepath.Join(dir, name)
			file = path
		}
	}

	// Fast-path if we have no files
	if file == "" {
		return nil
	}

	err = LoadConfigFile(file)
	if err != nil {
		return fmt.Errorf("Error loading %s: %s", file, err)
	}

	return nil
}

func LoadConfigFile(path string) error {

	viper.SetConfigFile(path)
	err := viper.MergeInConfig()

	// If a config file is found, read it in.
	if err == nil {
		fmt.Println("==> Loaded configuration from ", viper.ConfigFileUsed())
	} else {
		fmt.Println("==> Error loading configuration: ", err)
	}
	return err
}
