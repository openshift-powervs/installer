package powervs

// Metadata contains Power VS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	CISInstanceCRN string `json:"cisInstanceCRN"`
	Region         string `json:"region"`
	Zone           string `json:"zone"`
	// ServiceEndpoints list contains custom endpoints which will override default
	// service endpoint of Power VS/IBM Services.
	// There must be only one ServiceEndpoint for a service.
	ServiceEndpoints []ServiceEndpoint `json:"serviceEndpoints,omitempty"`
}
