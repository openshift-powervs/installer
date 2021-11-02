package powervs

// Metadata contains Power VS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
<<<<<<< HEAD
	CISInstanceCRN string `json:"cisInstanceCRN"`
	Region         string `json:"region"`
	Zone           string `json:"zone"`
=======
	Region string `json:"region"`
	Zone   string `json:"zone"`
>>>>>>> ce5d7615b (Squashing Power VS IPI commits)
}
