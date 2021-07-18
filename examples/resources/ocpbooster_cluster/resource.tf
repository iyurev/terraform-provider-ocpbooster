resource "ocpbooster_cluster" "example" {
  cluster_name = "example"
  base_domain = "ocp.local"
  pull_secret = base64decode(var.pullSecret)
  pub_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDBPhWdcNUiQWIkikLU1DkL...."
}