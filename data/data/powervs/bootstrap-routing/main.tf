locals {
  api_servers       = var.powervs_expose_bootstrap ? concat(var.control_plane_ips, [var.bootstrap_private_ip]): var.control_plane_ips
  tcp_ports         = concat([22623, 6443], (var.powervs_expose_bootstrap ? [22]: []))
}

provider "ibm" {
  ibmcloud_api_key = var.powervs_api_key
  region           = var.powervs_vpc_region
  zone             = var.powervs_vpc_zone
}

resource "ibm_is_lb_pool_member" "machine_config_member" {
  count      = length(local.api_servers)

  lb             = var.load_balancer_int_id
  pool           = var.machine_config_lb_pool
  port           = 22623
  target_address = local.api_servers[count.index]
}

resource "ibm_is_lb_pool_member" "api_member_int" {
  depends_on = [ibm_is_lb_pool_member.machine_config_member]
  count      = length(local.api_servers)

  lb             = var.load_balancer_int_id
  pool           = var.api_int_lb_pool
  port           = 6443
  target_address = local.api_servers[count.index]
}

resource "ibm_is_lb_pool_member" "api_member" {
  depends_on = [ibm_is_lb_pool_member.machine_config_member]
  count      = length(local.api_servers)

  lb             = var.load_balancer_id
  pool           = var.api_lb_pool
  port           = 6443
  target_address = local.api_servers[count.index]
}

resource "ibm_is_security_group_rule" "inbound_ports" {
  count     = length(local.tcp_ports)
  group     = var.security_group_id
  direction = "inbound"
  tcp {
    port_min = local.tcp_ports[count.index]
    port_max = local.tcp_ports[count.index]
  }
}

