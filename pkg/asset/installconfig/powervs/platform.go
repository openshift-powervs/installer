package powervs

import (
	"fmt"
	"os"

	"github.com/openshift/installer/pkg/types/powervs"
)

// Platform collects powervs-specific configuration.
func Platform() (*powervs.Platform, error) {

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	var p powervs.Platform
	if osOverride := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); len(osOverride) != 0 {
		p.ClusterOSImage = osOverride
	}

	p.Region = ssn.Session.Region
	p.Zone = ssn.Session.Zone
	p.APIKey = ssn.Session.IAMToken
	p.UserID = ssn.Session.UserAccount

	return &p, nil
}

//ValidateAccountPermissions function validates account type and returns error
func ValidateAccountPermissions(client *Client) error {
	accType, err := client.GetAccountType()
	if err != nil {
		return err
	}

	if accType == "TRIAL" {
		return fmt.Errorf("Provided account is Trial account")
	}
	return nil
}
