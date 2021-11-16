// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*
Datasource to get the list of dhcp servers in a power instance
*/

const PIDhcpServers = "servers"

func dataSourceIBMPIDhcps() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIDhcpServersRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PIDhcpServers: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PIDhcpId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the DHCP Server",
						},
						PIDhcpStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the DHCP Server",
						},
						PIDhcpNetwork: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DHCP Server private network",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIDhcpServersRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIDhcpClient(sess, cloudInstanceID)
	dhcpServers, err := client.GetAll(cloudInstanceID)
	if err != nil {
		log.Printf("[DEBUG] get all DHCP failed %v", err)
		return fmt.Errorf("failed to perform get all DHCP operation with error %v", err)
	}

	result := make([]map[string]interface{}, 0, len(dhcpServers))
	for _, dhcpServer := range dhcpServers {
		server := map[string]interface{}{
			PIDhcpId:     *dhcpServer.ID,
			PIDhcpStatus: *dhcpServer.Status,
		}

		dhcpNetwork := dhcpServer.Network
		if dhcpNetwork != nil {
			server[PIDhcpNetwork] = *dhcpNetwork.ID
		}
		result = append(result, server)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PIDhcpServers, result)

	return nil
}
