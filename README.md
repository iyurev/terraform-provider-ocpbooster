## OCP booster terraform provider 

#### Fast terraform way to init metadata for OpenShift cluster installation.


It's simple terraform provider can help you  with RH OpenShift UPI cluster installation    


If we through up the official UPI installation guide, then your find out that you need some tied manual actions with an installation tool - openshift-install

1. Download openshift-install itself from https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/
2. Prepare a correct install-config.yaml  in the installation directory
2. Run ```openshift-install create ignition-config``` then  grub ignition configs and  admin cluster credentials to use it later.

All these actions can be fulfilled  with this terraform module, e.g.

```bigquery
resource "ocpbooster_cluster" "my-cluster" {
  cluster_name = "my-cluster"
  base_domain = "ocp.local"
  pull_secret = var.pullSecret
  pub_ssh_key = var.pubKey
  offline_mirror = var.offlineMirror
  add_trust_ca_bundler = var.addTrustCA
}
```

After then resource ```my-cluster``` will be created we can get the all information described above from its attributes (ignition configs, cluster credentials etc.):

Example 1,  uploading bootstrap ignition to S3 bucket:
```bigquery
resource "aws_s3_bucket_object" "bootstrap-ign" {
  bucket = var.minioOpenshiftPubBucket
  key = local.bootstrapIgnName
  content_base64 = ocpbooster_cluster.my-cluster.bootstrap_ign
}
```

Example 2, just send bootstrap ignition and a new cluster CA bundle to the next module through output:
```bigquery
output "bootstrap_ign" {
  value = ocpbooster_cluster.my-cluster.bootstrap_ign
}

output "cluster_ca" {
  value = ocpbooster_cluster.my-cluster.cluster_ca
}
```