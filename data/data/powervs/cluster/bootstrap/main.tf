# TODO(mjturek): network and image data blocks can be in main module
#                as master and bootstrap will be using the same
#                network and image. Once we add in master module, make
#                the move.
data "ibm_pi_network" "network" {
  pi_network_name      = var.network_name
  pi_cloud_instance_id = var.cloud_instance_id
}

data "ibm_resource_group" "cos_group" {
  name = var.resource_group
}

resource "ibm_resource_instance" "cos_instance" {
  name              = "${var.cluster_id}-cos"
  resource_group_id = data.ibm_resource_group.cos_group.id
  service           = "cloud-object-storage"
  plan              = "standard"
  location          = var.cos_instance_location
  tags              = [var.cluster_id]
}

# Create an IBM COS Bucket to store ignition
resource "ibm_cos_bucket" "ignition" {
  bucket_name          = "${var.cluster_id}-bootstrap-ign"
  resource_instance_id = ibm_resource_instance.cos_instance.id
  region_location      = var.cos_bucket_location
  storage_class        = var.cos_storage_class
}

resource "ibm_resource_key" "cos_service_cred" {
  name                 = "${var.cluster_id}-cred"
  role                 = "Reader"
  resource_instance_id = ibm_resource_instance.cos_instance.id
  parameters           = { HMAC = true }
}

resource "ibms3presign" "bootstrap_ignition" {
  access_key_id     = ibm_resource_key.cos_service_cred.credentials["cos_hmac_keys.access_key_id"]
  secret_access_key = ibm_resource_key.cos_service_cred.credentials["cos_hmac_keys.secret_access_key"]
  bucket_name       = "${var.cluster_id}-bootstrap-ign"
  key               = "bootstrap.ign"
  region_location   = ibm_cos_bucket.ignition.region_location
  storage_class     = ibm_cos_bucket.ignition.storage_class
}

# Place the bootstrap ignition file in the ignition COS bucket
resource "ibm_cos_bucket_object" "ignition" {
  bucket_crn      = ibm_cos_bucket.ignition.crn
  bucket_location = ibm_cos_bucket.ignition.region_location
  content         = var.ignition
  key             = "bootstrap.ign"
}
