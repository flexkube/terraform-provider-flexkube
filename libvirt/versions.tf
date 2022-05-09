terraform {
  required_version = ">= 0.14"

  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "3.1.1"
    }
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = "0.6.14"
    }
    ct = {
      source  = "poseidon/ct"
      version = "0.6.1"
    }
  }
}
