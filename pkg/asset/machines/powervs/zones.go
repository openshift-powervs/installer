package powervs

import (
	"fmt"
	"github.com/openshift/installer/pkg/types/powervs"
)

// AvailabilityZones returns a list of supported zones for the specified region.
func AvailabilityZones(region string) ([]string, error) {
	var zones []string

	zones = powervs.ZonesForRegion(region)
	if zones == nil {
		return zones, fmt.Errorf("region not found %s", region)
	}

	return zones, nil
}
