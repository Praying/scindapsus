package config

import (
	"fmt"
	"testing"
)

func TestConfigEngine_ReadConfig(t *testing.T) {
	apiConfig := GetConfigEngine().ReadConfig()

	fmt.Println(apiConfig.ApiKey)
	fmt.Println(apiConfig.ApiSecretKey)
	fmt.Println(apiConfig.ApiPassphrase)
}
