package powervs

import (
	"github.com/openshift/installer/pkg/types/powervs"
)

// Platform returns a powervs.Platform object with session data
func Platform() (*powervs.Platform, error) {

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	var p powervs.Platform
	if len(ssn.OsOverride) != 0 {
		p.BootstrapOSImage = ssn.OsOverride
		p.ClusterOSImage = ssn.OsOverride
	}

	p.Region = ssn.Session.Region
	p.Zone = ssn.Session.Zone
	p.APIKey = ssn.Creds.APIKey
	p.UserID = ssn.Creds.UserID

	return &p, nil
}
