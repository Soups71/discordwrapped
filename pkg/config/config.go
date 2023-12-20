package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	Token     string
	BotPrefix string
	DBConn    string
	config    *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	DBConn    string `json:"DBConn"`
}

func ReadConfig() error {

	file, err := ioutil.ReadFile("config.json")

	if err != nil {
		log.Fatalln(err.Error())
		return err
	} else {
		log.Println("Successfully read config file")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	} else {
		log.Println("Successfully parsed config file")

	}
	Token = config.Token
	BotPrefix = config.BotPrefix
	DBConn = config.DBConn
	return nil

}
