package helper

import (
	"gopkg.in/yaml.v3"
	"os"
	"zeus/log"
	"zeus/model"
)

func LoadConfigFile(cfg *model.RotationKey) {
	f, err := os.ReadFile("config_file/schedules.yaml")
	if err != nil {
		log.Error(err)
	}

	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		log.Error(err)
	}
}
