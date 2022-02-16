data "ibm_is_vpc" "vpc" {
  name = var.vpc_name
}

resource "ibm_is_security_group" "ocp_security_group" {
  name           = "${var.cluster_id}-ocp-sec-group"
  resource_group = data.ibm_resource_group.resource_group.id
  vpc            = data.ibm_is_vpc.vpc.id
  tags           = [var.cluster_id]
}

resource "ibm_is_security_group_rule" "outbound_any" {
  group     = ibm_is_security_group.ocp_security_group.id
  direction = "outbound"
}
