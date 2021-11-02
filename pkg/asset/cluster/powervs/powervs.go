// Package powervs extracts Power VS metadata from install configurations.
package powervs

import (
<<<<<<< HEAD
	"context"

	icpowervs "github.com/openshift/installer/pkg/asset/installconfig/powervs"
=======
>>>>>>> ce5d7615b (Squashing Power VS IPI commits)
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// Metadata converts an install configuration to PowerVS metadata.
<<<<<<< HEAD
func Metadata(config *types.InstallConfig, meta *icpowervs.Metadata) *powervs.Metadata {
	cisCRN, _ := meta.CISInstanceCRN(context.TODO())

	return &powervs.Metadata{
		CISInstanceCRN: cisCRN,
		Region:         config.Platform.PowerVS.Region,
		Zone:           config.Platform.PowerVS.Zone,
=======
func Metadata(config *types.InstallConfig) *powervs.Metadata {
	return &powervs.Metadata{
		Region: config.Platform.PowerVS.Region,
		Zone:   config.Platform.PowerVS.Zone,
>>>>>>> ce5d7615b (Squashing Power VS IPI commits)
	}
}
