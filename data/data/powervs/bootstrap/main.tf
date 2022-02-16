provider "ibm" {
  alias            = "vpc"
  ibmcloud_api_key = var.powervs_api_key
  region           = var.powervs_vpc_region
  zone             = var.powervs_vpc_zone
}

provider "ibm" {
  alias = "powervs"
  ibmcloud_api_key = var.powervs_api_key
  region           = var.powervs_region
}

module "vm" {
  providers = {
    ibm = ibm.powervs
  }
  source            = "./vm"

  cloud_instance_id = var.powervs_cloud_instance_id
  cluster_id        = var.cluster_id
  cluster_key_id    = var.cluster_key_id
  ignition_url      = var.bootstrap_ignition_url
  boot_image        = var.boot_image_id

  memory            = var.powervs_bootstrap_memory
  processors        = var.powervs_bootstrap_processors
  sys_type          = var.powervs_sys_type
  proc_type         = var.powervs_proc_type
  network_id        = var.network_id
  network_name      = var.powervs_network_name
}

module "lb" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./lb"

  load_balancer_id     = var.load_balancer_id
  bootstrap_private_ip = module.vm.bootstrap_private_ip
}
