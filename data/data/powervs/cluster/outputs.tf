output "bootstrap_ip" {
  value = module.loadbalancer.powervs_lb_hostname
}

output "control_plane_ips" {
  value = module.master.master_ips
}

output "cluster_key_id" {
  value = ibm_pi_key.cluster_key.key_id
}

output "bootstrap_ignition_host" {
  value = module.bootstrap.bootstrap_ignition_host
}

output "bootstrap_ignition_bucket" {
  value = module.bootstrap.bootstrap_ignition_bucket
}

output "bootstrap_ignition_key" {
  value = module.bootstrap.bootstrap_ignition_key
}

output "boot_image_id" {
  value = ibm_pi_image.boot_image.image_id
}

output "network_id" {
  value = module.master.network_id
}

output "load_balancer_id" {
  value = module.loadbalancer.load_balancer_id
}

output "load_balancer_int_id" {
  value = module.loadbalancer.load_balancer_int_id
}

output "machine_config_lb_pool" {
  value = module.loadbalancer.machine_config_lb_pool
}

output "api_lb_pool" {
  value = module.loadbalancer.api_lb_pool
}

output "api_int_lb_pool" {
  value = module.loadbalancer.api_int_lb_pool
}
