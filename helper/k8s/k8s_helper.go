package k8s

type K8s interface {
	UpdateSecret(namespace, secretName, key, value string) error
}
