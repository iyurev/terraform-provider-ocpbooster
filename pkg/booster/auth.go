package booster

func (a *auth) KubeAdminPassword() string {
	return a.kubeAdminPassword
}

func (a *auth) Kubeconfig() string {
	return a.kubeconfig
}
