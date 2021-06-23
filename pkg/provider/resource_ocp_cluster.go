
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-modules-collection/ocpbooster/pkg/booster"
	"log"
)

func resourceOcpCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Initial configuration for installing a new OpenShift cluster",
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Optional:    false,
				Description: "Cluster name, '.metadata.name:' in the install-config.yaml",
			},
			"base_domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Optional:    false,
				Description: "Cluster base domain, '.baseDomain:' in the install-config.yaml",
			},
			"pub_ssh_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Optional:    false,
				Description: "Public SSH key, '.sshKey:' in the install-config.yaml",
			},
			"pull_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Optional:    false,
				Description: "Cluster pull secret, '.pullSecret: in the install-config.yaml",
			},
			"add_trust_ca_bundler": {
				Type:      schema.TypeString,
				Required:  false,
				Sensitive: false,
				ForceNew:  true,
				Optional:  true,
				Description: "Cluster additional CA trust bundle, '.additionalTrustBundle: in the install-config.yaml",

			},
			"offline_mirror": {
				Type:      schema.TypeString,
				Required:  false,
				Sensitive: false,
				ForceNew:  true,
				Optional:  true,
				Description: "Mirror for offline cluster, '.imageContentSources: in the install-config.yaml",
			},
			"bootstrap_ign": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: false,
				Description: "Bootstrap ignition configuration bundle, it's content of bootstrap.ign",
			},
			"cluster_ca": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: false,
				Description: "Self-signed cluster CA",
			},
			"master_ign": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: false,
				Description: "Master node ignition configuration bundle",
			},
			"worker_ign": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: false,
				Description: "Worker node ignition configuration bundle",
			},
			"kubeadmin_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
				Description: "Password for service user account - kubeadmin",
			},
			"kubeconfig": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
				Description: "kubeconfig with super admin credentials",
			},
			"cluster_id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: false,
				Description: "Cluster unique ID",
			},
		},
		CreateContext: resourceOcpClusterCreate,
		ReadContext:   resourceOcpClusterRead,
		Delete:        schema.RemoveFromState,
	}
}

func resourceOcpClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	clusterName, err := stringMake(d.Get("cluster_name"))
	if err != nil {
		return diag.FromErr(err)
	}
	baseDomain, err := stringMake(d.Get("base_domain"))
	if err != nil {
		return diag.FromErr(err)
	}

	pubSSHkey, err := stringMake(d.Get("pub_ssh_key"))
	if err != nil {
		return diag.FromErr(err)
	}
	pullSecret, err := stringMake(d.Get("pull_secret"))
	if err != nil {
		return diag.FromErr(err)
	}
	addTrustCABandle, err := stringMake(d.Get("add_trust_ca_bundler"))
	if err != nil {
		if err == errEmptyStr {
			log.Println("this is optional parameters")
		} else {
			return diag.FromErr(err)
		}
	}
	offlineMirror, err := stringMake(d.Get("offline_mirror"))
	if err != nil {
		if err == errEmptyStr {
			log.Println("this is optional parameters")
		} else {

			return diag.FromErr(err)
		}
	}
	bs, err := booster.NewBooster(installerPath, clusterName, pubSSHkey, pullSecret, baseDomain, offlineMirror, addTrustCABandle)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("bootstrap_ign", bs.Ignitions().BootstrapB64())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("master_ign", bs.Ignitions().MasterB64()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("worker_ign", bs.Ignitions().WorkerB64()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kubeadmin_password", bs.Auth().KubeAdminPassword()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kubeconfig", bs.Auth().Kubeconfig()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_id", bs.ClusterMetadata().ClusterID()); err != nil {
		return diag.FromErr(err)
	}
	clusterCA, err := bs.Ignitions().ClusterCA()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_ca", clusterCA); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(clusterName)
	return diags
}

func resourceOcpClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

