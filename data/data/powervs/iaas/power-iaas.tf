provider "ibm" {
 ibmcloud_api_key = var.ibmcloud_api_key
}

data "ibm_resource_group" "group" { 
  name = var.cloud_resource_group
}

resource "ibm_resource_instance" "resource_instance" { 


  name     = var.service_instance_name
  service  = "power-iaas"
  plan     = "power-virtual-server-group" 
  location = var.ibmcloud_region
  tags     = var.service_tags 
  tags = merge(
    {
      "Name" = "${var.cluster_id}-power-iaas"
    },
    var.tags,
  )
  resource_group_id = data.ibm_resource_group.group.id 
  
  timeouts { 
    create = "5m"
    update = "5m"
    delete = "5m"
  } 
} 
