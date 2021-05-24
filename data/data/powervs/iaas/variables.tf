variable "ibmcloud_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"

  validation{
    condition = var.ibmcloud_api_key != "" && lower(var.ibmcloud_api_key) != "<key>"
    error_message   = "The ibmcloud_api_key is required and cannot be empty."
  }
}

variable "cloud_resource_group" {
  type        = string
  description = "The cloud instance resource group"
  default     = ""
}

variable "ibmcloud_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = ""

  validation{
    condition       = var.ibmcloud_region != "" && lower(var.ibmcloud_region) != "<region>"
    error_message   = "The ibmcloud_region is required and cannot be empty."
  }
}

# TODO(cklokman): Ultimatly this should be named similarly to our virtual machines
variable "service_instance_name" {
  type        = string
  description = "A name for our resource instance (service)."
  default     = "ipi-powervs-iaas"
}

variable "service_tags" {
  type        = list(string)
  description = "A list of tags for our resource instance."
  default     = []
}
