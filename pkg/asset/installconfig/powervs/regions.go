package powervs 

import (
	"github.com/powervs/powervs-sdk-go/powervs/endpoints"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/rhcos"
)

// TODO(cklokman): aws uses another repo to provide endpoints, which seems to be used to query
//		   region information.  For my testing I implemented a very simple version of this
//		   that returns hard coded data to satisfy this, and return similar data for powervs
//		   I am uncertain the direction to take here, but don't believe hardcoding that data
//		   in this file is the correct solution.  It may be benificial to hold another repo
//		   in the same way as aws, or make api calls here to grab the required information.
//

// TODO(cklokman): This section came from aws and returns a map of regions, which will need to be
//                 a subset of powervs regions and the regions where RHEL CoreOS images are available.
//                 This process will be different from AWS, because of the method IBM is using for the
//                 CoreOS images.  It may be the case that all regions will have CoreOS images present
//                 in that case, this function becomes quite straight forward.  For the moment I am
//                 leaving this logic here, because this logic may be nessesary, and seems like it could
//                 be benificial to filter regions by those that are known to have CoreOS images in the
//                 case that not every region has a CoreOS, or the correct CoreOS images.
//
//                 I believe that the CoreOS filtering can happen on regions, and do not believe
//                 that CoreOS images would vary between zones within a region, but this should
//                 be verified, because it would change the logic here and in knownZones

// knownRegions is a list of AWS regions that the installer recognizes.
// This is subset of AWS regions and the regions where RHEL CoreOS images are published.
// The result is a map of region identifier to region description
func knownRegions() map[string]string {
	required := sets.NewString(rhcos.AMIRegions...)

	regions := make(map[string]string)

	// Partitions is probably incorrect termonology here, but I am mirroring endpoint funcitonality
	// until a more permanent solution to grabbing this information is nailed down.
	for _, partition := range endpoints.DefaultPartitions() {
		for _, partitionRegion := range partition.Regions() {
			partitionRegion := partitionRegion
			if required.Has(partitionRegion.ID()) {
				regions[partitionRegion.ID()] = partitionRegion.Description()
			}
		}
	}
	return regions
}

// IsKnownRegion return true is a specified region is Known to the installer.
// A known region is subset of AWS regions and the regions where RHEL CoreOS images are published.
func IsKnownRegion(region string) bool {
        if _, ok := knownRegions()[region]; ok {
                return true
        }
        return false
}

// TODO(cklokman): We may want to drop zone descriptions here, and I am just mirroring region
//		   functionality and using it for the survey prompts.
//

// knownZones is a list of powervs zones that coorespond to the region provided.
func knownZones(region string) map[string]string {
	zones := make(map[string]string)
	partitions := endpoints.DefaultPartitions()
	regions := partitions.Regions[region]
	for _, zone := range regions.Zones() {
		zone := zone
		zones[zone.ID()] = zone.Description()
	}
	return zones
}

// IsKnownZone return true is a specified zone is Known to the installer.
func IsKnownRegion(region string, zone string) bool {
	if _, ok := knownRegions()[region]; ok {
		if _, ok := knownZones(region)[zone]; ok {
			return true
		}
		return false
	}
	return false
}
