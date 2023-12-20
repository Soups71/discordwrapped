package main

import (
	"discordwrapped/pkg/bot"
	"discordwrapped/pkg/config"
	"log"
	"os"
)

func createFileIfNotExists(filename string) (*os.File, error) {
	// Check if the file exists
	_, err := os.Stat(filename)

	// If the file doesn't exist, create it
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	// If the file already exists, open it
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: go run main.go <filename>")
	}
	file, err := createFileIfNotExists(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	err = config.ReadConfig()

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
