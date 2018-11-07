package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
)

type app struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowTypes []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &app{}

type server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &server{}

type database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &database{}

type redisSetting struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &redisSetting{}

type tomlConfig struct {
	App      app
	Server   server
	Database database
	RedisSetting redisSetting
}


func Setup() {
	filePath, err := filepath.Abs("conf/app.toml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("parse toml file once. filePath: %s\n", filePath)

	var cfg tomlConfig
	viper.SetConfigName("app")   // name of config file (without extension)
	viper.AddConfigPath("conf/") // path to look for the config file in
	err2 := viper.ReadInConfig() // Find and read the config file
	if err2 != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.Unmarshal(&cfg)

	AppSetting = &cfg.App
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	ServerSetting = &cfg.Server
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second

	DatabaseSetting = &cfg.Database

	RedisSetting = &cfg.RedisSetting
}
