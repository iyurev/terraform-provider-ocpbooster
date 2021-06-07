package booster

func (m *metadata) RawMetadata() []byte {
	return m.rawMetadata
}

func (m *metadata) ClusterID() string {
	return m.clusterMetadata.ClusterID
}
func (m *metadata) InfraID() string {
	return m.clusterMetadata.InfraID
}
func (m *metadata) ClusterName() string {
	return m.clusterMetadata.ClusterName
}
