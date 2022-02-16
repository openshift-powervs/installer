variable "load_balancer_id" { type = string }
variable "load_balancer_int_id" { type = string }
variable "machine_config_lb_pool" { type = string }
variable "api_lb_pool" { type = string }
variable "api_int_lb_pool" { type = string }
variable "bootstrap_private_ip" { type = string }
variable "control_plane_ips" { type = list }
variable "security_group_id" { type = string }
