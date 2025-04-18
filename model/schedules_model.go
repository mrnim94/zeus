package model

type RotationKey struct {
	Schedules []Schedule `yaml:"schedules"`
}

type RestartWorkloads struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
	NamespaceOnK8s string `yaml:"namespaceOnK8s"`
}

type Locations struct {
	SecretName      string `yaml:"secretName"`
	Style           string `yaml:"style"`
	CredentialOnK8S string `yaml:"credentialOnK8s,omitempty"`
	Profile         string `yaml:"profile,omitempty"`
	AccessKeyOnK8S  string `yaml:"accessKeyOnK8s,omitempty"`
	SecretKeyOnK8S  string `yaml:"secretKeyOnK8s,omitempty"`
	NamespaceOnK8s  string `yaml:"namespaceOnK8s"`
}

type Schedule struct {
	Name             string             `yaml:"name"`
	Cron             string             `yaml:"cron"`
	UsernameOnAws    string             `yaml:"usernameOnAws"`
	Action           string             `yaml:"action"`
	Locations        []Locations        `yaml:"locations"`
	RestartWorkloads []RestartWorkloads `yaml:"restartWorkloads"`
}
