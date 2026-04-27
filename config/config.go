package config

import (
	"exchangeapp/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name     string
		Port     string
		Mode     string
		LogLevel string
		LogFile  string
	}
	Database struct {
		Dsn             string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime int
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
		PoolSize int
	}
	JWT struct {
		Secret          string
		ExpirationHours int
	}
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	setGinMode()
	initLogger()
	initDB()
	initRedis()
	utils.InitJWTConfig(AppConfig.JWT.Secret, AppConfig.JWT.ExpirationHours)
	log.Println("✅ 配置初始化完成")
}

func setGinMode() {
	mode := strings.ToLower(AppConfig.App.Mode)
	if mode == "release" || mode == "prod" {
		os.Setenv("GIN_MODE", "release")
	} else {
		os.Setenv("GIN_MODE", "debug")
	}
}

func initLogger() {
	logLevel := getLogLevel(AppConfig.App.LogLevel)
	logFile := AppConfig.App.LogFile

	if err := utils.InitLogger(logFile, logLevel); err != nil {
		log.Printf("Warning: Failed to initialize file logger, using stdout: %v", err)
		utils.InitLogger("", logLevel)
	}

	utils.Info("日志系统初始化完成，日志级别: %s", AppConfig.App.LogLevel)
}

func getLogLevel(level string) utils.LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return utils.DEBUG
	case "info":
		return utils.INFO
	case "warn":
		return utils.WARN
	case "error":
		return utils.ERROR
	case "fatal":
		return utils.FATAL
	default:
		return utils.INFO
	}
}
