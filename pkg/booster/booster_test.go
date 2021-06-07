package booster

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	boosterTestsClusterName = "ocp-test"
	boosterTestsSSHKey      = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCx6eDY9OekCJiJJrW5a5zwExDT1aEhBIM2aF8To2WMoDY4HI9ScWyP72cewQ6sdVNUcS0RtuCI3bvUz2wwEKgvxBK6B19Shv4/DTBKQT5vVDoyiuvxKaxaNKUpmpRRfoxMCrQjTArqf08nBlRK+dZACa/GYpOQGTBcH6iqxMAtGQPrqUUMRjyONqqfxwUgs39X0TiQLX2UbeSgeyPb/PmBxpf3LRRJG+2KN+jdBhYntjg8i3PNNh5MiTy2Ya3xBqBmSiqL9pELrnPocbEjDgm59oC0Ll4CWMFpv2EnEAJlNC+egqxfUTZ46rODZPJ+WBAmg332wKpqaLppzYzSA4TF ocp-test@local"
	boosterTestsBaseDomain  = "test.local"
	//boosterTestsOfflineMirr = "registry.local/openshift-release-dev/ocp-release"
	boosterTestsOfflineMirr = ""
	boosterTestsAddCaBundle = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM3VENDQWRXZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFtTVNRd0lnWURWUVFEREJ0cGJtZHkKWlhOekxXOXdaWEpoZEc5eVFERTFPVGs1T1RFd056Z3dIaGNOTWpBd09URXpNRGsxT0RBeFdoY05Nakl3T1RFegpNRGsxT0RBeVdqQW1NU1F3SWdZRFZRUUREQnRwYm1keVpYTnpMVzl3WlhKaGRHOXlRREUxT1RrNU9URXdOemd3CmdnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUURVU0owbnFnUml4bVZXVE9wQnZHMzIKbE0wMVBoMjkrQzF1bExQazRRYTJIUmUvTVJlcURtdzVYclhQZCtRT2xzd0wrNmgybDAycDRqNWhBakFuZnBUUwpJcUpybzBTS2pwRXpGREJyUDhka09OS3RtT3N6UVJqNEJhbEVVbEx1UFFLUms1WjlxZWdhay9QTjQ3YU1yMmJ2CkZRYVNsSkVnL3hQNjdEb0M3SlZKanY0YjlMTVZoQld3NWZhUHorbGVOdHVsRE1ZQzNyOXVVS1BPZ1FPakJ1WEgKYnhjVUE5VldiS3hjamJCU0hSKzJtNmtCdG5ZZ2tER21VZWIrR204clh6R2tXdkdqSGg2T1lIV0Y0MzNRWHlhSgpKM1FtWkQ5TitQdWNleVM2UGM5cGRUMXpITGRIMWhMQWRVdEtrN3A3RTFGeVBkVG1UdHJ1amJZRWMxNkxObjVuCkFnTUJBQUdqSmpBa01BNEdBMVVkRHdFQi93UUVBd0lDcERBU0JnTlZIUk1CQWY4RUNEQUdBUUgvQWdFQU1BMEcKQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUNlV3hESFlkWXB2eC9YMjIvS3V6cngwakU3TmlBaGsxMDF4Nlp0bmxyMQo1d2w1dG81SE1TUitqTzdUeVZ2L1JNMklFMXZNT2ZUUVFZdmVKWXB6NVdKKzYyemZaZkJkSEJ0dHhkT2QvTkRJCjRTOE5GS3ZmanVOT2pldCtCeXgrVDRLSzVjWGVwL3Z2RnNnMUxDZXRPQ0RiVHBxODdHQXJvRGVDMVd4VENpaEMKZ1NTQTQ5bnZqendIeGlJYTBPWExQOG9XbkkveDRiRlVDaWpQWkNrKzRhSHdkQ0RzRmVJVXdiTEhsRmVydjNaVgpVSG80VEhMMzYzYUhVUlk5YjRIb0xUOWxOdVJoNkRGb0cwL0VLUnFYMTZZL1RYK2d3SlBkYmY4ZG5pWUlwY3lWCjZ3ZWxMOUM0OVpvcDlHNFNpVVU0NFMwRHorWGtROHFjSDdRekkyQnBwSmtuCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
)

func boosterTestsInstallerPath() string {
	home := os.Getenv("HOME")
	if home != "" {
		return fmt.Sprintf("%s/%s", home, "openshift-install")
	} else {
		return fmt.Sprintf("%s/%s", "/tmp", "openshift-install")
	}
}

func boosterTestsPullSecret() string {
	pullSecret := os.Getenv("OCP_TEST_PULL_SECRET")
	if pullSecret == "" {
		log.Fatal("empty pull secret, you must define OCP_TEST_PULL_SECRET variable")
	}
	return pullSecret
}

func TestNewBooster(t *testing.T) {
	t.Log("Init new booster")
	booster, err := NewBooster(boosterTestsInstallerPath(), boosterTestsClusterName, boosterTestsSSHKey, boosterTestsPullSecret(), boosterTestsBaseDomain, boosterTestsOfflineMirr, boosterTestsAddCaBundle)
	if err != nil {
		t.Fatal(err)
	}
	ca, err := booster.Ignitions().ClusterCA()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Cluster CA: %s\n", ca)
	t.Logf("Cluster meta: %s\n", booster.ClusterMetadata().RawMetadata())
}
