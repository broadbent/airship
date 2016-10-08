package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	MongoURI          string
	Interval          string
	BidIncrement      int
	MemorySplit       int
	FinalStage        int
	ProvisionerPath   string
	DatabaseName      string
	StartingValuation int
}

func Read(path string) Configuration {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	if e := decoder.Decode(&configuration); e != nil {
		log.Panic(e)
	}
	return configuration
}
