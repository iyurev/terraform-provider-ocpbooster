package booster

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/Masterminds/sprig"
	installertypes "github.com/openshift/installer/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"os/exec"
	"text/template"
)

var (
	//go:embed install-config.tmpl
	installConfigTmpl     string
	bootstrapIgnFileName  = "bootstrap.ign"
	masterIgnFileName     = "master.ign"
	workerIgnFileName     = "worker.ign"
	metadataFileName      = "metadata.json"
	authDirName           = "auth"
	kubeAdminPassFileName = "kubeadmin-password"
	kubeconfigFileName    = "kubeconfig"
)

func NewBooster(installerPath, clusterName, pubSSHkey, pullSecret, baseDomain, offMirror, trustBundle string) (*Booster, error) {
	booster := &Booster{}
	if installerPath == "" {
		return booster, errInstallerIsNotInstalled
	}
	if offMirror != "" {
		booster.OfflineMirror = offMirror
	}
	//TODO we need to refactoring decoding function from B64
	if trustBundle != "" {
		decoded, err := DecodeBytesB64(trustBundle)
		if err != nil {
			return nil, err
		}
		booster.AddTrustBundle = decoded
	}
	booster.InstallerPath = installerPath
	booster.ClusterName = clusterName
	booster.PubSshKey = pubSSHkey
	booster.PullSecret = pullSecret
	booster.BaseDomain = baseDomain
	booster.workDir = booster.setWorkDir(clusterName)
	if err := booster.initWorkDir(); err != nil {
		return nil, err
	}
	if err := booster.createIgnitions(); err != nil {
		return nil, err
	}
	if err := booster.fetchMetadata(); err != nil {
		return nil, err
	}
	if err := booster.fetchAuth(); err != nil {
		return nil, err
	}
	return booster, nil
}

func (b *Booster) Ignitions() *ignitions {
	return b.ignitions
}

func (b *Booster) ClusterMetadata() *metadata {
	return b.metadata
}

func (b *Booster) Auth() *auth {
	return b.auth
}

func (b *Booster) initWorkDir() error {

	if err := os.RemoveAll(b.workDir); err != nil {
		return err
	}
	if err := os.MkdirAll(b.workDir, 0777); err != nil {
		return err
	}
	instConf, err := b.genInstallConfig()
	if err != nil {
		return err
	}
	installConfigFile, err := os.Create(fmt.Sprintf("%s/%s", b.workDir, installConfigFileName))
	if err != nil {
		return err
	}
	if _, err := installConfigFile.Write(instConf); err != nil {
		return err
	}
	return nil
}

func (b *Booster) fetchMetadata() error {
	metadata := &metadata{}
	raw, err := readFile(fmt.Sprintf("%s/%s", b.workDir, metadataFileName))
	if err != nil {
		return err
	}
	metadata.rawMetadata = raw
	clusterMetadata := installertypes.ClusterMetadata{}
	if err := json.Unmarshal(raw, &clusterMetadata); err != nil {
		return err
	}
	metadata.clusterMetadata = clusterMetadata
	b.metadata = metadata
	return nil
}

func (b *Booster) fetchAuth() error {
	auth := &auth{}
	rawKubeAdminPass, err := readFile(fmt.Sprintf("%s/%s/%s", b.workDir, authDirName, kubeAdminPassFileName))
	if err != nil {
		return err
	}
	auth.kubeAdminPassword = fmt.Sprintf("%s", rawKubeAdminPass)
	rawKubeconfig, err := readFile(fmt.Sprintf("%s/%s/%s", b.workDir, authDirName, kubeconfigFileName))
	if err != nil {
		return err
	}
	auth.kubeconfig = fmt.Sprintf("%s", rawKubeconfig)
	b.auth = auth
	return nil
}

func (b *Booster) createIgnitions() error {
	ignitions := &ignitions{}
	openshiftInst := b.InstallerPath
	args := []string{"--dir", b.workDir, "create", "ignition-configs"}
	cmd := exec.Command(openshiftInst, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	if bsIgn, err := readFile(fmt.Sprintf("%s/%s", b.workDir, bootstrapIgnFileName)); err != nil {
		return err
	} else {
		ignitions.bootstrap = bsIgn
	}
	if masterIgn, err := readFile(fmt.Sprintf("%s/%s", b.workDir, masterIgnFileName)); err != nil {
		return err
	} else {
		ignitions.master = masterIgn
	}
	if workerIgn, err := readFile(fmt.Sprintf("%s/%s", b.workDir, workerIgnFileName)); err != nil {
		return err
	} else {
		ignitions.worker = workerIgn
	}
	b.ignitions = ignitions

	return nil
}

func (b *Booster) setWorkDir(dirName string) string {
	if home := os.Getenv("HOME"); home != "" {
		return fmt.Sprintf("%s/%s", home, dirName)
	}
	return fmt.Sprintf("%s/%s", "/tmp", dirName)
}

func (b *Booster) genInstallConfig() ([]byte, error) {
	templ, err := template.New("install-config").Funcs(sprig.TxtFuncMap()).Parse(installConfigTmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := templ.Execute(&buf, *b); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//func (b *Booster) generateInstallConfig() *installertypes.InstallConfig {
//	platform := installertypes.Platform{
//		None: &none.Platform{},
//	}
//	masters := &installertypes.MachinePool{
//		Name:     "master",
//		Replicas: &masterReplicas,
//	}
//	computes := []installertypes.MachinePool{
//		{
//			Name:     "worker",
//			Replicas: &workerReplicas,
//		},
//	}
//	networking := &installertypes.Networking{
//		ClusterNetwork: []installertypes.ClusterNetworkEntry{
//			{
//				CIDR: ipnet.IPNet{
//					IPNet: defaultIPNet(),
//				},
//				HostPrefix: 23,
//			},
//		},
//		NetworkType: "OpenShiftSDN",
//		ServiceNetwork: []ipnet.IPNet{
//			defaultServiceNetwork(),
//		},
//	}
//	installConfig := &installertypes.InstallConfig{
//		TypeMeta: metav1.TypeMeta{
//			APIVersion: installertypes.InstallConfigVersion,
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name: b.ClusterName,
//		},
//		ControlPlane: masters,
//		Compute:      computes,
//		Platform:     platform,
//		Networking:   networking,
//		SSHKey:       b.PubSshKey,
//		PullSecret:   b.PullSecret,
//		BaseDomain:   b.BaseDomain,
//	}
//	return installConfig
//}
