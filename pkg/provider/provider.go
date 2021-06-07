package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-modules-collection/ocpbooster/pkg/booster"
	"github.com/terraform-modules-collection/ocpbooster/pkg/config"
)

type configureVal struct {
	pubKey     string
	pullSecret string
}

var errEmptyPubKey = errors.New("empty SSH public key")
var errEmptyPullSecret = errors.New("empty image pull secret")
var errEmptyStr = errors.New("it seems like someone give us empty string")

var installerPath = ""

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ocpbooster_cluster": resourceOcpCluster(),
		},
		//DataSourcesMap: map[string]*schema.Resource{},
		Schema: map[string]*schema.Schema{
			"installer_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "openshift-install path",
				Sensitive:   false,
				Required:    false,
			},
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	_ = ctx
	var diags diag.Diagnostics

	if ip := d.Get("installer_path").(string); ip != "" {
		installerPath = ip
	} else {
		installerPath = config.InstallerPath
	}
	err := booster.ExtractInstaller(installerPath)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return nil, diags
}

func stringMake(s interface{}) (string, error) {
	str := fmt.Sprintf("%s", s)
	if str == "" {
		return str, errEmptyStr
	}
	return str, nil
}
