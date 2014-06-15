package settings

import (
	"errors"
	"github.com/jinzhu/gorm"
	"paintToWin/storage"
	"reflect"
	"strconv"
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
			switch field.Kind() {
			case reflect.Bool:
				if value, err := strconv.ParseBool(settingValue); err != nil {
					errs = append(errs, "Cannot convert "+settingValue+" to bool for setting "+settingName)
				} else {
					field.SetBool(value)
				}
			case reflect.Int:
				if value, err := strconv.ParseInt(settingValue, 10, 32); err != nil {
					errs = append(errs, "Cannot convert "+settingValue+" to int for setting "+settingName)
				} else {
					field.SetInt(value)
				}
			case reflect.Int64:
				if value, err := strconv.ParseInt(settingValue, 10, 64); err != nil {
					errs = append(errs, "Cannot convert "+settingValue+" to int64 for setting "+settingName)
				} else {
					field.SetInt(value)
				}
			case reflect.String:
				field.SetString(settingValue)
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
		return nil
	}
}

func loadSettings(serverName string, db gorm.DB) map[string]string {
	serverSettings := []storage.Setting{}
	globalSettings := []storage.Setting{}
	db.Where(&storage.Setting{Server: serverName}).Find(&serverSettings)
	db.Where(&storage.Setting{Server: "global"}).Find(&globalSettings)

	settings := make(map[string]string)

	for _, setting := range serverSettings {
		settings[setting.Key] = setting.Value
	}

	for _, setting := range globalSettings {
		if _, ok := settings[setting.Key]; !ok {
			settings[setting.Key] = setting.Value
		}
	}

	return settings
}
