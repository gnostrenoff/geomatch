package configs

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// this is necessary to get the path of the config file regardless of the current working directory (in tests for example)
var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

var Config *viper.Viper

func init() {
	ResetByEnv()
}
func Init(env string) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(env)
	v.AutomaticEnv() // Give priority to ENV variables
	if env == "local" {
		v.AddConfigPath(basePath)
		err := v.ReadInConfig()
		if err != nil {
			log.Fatal("error on parsing configuration file ", err)
		}
	}
	Config = v
}

func ResetByEnv() {
	env := os.Getenv("ENV")
	Init(env)
}
