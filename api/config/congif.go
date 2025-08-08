package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort         string `mapstructure:"SERVER_PORT"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPassword         string `mapstructure:"DB_PASSWORD"`
	DBName             string `mapstructure:"DB_NAME"`
	DBPort             string `mapstructure:"DB_PORT"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	EncryptionKey      string `mapstructure:"ENCRYPTION_KEY"`
	TripayApiUrl       string `mapstructure:"TRIPAY_API_URL"`
	TripayApiKey       string `mapstructure:"TRIPAY_API_KEY"`
	TripayPrivateKey   string `mapstructure:"TRIPAY_PRIVATE_KEY"`
	TripayMerchantCode string `mapstructure:"TRIPAY_MERCHANT_CODE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
