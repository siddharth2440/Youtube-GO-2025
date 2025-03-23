package configs

import "github.com/spf13/viper"

type Env struct {
	MONGO_URL      string
	REDIS_URL      string
	JWT_SECRET_KEY string
}

func Config() (*Env, error) {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Env{
		MONGO_URL:      viper.GetString("MONGO_URI"),
		REDIS_URL:      viper.GetString("REDIS_URI"),
		JWT_SECRET_KEY: viper.GetString("JWT_SECRET"),
	}, nil
}
