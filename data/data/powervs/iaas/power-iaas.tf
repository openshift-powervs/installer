provider "ibm" {
 ibmcloud_api_key = var.ibmcloud_api_key
}

data "ibm_resource_group" "group" { 
  name = var.cloud_resource_group
}

resource "ibm_resource_instance" "resource_instance" { 
  name     = "${var.cluster_id}-power-iaas" 
  service  = "power-iaas"
  plan     = "power-virtual-server-group" 
  location = var.ibmcloud_region
  tags     = concat( var.service_tags, [ "${var.cluster_id}-power-iaas" ] )
  resource_group_id = data.ibm_resource_group.group.id 
  
  timeouts { 
    create = "10m"
    update = "10m"
    delete = "10m"
  } 
} 
