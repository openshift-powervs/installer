// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

func dataSourceIBMPICloudConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMPICloudConnectionRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PICloudConnectionName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Cloud Connection Name to be used",
				ValidateFunc: validation.NoZeroValues,
			},

			// Start of Computed Attributes
			helpers.PICloudConnectionSpeed: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			helpers.PICloudConnectionGlobalRouting: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			helpers.PICloudConnectionMetered: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			helpers.PICloudConnectionStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			helpers.PICloudConnectionIBMIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			helpers.PICloudConnectionUserIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			helpers.PICloudConnectionPort: {
				Type:     schema.TypeString,
				Computed: true,
			},
			helpers.PICloudConnectionNetworks: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of Networks attached to this cloud connection",
			},
			helpers.PICloudConnectionClassicEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable classic endpoint destination",
			},
			// helpers.PICloudConnectionClassicGreCidr: {
			// 	Type:         schema.TypeString,
			// 	Computed:     true,
			// 	Description:  "GRE network in CIDR notation",
			// },
			helpers.PICloudConnectionClassicGreDest: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GRE destination IP address",
			},
			helpers.PICloudConnectionClassicGreSource: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GRE auto-assigned source IP address",
			},
			helpers.PICloudConnectionVPCEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable VPC for this cloud connection",
			},
			helpers.PICloudConnectionVPCCRNs: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of VPCs attached to this cloud connection",
			},
		},
	}
}

func dataSourceIBMPICloudConnectionRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	cloudConnectionName := d.Get(helpers.PICloudConnectionName).(string)
	client := instance.NewIBMPICloudConnectionClient(sess, cloudInstanceID)

	// cloudConnectionD, err := client.Get(cloudConnectionName, cloudInstanceID)
	// if err != nil {
	// 	log.Printf("[DEBUG] get cloud connection failed %v", err)
	// 	return fmt.Errorf(errors.GetCloudConnectionOperationFailed, cloudConnectionName, err)
	// }

	// Work around for GET by Name not working
	cloudConnections, err := client.GetAll(cloudInstanceID)
	if err != nil {
		log.Printf("[DEBUG] get cloud connections failed %v", err)
		return fmt.Errorf("failed to perform Get Cloud Connections Operation with error %v", err)
	}
	var cloudConnection *models.CloudConnection
	for _, cc := range cloudConnections.CloudConnections {
		if cloudConnectionName == *cc.Name {
			cloudConnection = cc
			break
		}
	}
	if cloudConnection == nil {
		log.Printf("[DEBUG] get cloud connections failed %v", err)
		return fmt.Errorf("failed to perform Get Cloud Connections Operation for name %s with error %v", cloudConnectionName, err)
	}

	d.SetId(*cloudConnection.CloudConnectionID)
	d.Set(helpers.PICloudConnectionName, cloudConnection.Name)
	d.Set(helpers.PICloudConnectionGlobalRouting, cloudConnection.GlobalRouting)
	d.Set(helpers.PICloudConnectionMetered, cloudConnection.Metered)
	d.Set(helpers.PICloudConnectionIBMIPAddress, cloudConnection.IbmIPAddress)
	d.Set(helpers.PICloudConnectionUserIPAddress, cloudConnection.UserIPAddress)
	d.Set(helpers.PICloudConnectionStatus, cloudConnection.LinkStatus)
	d.Set(helpers.PICloudConnectionPort, cloudConnection.Port)
	d.Set(helpers.PICloudConnectionSpeed, cloudConnection.Speed)
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)
	if cloudConnection.Networks != nil {
		networks := make([]string, len(cloudConnection.Networks))
		for i, ccNetwork := range cloudConnection.Networks {
			if ccNetwork != nil {
				networks[i] = *ccNetwork.NetworkID
			}
		}
		d.Set(helpers.PICloudConnectionNetworks, networks)
	}
	if cloudConnection.Classic != nil {
		d.Set(helpers.PICloudConnectionClassicEnabled, cloudConnection.Classic.Enabled)
		if cloudConnection.Classic.Gre != nil {
			d.Set(helpers.PICloudConnectionClassicGreDest, cloudConnection.Classic.Gre.DestIPAddress)
			d.Set(helpers.PICloudConnectionClassicGreSource, cloudConnection.Classic.Gre.SourceIPAddress)
		}
	}
	if cloudConnection.Vpc != nil {
		d.Set(helpers.PICloudConnectionVPCEnabled, cloudConnection.Vpc.Enabled)
		if cloudConnection.Vpc.Vpcs != nil && len(cloudConnection.Vpc.Vpcs) > 0 {
			vpcCRNs := make([]string, len(cloudConnection.Vpc.Vpcs))
			for i, vpc := range cloudConnection.Vpc.Vpcs {
				vpcCRNs[i] = *vpc.VpcID
			}
			d.Set(helpers.PICloudConnectionVPCCRNs, vpcCRNs)
		}
	}

	return nil
}
