package booster

import (
	"errors"
	"fmt"
	ign "github.com/coreos/ignition/v2/config/v3_2"
)

var (
	ErrWrongSecurityTLS = errors.New("wrong security CA TLS in the master.ign")
)

func (i *ignitions) Bootstrap() string {
	return fmt.Sprintf("%s", i.bootstrap)
}

func (i *ignitions) Master() string {
	return fmt.Sprintf("%s", i.master)
}

func (i *ignitions) Worker() string {
	return fmt.Sprintf("%s", i.worker)
}

func (i *ignitions) BootstrapB64() string {
	return toB64(i.bootstrap)
}

func (i *ignitions) MasterB64() string {
	return toB64(i.master)
}

func (i *ignitions) WorkerB64() string {
	return toB64(i.worker)
}

func (i *ignitions) ClusterCA() (string, error) {
	conf, report, err := ign.Parse(i.master)
	if err != nil {
		return "", err
	}
	_ = report
	ca := conf.Ignition.Security.TLS.CertificateAuthorities
	if len(ca) == 1 {
		return *ca[0].Source, nil
	}
	return "", ErrWrongSecurityTLS
}
