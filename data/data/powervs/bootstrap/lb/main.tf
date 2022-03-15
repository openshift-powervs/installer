# bootstrap listener and backend pool
resource "ibm_is_lb_listener" "bootstrap_listener" {
  lb             = var.load_balancer_id
  port           = 22
  protocol       = "tcp"
  default_pool   = ibm_is_lb_pool.bootstrap_pool.id
}
resource "ibm_is_lb_pool" "bootstrap_pool" {
  #depends_on = [ibm_is_lb.load_balancer]

  name           = "bootstrap-node"
  lb             = var.load_balancer_id
  algorithm      = "round_robin"
  protocol       = "tcp"
  health_delay   = 5
  health_retries = 2
  health_timeout = 2
  health_type    = "tcp"
}
resource "ibm_is_lb_pool_member" "bootstrap" {
  depends_on = [ibm_is_lb_listener.bootstrap_listener]

  lb             = var.load_balancer_id
  pool           = ibm_is_lb_pool.bootstrap_pool.id
  port           = 22
  target_address = var.bootstrap_private_ip
}


