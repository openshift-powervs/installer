output "powervs_lb_hostname" {
  value = ibm_is_lb.load_balancer.hostname
}

output "powervs_lb_int_hostname" {
  value = ibm_is_lb.load_balancer_int.hostname
}

output "load_balancer_id" {
  value = ibm_is_lb.load_balancer.id
}

output "load_balancer_int_id" {
  value = ibm_is_lb.load_balancer_int.id
}

output "machine_config_lb_pool" {
  value = ibm_is_lb_pool.machine_config_pool.id
}

output "api_lb_pool" {
  value = ibm_is_lb_pool.api_pool.id
}

output "api_int_lb_pool" {
  value = ibm_is_lb_pool.api_pool_int.id
}

output "security_group_id" {
  value = ibm_is_security_group.ocp_security_group.id
}
