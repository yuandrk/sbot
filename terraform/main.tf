provider "flux" {
  kubernetes = {
    config_path = "${path.cwd}/.terraform/modules/gke_cluster/kubeconfig"
  }
  git = {
    url  = "https://github.com/${var.GITHUB_OWNER}/${var.FLUX_GITHUB_REPO}.git"
    http = {
      username    = "git"
      password    =  var.GITHUB_TOKEN
    }
  }
}

resource "flux_bootstrap_git" "this" {
  path = "./clusters"
  depends_on = [ module.gke_cluster ]
}

module "github_repository" {
  source                   = "github.com/den-vasyliev/tf-github-repository"
  github_owner             = var.GITHUB_OWNER
  github_token             = var.GITHUB_TOKEN
  repository_name          = var.FLUX_GITHUB_REPO
  public_key_openssh       = module.tls_private_key.public_key_openssh
  public_key_openssh_title = "flux"
}

module "tls_private_key" {
  source = "github.com/den-vasyliev/tf-hashicorp-tls-keys"
}

module "gke_cluster" {
  source         = "github.com/yuandrk/tf-google-gke-cluster"
  GOOGLE_REGION  = var.GOOGLE_REGION
  GOOGLE_PROJECT = var.GOOGLE_PROJECT
  GKE_CLUSTER_NAME = var.GKE_CLUSTER_NAME
  GKE_MACHINE_TYPE = var.GKE_MACHINE_TYPE
  GKE_NUM_NODES  = 1
}

terraform {
  backend "gcs" {
    bucket = "yuandrk-terraform-state-bucket"
    prefix = "terraform/state"
  }
  required_providers {
    flux = {
      source = "fluxcd/flux"
      version = "1.2.1"
    }
  }
}