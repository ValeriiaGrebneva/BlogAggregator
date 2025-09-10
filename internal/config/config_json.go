package config

type Config struct {
	Connection string `json:"db_url"`
	Username   string `json:"current_user_name"`
}
