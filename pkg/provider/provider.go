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


var errEmptyStr = errors.New("it seems like someone give us empty string")
var installerPath = ""

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ocpbooster_cluster": resourceOcpCluster(),
		},
		Schema: map[string]*schema.Schema{
			"installer_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local filesystem path to install - openshift-install tool",
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
