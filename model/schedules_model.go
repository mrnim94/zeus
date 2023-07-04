package model

type RotationKey struct {
	Tasks []struct {
		Name           string `yaml:"name"`
		Cron           string `yaml:"cron"`
		UsernameOnAws  string `yaml:"usernameOnAws"`
		AccessKeyOnK8S struct {
			Name string `yaml:"name"`
			Key  string `yaml:"key"`
		} `yaml:"accessKeyOnK8s"`
		SecretKeyOnK8S struct {
			Name string `yaml:"name"`
			Key  string `yaml:"key"`
		} `yaml:"secretKeyOnK8s"`
	} `yaml:"schedules"`
}
