package config

import (
    "fmt"
    "github.com/spf13/viper"
)

type Config struct {
    Server struct {
        Port int    `mapstructure:"port"`
        Env  string `mapstructure:"env"`
    } `mapstructure:"server"`

    Database struct {
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        Name     string `mapstructure:"name"`
        User     string `mapstructure:"user"`
        Password string `mapstructure:"password"`
    } `mapstructure:"database"`

    Redis struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"redis"`

    JWT struct {
        Secret    string `mapstructure:"secret"`
        ExpiresIn string `mapstructure:"expires_in"`
    } `mapstructure:"jwt"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("error reading config file: %w", err)
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("error unmarshaling config: %w", err)
    }

    return &config, nil
}