package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string //for okex.com v3 api
	ClientId      string //for bitstamp.net , huobi.pro

	Lever float64 //杠杆倍数 , for future
}

var once sync.Once
var instance *ConfigEngine

type ConfigEngine struct {
}

func NewConfigEngine() *ConfigEngine {
	return &ConfigEngine{}
}

func GetConfigEngine() *ConfigEngine {
	once.Do(func() {
		instance = NewConfigEngine()
	})
	return instance
}

func (configEngine *ConfigEngine) ReadConfig() *APIConfig {
	//curPath, _ := os.Getwd()
	viper.SetConfigFile("/Users/ran/SourceCode/quant/scindapsus/config/api.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("[config] parse config file failed:%s", err.Error())
		return nil
	}
	apiConfig := &APIConfig{
		HttpClient:    nil,
		Endpoint:      "",
		ApiKey:        viper.GetString("apikey"),
		ApiSecretKey:  viper.GetString("secretkey"),
		ApiPassphrase: viper.GetString("password"),
		ClientId:      "",
		Lever:         0,
	}
	return apiConfig
}
