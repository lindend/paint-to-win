package settings

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"

	"paintToWin/storage"
	"paintToWin/util"
)

func Load(serverName string, db gorm.DB, settings interface{}) error {
	settingsMap := loadSettings(serverName, db)

	errs := []string{}

	settingsValue := reflect.Indirect(reflect.ValueOf(settings))
	settingsType := settingsValue.Type()
	for i := 0; i < settingsValue.NumField(); i++ {
		field := settingsValue.Field(i)
		settingName := settingsType.Field(i).Name

		if settingValue, ok := settingsMap[settingName]; ok {
			err := util.ParseStringToValue(settingValue, field)
			if err != nil {
				errs = append(errs, err.Error())
			}
		} else {
			errs = append(errs, "No setting value registered for "+settingName)
		}
	}

	if len(errs) > 0 {
		errString := ""
		for _, err := range errs {
			errString += err + "\n"
		}
		return errors.New(errString)
	} else {
		fmt.Printf("Using settings %+v \n", settings)
		return nil
	}
}

func loadSettings(serverName string, db gorm.DB) map[string]string {
	serverSettings := []storage.Setting{}
	globalSettings := []storage.Setting{}

	db.Where(&storage.Setting{Server: serverName}).Find(&serverSettings)
	db.Where(&storage.Setting{Server: "global"}).Find(&globalSettings)

	settings := make(map[string]string)

	argSettings := argumentSettings()
	envSettings := environmentSettings()

	for _, setting := range globalSettings {
		settings[setting.Key] = setting.Value
	}

	for _, setting := range serverSettings {
		settings[setting.Key] = setting.Value
	}

	for settingKey, settingValue := range envSettings {
		settings[settingKey] = settingValue
	}

	for settingKey, settingValue := range argSettings {
		settings[settingKey] = settingValue
	}

	return settings
}

func environmentSettings() map[string]string {
	return parseKeyValueSettings(os.Environ())
}

func argumentSettings() map[string]string {
	return parseKeyValueSettings(os.Args[1:])
}

func parseKeyValueSettings(settings []string) map[string]string {
	result := make(map[string]string)
	for _, settingValue := range settings {
		keyValue := strings.SplitN(settingValue, "=", 2)
		if len(keyValue) == 2 && len(keyValue[0]) > 0 {
			result[keyValue[0]] = keyValue[1]
		}
	}
	return result
}
