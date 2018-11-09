package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
)

var Conf *tomlConfig

type tomlConfig struct {
	App      app
	Server   server
	Database database
	Redis    redis
}

type app struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	ImagePrefixUrl  string
	ImageSavePath   string
	ImageMaxSize    int
	ImageAllowTypes []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

type server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

type redis struct {
	Host      string
	Port      string
	Password  string
	MaxActive int
}

func Setup() {
	filePath, err := filepath.Abs("conf/app.toml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("parse toml file once. filePath: %s\n", filePath)

	viper.SetConfigName("app")   // name of config file (without extension)
	viper.AddConfigPath("conf/") // path to look for the config file in
	err2 := viper.ReadInConfig() // Find and read the config file
	if err2 != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.Unmarshal(&Conf)

	Conf.App.ImageMaxSize = Conf.App.ImageMaxSize * 1024 * 1024

	Conf.Server.ReadTimeout = Conf.Server.ReadTimeout * time.Second
	Conf.Server.WriteTimeout = Conf.Server.ReadTimeout * time.Second
}
