package powervs

// Metadata contains Power VS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	APIKey         string `json:"APIKey"`
	BaseDomain     string `json:"BaseDomain"`
	CISInstanceCRN string `json:"cisInstanceCRN"`
	Region         string `json:"region"`
	VPCRegion      string `json:"vpcRegion"`
	Zone           string `json:"zone"`
}
