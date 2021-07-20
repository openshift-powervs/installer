// Package powervs contains Power Virtual Servers-specific Terraform-variable logic.
package powervs

import (
	"encoding/json"

	"github.com/openshift/cluster-api-provider-powervs/pkg/apis/powervsprovider/v1alpha1"
)

type IBMCloud struct {
	IBMCloudAPIKey      string `json:"ibmcloud_api_key"`
	IBMCloudRegion      string `json:"ibmcloud_region"`
	IBMCloudZone        string `json:"ibmcloud_zone"`
	IAASClassicUserName string `json:"iaas_classic_username"`
	IAASClassicAPIKey   string `json:"iaas_classic_api_key"`
}

type config struct {
	IBMCloud             `json:",inline"`
	PowerVSResourceGroup string `json:"powervs_resource_group"`
	ImageID              string `json:"powervs_image_id"`
	NetworkName          string `json:"powervs_network_ids"`
	BootstrapMemory      string `json:"powervs_bootstrap_memory"`
	BootstrapProcessors  string `json:"powervs_bootstrap_processors"`
	MasterMemory         string `json:"powervs_master_memory"`
	MasterProcessors     string `json:"powervs_master_processors"`
	ProcType             string `json:"powervs_proc_type"`
	SysType              string `json:"powervs_sys_type"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterConfigs    []*v1alpha1.PowerVSMachineProviderConfig
	IBMCloud         IBMCloud
}

// TFVars generates Power VS-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	//@TODO: Add resource group to platform
	//  -- change ImageID to ImageURL here?
	//  --
	cfg := &config{
		IBMCloud:             sources.IBMCloud,
		PowerVSResourceGroup: "powervs-ipi-resourcegroup",
		ImageID:              masterConfig.ImageID,
		BootstrapMemory:      masterConfig.Memory,
		BootstrapProcessors:  masterConfig.Processors,
		MasterMemory:         masterConfig.Memory,
		MasterProcessors:     masterConfig.Processors,
		ProcType:             masterConfig.ProcType,
		SysType:              masterConfig.SysType,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
