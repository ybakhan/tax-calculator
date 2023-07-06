package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents configurable properties of tax-calculator
type Config struct {
	Port     int            `yaml:"port"`
	ApiToken ApiTokenConfig `yaml:"apiToken"`

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"redis"`

	InterviewServer struct {
		BaseURL string `yaml:"baseUrl"`
	} `yaml:"interviewServer"`

	HTTPClient struct {
		TimeoutMs int `yaml:"timeoutMs"`

		Retry struct {
			Max int `yaml:"max"`

			Wait struct {
				MinMs int `yaml:"minMs"`
				MaxMs int `yaml:"maxMs"`
			} `yaml:"wait"`
		} `yaml:"retry"`
	} `yaml:"httpClient"`
}

type ApiTokenConfig struct {
	ExpirationMinutes int    `yaml:"expirationMinutes"`
	Secret            string `yaml:"secret"`
}

func readConfig() *Config {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
