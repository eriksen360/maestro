package lib

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// TODO: make as interface
func GetConfigurationFileSettings() map[string]interface{} {

	viper.SetConfigName("jobs")
	viper.SetConfigType("yaml")

	pwd, _ := os.Getwd()
	viper.AddConfigPath(pwd + "/conf")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return viper.AllSettings()
}

// TODO: make as interface
func Validate(settings map[string]interface{}) (bool, error) {

	return true, nil
}

// FormatConfigurationFileSettings This function translates the map of settings into a
// more explicit list, which can be further formatted to provide a concise instruction set
// for the executor algorithm
func FormatConfigurationFileSettings(f string) {

	settings := GetConfigurationFileSettings()
	ok, err := Validate(settings)
	if !ok {
		panic(fmt.Errorf("invalid configuration file: %w", err))
	}

}

func CreateInstructionList() {

	/* Uses the formatted conf file to create a list of instructions specific to a host, and passes
	   them to a context where the orchestrator can execute them */

}
