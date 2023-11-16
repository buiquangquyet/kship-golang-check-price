package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode   string `mapstructure:"mode"`
	AppKey string `mapstructure:"app_key"`

	Log struct {
		EnableFile bool   `mapstructure:"enable_file"`
		File       string `mapstructure:"file"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
	} `mapstructure:"log"`

	Server struct {
		Name    string `mapstructure:"name"`
		Address string `mapstructure:"address"`
		Prefix  string `mapstructure:"prefix"`
	} `mapstructure:"server"`

	Mysql struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		User        string `mapstructure:"user"`
		DbName      string `mapstructure:"db_name"`
		SslMode     string `mapstructure:"ssl_mode"`
		Password    string `mapstructure:"password"`
		AutoMigrate bool   `mapstructure:"auto_migrate"`
		MaxLifeTime int    `mapstructure:"max_life_time"`
	} `mapstructure:"mysql"`

	Mongo struct {
		Uri string `mapstructure:"uri"`
		DB  string `mapstructure:"db"`
	} `mapstructure:"mongo"`

	Redis struct {
		Hosts    []string `mapstructure:"hosts"`
		Username string   `mapstructure:"username"`
		Password string   `mapstructure:"password"`
	} `mapstructure:"redis"`

	ExtService struct {
		GHTK    *GHTK    `mapstructure:"ghtk"`
		AHAMOVE *AHAMOVE `mapstructure:"ahamove"`
		Voucher *Voucher `mapstructure:"voucher"`
	} `mapstructure:"ext_service"`

	Tracer struct {
		Enabled bool `mapstructure:"enabled"`
		Jaeger  struct {
			Endpoint string `mapstructure:"endpoint"`
			Active   bool   `mapstructure:"active"`
		} `mapstructure:"jaeger"`
	} `mapstructure:"tracer"`
}

type GHTK struct {
	Host  string `mapstructure:"host"`
	Token string `mapstructure:"token"`
}

type AHAMOVE struct {
	Host string `mapstructure:"host"`
	Key  string `mapstructure:"key"`
}

type Voucher struct {
	Host string `mapstructure:"host"`
}

var common *Config

func Get() *Config {
	return common
}

func LoadConfig(pathConfig string) error {
	viper.SetConfigFile(pathConfig)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&common)

	return nil
}
