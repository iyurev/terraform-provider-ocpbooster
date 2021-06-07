package booster

import (
	"errors"
	"github.com/openshift/installer/pkg/ipnet"
	installertypes "github.com/openshift/installer/pkg/types"
	"net"
)

const installConfigFileName = "install-config.yaml"

var (
	masterReplicas = int64(3)
	workerReplicas = int64(0)

	errInstallerIsNotInstalled = errors.New("openshift-installer is not installed")
)

type ignitions struct {
	master    []byte
	worker    []byte
	bootstrap []byte
}

type metadata struct {
	rawMetadata     []byte
	clusterMetadata installertypes.ClusterMetadata
}

type auth struct {
	kubeAdminPassword string
	kubeconfig        string
}

type Booster struct {
	InstallerPath  string
	ClusterName    string
	PubSshKey      string
	PullSecret     string
	OfflineMirror  string
	AddTrustBundle string

	BaseDomain string
	workDir    string
	ignitions  *ignitions
	auth       *auth
	metadata   *metadata
}

func defaultIPNet() net.IPNet {
	_, ipnetwork, _ := net.ParseCIDR("10.128.0.0/14")
	return *ipnetwork
}
func defaultServiceNetwork() ipnet.IPNet {
	_, ipnetwork, _ := net.ParseCIDR("172.30.0.0/16")
	return ipnet.IPNet{
		IPNet: *ipnetwork,
	}
}
