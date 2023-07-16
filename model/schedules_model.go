package model

type RotationKey struct {
	Schedules []Schedule `yaml:"schedules"`
}

type RestartWorkloads struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}

type Schedule struct {
	Name           string `yaml:"name"`
	Cron           string `yaml:"cron"`
	UsernameOnAws  string `yaml:"usernameOnAws"`
	NamespaceOnK8s string `yaml:"namespaceOnK8s"`
	AccessKeyOnK8S struct {
		Name string `yaml:"name"`
		Key  string `yaml:"key"`
	} `yaml:"accessKeyOnK8s"`
	SecretKeyOnK8S struct {
		Name string `yaml:"name"`
		Key  string `yaml:"key"`
	} `yaml:"secretKeyOnK8s"`
	RestartWorkloads []RestartWorkloads `yaml:"restartWorkloads"`
}
