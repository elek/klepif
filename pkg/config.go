package pkg

import (
	"github.com/spf13/viper"
	"os"
)

type JenkinsConfig struct {
	Username string
	Url      string
	Token    string
}

type GithubConfig struct {
	Token string
	Org   string
	Repo  string
}

func ReadConfig() (error) {
	dir, err := os.Getwd()
	if err != nil {
		viper.SetDefault("persistence.path", dir);
	}
	viper.SetConfigName("klepif")
	viper.AddConfigPath("/etc/klepif/")
	viper.AddConfigPath("$HOME/.config/klepif")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	return err
}
