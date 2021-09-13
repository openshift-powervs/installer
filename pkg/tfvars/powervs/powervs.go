// Package powervs contains Power Virtual Servers-specific Terraform-variable logic.
package powervs

import (
	"encoding/json"

	"github.com/openshift/cluster-api-provider-powervs/pkg/apis/powervsprovider/v1alpha1"
)

// powervsRegionToVPCRegion based on:
// https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server&locale=en#creating-service
// https://github.com/ocp-power-automation/ocp4-upi-powervs/blob/master/docs/var.tfvars-doc.md
var powervsRegionToIBMRegion = map[string]string{
	"dal":     "us-south",
	"us-east": "us-east",
	"sao":     "br-sao",
	"tor":     "ca-tor",
	"mon":     "ca-mon",
	"eu-de-1": "eu-de",
	"eu-de-2": "eu-de",
	"lon":     "eu-gb",
	"syd":     "au-syd",
	"tok":     "jp-tok",
	"osa":     "jp-osa",
}

type config struct {
	ServiceInstanceID    string `json:"powervs_cloud_instance_id"`
	APIKey               string `json:"powervs_api_key"`
	PowerVSRegion        string `json:"powervs_region"`
	PowerVSZone          string `json:"powervs_zone"`
	VPCRegion            string `json:"powervs_vpc_region"`
	PowerVSResourceGroup string `json:"powervs_resource_group"`
	CISInstanceCRN       string `json:"powervs_cis_crn"`
	SSHKey               string `json:"powervs_ssh_key"`
	ImageID              string `json:"powervs_image_name"`
	NetworkIDs           string `json:"powervs_network_name"`
	VPCID                string `json:"powervs_vpc_id"`
	VPCSubnetName        string `json:"powervs_vpc_subnet_name"`
	BootstrapMemory      string `json:"powervs_bootstrap_memory"`
	BootstrapProcessors  string `json:"powervs_bootstrap_processors"`
	MasterMemory         string `json:"powervs_master_memory"`
	MasterProcessors     string `json:"powervs_master_processors"`
	ProcType             string `json:"powervs_proc_type"`
	SysType              string `json:"powervs_sys_type"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterConfigs        []*v1alpha1.PowerVSMachineProviderConfig
	APIKey               string
	PowerVSZone          string
	PowerVSResourceGroup string
	CISInstanceCRN       string
	VPCID                string
	VPCSubnetName        string
}

// TFVars generates Power VS-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	//@TODO: Add resource group to platform
	cfg := &config{
		ServiceInstanceID:    masterConfig.ServiceInstanceID,
		APIKey:               sources.APIKey,
		PowerVSRegion:        masterConfig.Region,
		PowerVSZone:          sources.PowerVSZone,
		VPCRegion:            powervsRegionToIBMRegion[masterConfig.Region],
		PowerVSResourceGroup: sources.PowerVSResourceGroup,
		CISInstanceCRN:       sources.CISInstanceCRN,
		SSHKey:               *masterConfig.KeyPairName,
		ImageID:              masterConfig.ImageID,
		NetworkIDs:           masterConfig.NetworkIDs[0],
		VPCID:                sources.VPCID,
		VPCSubnetName:        sources.VPCSubnetName,
		BootstrapMemory:      masterConfig.Memory,
		BootstrapProcessors:  masterConfig.Processors,
		MasterMemory:         masterConfig.Memory,
		MasterProcessors:     masterConfig.Processors,
		ProcType:             masterConfig.ProcType,
		SysType:              masterConfig.SysType,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
