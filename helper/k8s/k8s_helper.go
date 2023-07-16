package k8s

type K8s interface {
	UpdateSecret(namespace, secretName, key, value string) error
	RestartWorkloads(namespace, kind, workload string) error
}
