data "ibm_iam_auth_token" "iam_token" {}

# Create the bootstrap instance
resource "ibm_pi_instance" "bootstrap" {
  pi_memory            = var.memory
  pi_processors        = var.processors
  pi_instance_name     = "${var.cluster_id}-bootstrap"
  pi_proc_type         = var.proc_type
  pi_image_id          = var.boot_image
  pi_sys_type          = var.sys_type
  pi_cloud_instance_id = var.cloud_instance_id
  pi_network {
    network_id = var.network_id
  }
  pi_user_data         = base64encode(templatefile("${path.module}/templates/bootstrap.ign", {
    HOSTNAME    = var.ignition_host
    BUCKET_NAME = var.ignition_bucket
    OBJECT_NAME = var.ignition_key
    IAM_TOKEN   = data.ibm_iam_auth_token.iam_token.iam_access_token
  }))
  pi_key_pair_name     = var.cluster_key_id
  pi_health_status     = "WARNING"
}

data "ibm_pi_instance_ip" "bootstrap_ip" {
  depends_on = [ibm_pi_instance.bootstrap]

  pi_instance_name     = ibm_pi_instance.bootstrap.pi_instance_name
  pi_network_name      = var.network_name
  pi_cloud_instance_id = var.cloud_instance_id
}

