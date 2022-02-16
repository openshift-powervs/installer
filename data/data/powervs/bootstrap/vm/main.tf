data "ignition_config" "bootstrap" {
  merge {
    source = var.ignition_url
  }
}

# Create the bootstrap instance
resource "ibm_pi_instance" "bootstrap" {
  pi_memory            = var.memory
  pi_processors        = var.processors
  pi_instance_name     = "${var.cluster_id}-bootstrap"
  pi_proc_type         = var.proc_type
  pi_image_id          = var.boot_image
  pi_sys_type          = var.sys_type
  pi_cloud_instance_id = var.cloud_instance_id
  pi_network_ids       = [var.network_id]

  pi_user_data         = base64encode(data.ignition_config.bootstrap.rendered)
  pi_key_pair_name     = var.cluster_key_id
  pi_health_status     = "WARNING"
}

data "ibm_pi_instance_ip" "bootstrap_ip" {
  depends_on = [ibm_pi_instance.bootstrap]

  pi_instance_name     = ibm_pi_instance.bootstrap.pi_instance_name
  pi_network_name      = var.network_name
  pi_cloud_instance_id = var.cloud_instance_id
}

