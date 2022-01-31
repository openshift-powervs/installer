terraform {
  required_version = ">= 0.14"
  required_providers {
    ibm = {
      source = "openshift/local/ibm"
    }
    ibms3presign = {
      source = "openshift/local/ibms3presign"
    }
    ignition = {
      source = "openshift/local/ignition"
    }
  }
}
