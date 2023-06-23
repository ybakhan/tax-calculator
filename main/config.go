package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents configurable properties of tax-calculator
type Config struct {
	Port            int `yaml:"port"`
	InterviewServer struct {
		BaseURL string `yaml:"baseUrl"`
	} `yaml:"interviewServer"`

	HTTPClient struct {
		Retry struct {
			Max  int `yaml:"max"`
			Wait struct {
				MinSeconds int `yaml:"minSeconds"`
				MaxSeconds int `yaml:"maxSeconds"`
			} `yaml:"wait"`
		} `yaml:"retry"`
	} `yaml:"httpClient"`
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
