package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
)

var conf *TomlConfig

type TomlConfig struct {
	App      app
	Server   server
	Database database
	Redis    map[string]redis
}

type app struct {
	JwtSecret string
	PageSize  int

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
	Type         string
	User         string
	Password     string
	Host         string
	Name         string
	TablePrefix  string
	MaxIdleConns int
	MaxOpenConns int
}

type redis struct {
	Host      string
	Port      string
	Password  string
	MaxActive int
}

func init() {
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
	viper.Unmarshal(&conf)

	conf.App.ImageMaxSize = conf.App.ImageMaxSize * 1024 * 1024

	conf.Server.ReadTimeout = conf.Server.ReadTimeout * time.Second
	conf.Server.WriteTimeout = conf.Server.ReadTimeout * time.Second
}

func New() *TomlConfig {
	return conf
}

func (t *TomlConfig) GetString(key string) string {
	return viper.GetString(key)
}

func (t *TomlConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

func (t *TomlConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (t *TomlConfig) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
