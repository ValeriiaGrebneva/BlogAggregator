package config

import "os"

type Config struct {
	Connection string `json:"db_url"`
	Username   string `json:"current_user_name"`
}

func Read(path string) (Config, error) {
	link := os.UserHomeDir + path

	data, err := os.ReadFile(link)
	if err != nil {
        return nil, err
	}

	var configStruct Config
    decoder := json.NewDecoder(data)
    if err := decoder.Decode(&configStruct); err != nil {
        return nil, err
    }

    return configStruct, nil
}