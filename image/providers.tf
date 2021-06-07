terraform {
  required_providers {
    vra = {
      source = "vmware/vra"
      version = "0.3.4"
    }
    ct = {
      source  = "poseidon/ct"
      version = "0.7.1"
    }
    matchbox = {
      source  = "poseidon/matchbox"
      version = "0.4.1"
    }
    aws = {
      source = "hashicorp/aws"
      version = "3.42.0"
    }
    template = {
      source = "hashicorp/template"
      version = "2.2.0"
    }
  }
}

