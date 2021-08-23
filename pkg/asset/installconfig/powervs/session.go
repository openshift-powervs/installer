package powervs

import (
	"fmt"
	"sort"
	"os"
	"time"
	"strings"

	"github.com/pkg/errors"
        "github.com/sirupsen/logrus"

	survey "github.com/AlecAivazis/survey/v2"
        "github.com/AlecAivazis/survey/v2/core"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	
)

var (
	//reqAuthEnvs = []string{"IBMID", "IBMID_PASSWORD"}
	//optAuthEnvs = []string{"IBMCLOUD_REGION", "IBMCLOUD_ZONE"}
	//debug = false
	defSessionTimeout time.Duration = 9000000000000000000.0
	defRegion                       = "us_south"
)

// Session is an object representing a session for the IBM Power VS API.
type Session struct {
	Session *ibmpisession.IBMPISession
	Creds   *UserCredentials
	OsOverride string
}

// UserCredentials is an object representing the credentials used for IBM Power VS during
// the creation of the install_config.yaml
type UserCredentials struct {
	APIKey string
	UserID string
}

// GetSession returns an IBM Cloud session by using credentials found in default locations in order:
// env IBMID & env IBMID_PASSWORD,
// ~/.bluemix/config.json ? (see TODO below)
// and, if no creds are found, asks for them
/* @TODO: if you do an `ibmcloud login` (or in my case ibmcloud login --sso), you get
//  a very nice creds file at ~/.bluemix/config.json, with an IAMToken. There's no username,
//  though (just the account's owner id, but that's not the same). It may be necessary
//  to use the IAMToken vs the password env var mentioned here:
//  https://github.com/IBM-Cloud/power-go-client#ibm-cloud-sdk-for-power-cloud
//  Yes, I think we'll need to use the IAMToken. There's a two-factor auth built into the ibmcloud login,
//  so the password alone isn't enough. The IAMToken is generated as a result. So either:
     1) require the user has done this already and pull from the file
     2) ask the user to paste in their IAMToken.
     3) let the password env var be the IAMToken? (Going with this atm since it's how I started)
     4) put it into Platform {userid: , iamtoken: , ...}
*/
func GetSession() (*Session, error) {
	region, zone, err := getRegionInfo()
	if err != nil {
                return nil, errors.Wrap(err, "failed to get region info")
        }
	
	uc, err := getUserCreds()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credentials")
	}

        s, err := getPISession( region, zone, uc )
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ")
	}
        
	oso := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE")

	return &Session{Session: s, Creds: uc, OsOverride: oso}, nil
}

func getUserCreds()  (*UserCredentials, error) {
	var apikey, id string
	var err error
	
	if id = os.Getenv("IBMID"); len(id) == 0 {
                err = survey.Ask([]*survey.Question{
                        {
                                Prompt: &survey.Input{
                                        Message: "IBM Cloud User ID",
                                        Help:    "The login for \nhttps://cloud.ibm.com/",
                                },
                        },
                }, &id)
                if err != nil {
                        return nil, errors.New("Error saving the IBMID variable")
                }
        }
	
        if apikey = os.Getenv("API_KEY"); len(apikey) == 0 {
                err = survey.Ask([]*survey.Question{
                        {
                                Prompt: &survey.Password{
                                        Message: "IBM Cloud API Key",
                                        Help:    "The api key installation.\nhttps://cloud.ibm.com/iam/apikeys",
                                },
                        },
                }, &apikey)
                if err != nil {
                        return nil, errors.New("Error saving the API_KEY variable")
                }
        }
	
	uc := &UserCredentials{UserID: id, APIKey: apikey}

	return uc, err	
}

/*
//  https://github.com/IBM-Cloud/power-go-client/blob/master/ibmpisession/ibmpowersession.go
*/
func getPISession( region string, zone string, uc *UserCredentials ) (*ibmpisession.IBMPISession, error) {

	// @TOOD: query if region is multi-zone? or just pass through err...
	// @TODO: pass through debug?
	s, err := ibmpisession.New(uc.APIKey, region, false, defSessionTimeout, uc.UserID, zone)
	return s, err
}

// getRegionInfo   returns information usually collected by platform.go but is required here
// 		   because ibmpisession requires this information as part of session.go
//		   this is found here because session.go is responsible for collecting
//		   os.Getenv variables, and therefore isnt found in region.go
func getRegionInfo() (string, string, error) {
	var region, sessionRegion, zone, sessionZone string
	var err error
	
	// -------------------------
	// Region
	// -------------------------

	// Retrieve sessionRegion from the current enviornment context
	sessionRegion = os.Getenv("IBMCLOUD_REGION")
	// this can also be pulled from  ~/bluemix/config.json
	if r2 := os.Getenv("IC_REGION"); len(r2) > 0 {
		if len(region) > 0 && region != r2 {
			return "", "", errors.New(fmt.Sprintf("conflicting values for IBM Cloud Region: IBMCLOUD_REGION: %s and IC_REGION: %s", region, r2))
		}
		if len(region) == 0 {
			sessionRegion = r2
		}
	}

	if len(sessionRegion) > 0 {
		if IsKnownRegion(sessionRegion) {
			region = sessionRegion
		} else {
			logrus.Warnf("Unrecognized Power VS region %s, ignoring IC_REGION and/or IBMCLOUD_REGION", sessionRegion)
		}
	}

	// Prompt the user if a region was not found in the current enviornment context
	if region == "" {
		regions := knownRegions()

		longRegions := make([]string, 0, len(regions))
		shortRegions := make([]string, 0, len(regions))
		for id, location := range regions {
			longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
			shortRegions = append(shortRegions, id)
		}

		sort.Strings(longRegions)
		sort.Strings(shortRegions)

		var regionTransform survey.Transformer = func(ans interface{}) interface{} {
			switch v := ans.(type) {
			case core.OptionAnswer:
				return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
			case string:
				return strings.SplitN(v, " ", 2)[0]
			}
			return ""
		}

		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Select{
					Message: "Region",
					Help:    "The Power VS region to be used for installation.",
					// Default: fmt.Sprintf("%s (%s)", defaultRegion, regions[defaultRegion]),
					Options: longRegions,
				},
				Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
					choice := regionTransform(ans).(core.OptionAnswer).Value
					i := sort.SearchStrings(shortRegions, choice)
					if i == len(shortRegions) || shortRegions[i] != choice {
						return errors.Errorf("invalid region %q", choice)
					}
					return nil
				}),
				Transform: regionTransform,
			},
		}, &region)
		if err != nil {
			return "", "", err
		}
	}

	// -------------------------
	// Zone
	// -------------------------
	// There is an effort to remove zone, I am uncertain how this can be done as it is required by
	// ibmpisession, if zone is not passed we may want to assign zone to region, however I am using
	// the code from platform.go to prompt the user for zone here instead

	sessionZone = os.Getenv("IBMCLOUD_ZONE")

	if len(sessionZone) > 0 {
                if IsKnownZone(region, sessionZone) {
                        zone = sessionZone
                } else {
                        logrus.Warnf("Unrecognized Power VS zone %s, ignoring IBMCLOUD_ZONE", sessionZone)
                }
        }

	// Prompt the user if zone was not found in the current enviornment context
	if zone == "" {
		zones := knownZones(region)
		defaultZone := zones[0]

		var zoneTransform survey.Transformer = func(ans interface{}) interface{} {
			switch v := ans.(type) {
			case core.OptionAnswer:
				return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
			case string:
				return strings.SplitN(v, " ", 2)[0]
			}
			return ""
		}

		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Select{
					Message: "Zone",
					Help:    "The powervs zone within the region to be used for installation.",
					Default: fmt.Sprintf("%s", defaultZone),
					Options: zones,
				},
				Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
					choice := zoneTransform(ans).(core.OptionAnswer).Value
					i := sort.SearchStrings(zones, choice)
					if i == len(zones) || zones[i] != choice {
						return errors.Errorf("invalid zone %q", choice)
					}
					return nil
				}),
				Transform: zoneTransform,
			},
		}, &zone)
		if err != nil {
			return "", "", err
		}
	}
	
	return region, zone, err
}
