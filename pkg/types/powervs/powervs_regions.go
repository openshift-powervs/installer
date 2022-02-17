package powervs

// Since there is no API to query these, we have to hard-code them here.

// Region describes resources associated with a region in Power VS.
// We're using a few items from the IBM Cloud VPC offering. The region names
// for VPC are different so another function of this is to correlate those.
type Region struct {
	Name        string
	Description string
	VPCRegion   string
	Zones       []string
}

// Regions holds the regions for IBM Power VS, and descriptions used during the survey
var Regions = map[string]Region{
	"dal": {
		Name:        "dal",
		Description: "Dallas, USA",
		VPCRegion:   "us-south",
		Zones:       []string{"dal12"},
	},
	"eu-de": {
		Name:        "eu-de",
		Description: "Frankfurt, Germany",
		VPCRegion:   "eu-de",
		Zones: []string{
			"eu-de-1",
			"eu-de-2",
		},
	},
	"lon": {
		Name:        "lon",
		Description: "London, UK.",
		VPCRegion:   "eu-gb",
		Zones: []string{
			"lon04",
			"lon06",
		},
	},
	"osa": {
		Name:        "osa",
		Description: "Osaka, Japan",
		VPCRegion:   "jp-osa",
		Zones:       []string{"osa21"},
	},
	"syd": {
		Name:        "syd",
		Description: "Sydney, Australia",
		VPCRegion:   "au-syd",
		Zones:       []string{"syd04"},
	},
	"sao": {
		Name:        "sao",
		Description: "São Paulo, Brazil",
		VPCRegion:   "br-sao",
		Zones:       []string{"sao01"},
	},
	"tor": {
		Name:        "tor",
		Description: "Toronto, Canada",
		VPCRegion:   "ca-tor",
		Zones:       []string{"tor01"},
	},
	"tok": {
		Name:        "tok",
		Description: "Tokyo, Japan",
		VPCRegion:   "jp-tok",
		Zones:       []string{"tok04"},
	},
	"us-east": {
		Name:        "us-east",
		Description: "Washington DC, USA",
		VPCRegion:   "us-east",
		Zones:       []string{"us-east"},
	},
}

// Zones retrieves a slice of all zones in Power VS
func Zones() []string {
	var zones []string
	for _, r := range Regions {
		zones = append(zones, r.Zones...)
	}
	return zones
}

// ZonesForRegion returns Zones for a given region
func ZonesForRegion(region string) []string {
	var zones []string
	for _, r := range Regions {
		if r.Name != region {
			continue
		}
		zones = append(zones, r.Zones...)
	}
	return zones
}

func VPCRegionForRegion(region string) string {
	var vpcRegion string
	for _, r := range Regions {
		if r.Name == region {
			vpcRegion = r.VPCRegion
		}
	}
	return vpcRegion
}
