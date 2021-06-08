################################################################
# Configure the IBM Cloud provider
################################################################
<<<<<<< HEAD
variable "powervs_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"
}

variable "powervs_vpc_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = "eu-gb"
  ##default     = "eu-gb"
}

variable "powervs_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = "lon"
=======
variable "ibmcloud_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"

  validation {
    condition     = var.ibmcloud_api_key != "" && lower(var.ibmcloud_api_key) != "<key>"
    error_message = "The ibmcloud_api_key is required and cannot be empty."
  }
}

variable "ibmcloud_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = ""

  validation {
    condition     = var.ibmcloud_region != "" && lower(var.ibmcloud_region) != "<region>"
    error_message = "The ibmcloud_region is required and cannot be empty."
  }
}

variable "ibmcloud_zone" {
  type        = string
  description = "The zone of an IBM Cloud region where you want to create Power System resources"
  default     = ""

  validation {
    condition     = var.ibmcloud_zone != "" && lower(var.ibmcloud_zone) != "<zone>"
    error_message = "The ibmcloud_zone is required and cannot be empty."
  }
>>>>>>> 22eef076f... Adding power iaas service (#5)
}

variable "powervs_resource_group" {
  type        = string
  description = "The cloud instance resource group"
  default     = ""
}

<<<<<<< HEAD
variable "powervs_cloud_instance_id" {
  type        = string
  description = "The cloud instance ID of your account"
  ## TODO: erase default and set via install-config
  default = "e449d86e-c3a0-4c07-959e-8557fdf55482"
}

################################################################
# Configure storage
################################################################
variable "powervs_cos_instance_location" {
  type        = string
  description = "The location of your COS instance"
  default     = "global"
}

variable "powervs_cos_bucket_location" {
  type        = string
  description = "The location to create your COS bucket"
  default     = "us-east"
}

variable "powervs_cos_storage_class" {
  type        = string
  description = "The plan used for your COS instance"
  default     = "smart"
=======
variable "cloud_instance_id" {
  type        = string
  description = "The cloud instance ID of your account"
  default     = ""
>>>>>>> 22eef076f... Adding power iaas service (#5)
}

################################################################
# Configure instances
################################################################
<<<<<<< HEAD
variable "powervs_image_name" {
=======
variable "image_name" {
>>>>>>> 22eef076f... Adding power iaas service (#5)
  type        = string
  description = "Name of the image used by all nodes in the cluster."
}

<<<<<<< HEAD
variable "powervs_network_name" {
  type        = string
  description = "Name of the network used by the all nodes in the cluster."
  default     = "pvs-ipi-net"
}

variable "powervs_bootstrap_memory" {
  type        = string
  description = "Amount of memory, in  GiB, used by the bootstrap node."
  default     = "32"
}

variable "powervs_bootstrap_processors" {
  type        = string
  description = "Number of processors used by the bootstrap node."
  default     = "0.5"
}

variable "powervs_master_memory" {
  type        = string
  description = "Amount of memory, in  GiB, used by each master node."
  default     = "32"
}

variable "powervs_master_processors" {
  type        = string
  description = "Number of processors used by each master node."
  default     = "0.5"
}

variable "powervs_proc_type" {
=======
variable "network_name" {
  type        = string
  description = "Name of the network used by the all nodes in the cluster."
}

variable "bootstrap" {
  type = object({ memory = string, processors = string })
  default = {
    memory     = "32"
    processors = "0.5"
  }
}

variable "proc_type" {
>>>>>>> 22eef076f... Adding power iaas service (#5)
  type        = string
  description = "The type of processor mode for all nodes (shared/dedicated)"
  default     = "shared"
}

<<<<<<< HEAD
variable "powervs_sys_type" {
=======
variable "sys_type" {
>>>>>>> 22eef076f... Adding power iaas service (#5)
  type        = string
  description = "The type of system (s922/e980)"
  default     = "s922"
}

<<<<<<< HEAD
variable "powervs_ssh_key" {
  type        = string
  description = "Public key for keypair used to access cluster. Required when creating 'ibm_pi_instance' resources."
}

## TODO: Set this in install-config instead
variable "powervs_vpc_name" {
  type        = string
  description = "Name of the IBM Cloud Virtual Private Cloud (VPC) to setup the load balancer."
  default     = "powervs-ipi"
=======
# Must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character
# Length cannot exceed 14 characters when combined with cluster_id_prefix
variable "cluster_id" {
  type    = string
  default = ""

  validation {
    condition     = can(regex("^$|^[a-z0-9]+[a-zA-Z0-9_\\-.]*[a-z0-9]+$", var.cluster_id))
    error_message = "The cluster_id value must be a lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
  }

  validation {
    condition     = length(var.cluster_id) <= 14
    error_message = "The cluster_id value shouldn't be greater than 14 characters."
  }
}

variable "powervs_vpc_name" {
  type        = string
  description = "Name of the IBM Cloud Virtual Private Cloud (VPC) to setup the load balancer."
  default     = ""
>>>>>>> 22eef076f... Adding power iaas service (#5)
}

variable "powervs_vpc_subnet_name" {
  type        = string
  description = "Name of the VPC subnet having DirectLink access to the PowerVS private network"
<<<<<<< HEAD
  default     = "subnet2"
}

## TODO: Pass the CIS CRN from the installer program, refer the IBM Cloud code to see the implementation.
variable "powervs_cis_crn" {
  type        = string
  description = "The CRN of CIS instance to use."
}
=======
  default     = ""
}

>>>>>>> 22eef076f... Adding power iaas service (#5)
