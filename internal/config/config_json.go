package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Connection string `json:"db_url"`
	Username   string `json:"current_user_name"`
}

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	link := homeDir + string(configFileName)

	data, err := os.ReadFile(link)
	if err != nil {
        return Config{}, err
	}

	var configStruct Config

	if err := json.Unmarshal(data, &configStruct); err != nil {
        return Config{}, err
    }
    //decoder := json.NewDecoder(strings.NewReader(data))
    //if err := decoder.Decode(&configStruct); err != nil {
    //    return Config{}, err
    //}

    return configStruct, nil
}

func (configStruct Config) SetUser(current_user_name string) error {
	configStruct.Username = current_user_name
	jsonStruct, err := json.Marshal(&configStruct)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFileName, jsonStruct, 0666) 
	if err != nil {
		return err
	}

	return nil
}