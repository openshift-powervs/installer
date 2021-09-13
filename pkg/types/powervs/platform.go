package powervs

// Platform stores all the global configuration that all machinesets
// use.
/// used by the installconfig, and filled in by the installconfig/platform/powervs::Platform() func
type Platform struct {

	// ServiceInstanceID is the ID of the Power IAAS instance created from the IBM Cloud Catalog
	ServiceInstanceID string `json:"serviceInstance"`

	// PowerVSResourceGroup is the resource group for creating Power VS resources.
	PowerVSResourceGroup string `json:"powervsResourceGroup"`

	// Region specifies the IBM Cloud region where the cluster will be created.
	Region string `json:"region"`

	// Zone specifies the IBM Cloud colo region where the cluster will be created.
	// Required for multi-zone regions.
	Zone string `json:"zone"`

	// UserID is the login for the user's IBM Cloud account.
	UserID string `json:"userID"`

	// APIKey is the API key for the user's IBM Cloud account.
	//
	// +optional
	APIKey string `json:"APIKey,omitempty"`

	// VPC is a VPC inside IBM Cloud. Needed in order to create VPC Load Balancers.
	//
	// +optional
	VPC string `json:"vpc,omitempty"`

	// Subnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	// @TODO: Rename VPCSubnetID & make into string
	//
	// +optional
	Subnets []string `json:"subnets,omitempty"`

	// PVSNetworkID specifies an existing network withing the Power VS Service Instance.
	// Leave unset to have the installer create the network.
	//
	// +optional
	PVSNetworkID string `json:"pvsNetworkID,omitempty"`

	// UserTags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	// BootstrapOSImage is a URL to override the default OS image
	// for the bootstrap node. The URL must contain a sha256 hash of the image
	// e.g https://mirror.example.com/images/image.ova.gz?sha256=a07bd...
	//
	// +optional
	BootstrapOSImage string `json:"bootstrapOSImage,omitempty"`

	// ClusterOSImage is a URL to override the default OS image
	// for cluster nodes. The URL must contain a sha256 hash of the image
	// e.g https://mirror.example.com/images/powervs.ova.gz?sha256=3b5a8...
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Power VS for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
