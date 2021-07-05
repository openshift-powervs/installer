package powervs

import (
	"fmt"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/powervs"
)

// Platform collects powervs-specific configuration.
func Platform() (*powervs.Platform, error) {
	regions := knownRegions()
	
        // TODO(cklokman): This section came from aws and transforms the response from knownRegions
        //                 into long and short regions to prompt the user for region select this section
        //                 need need to be different based on powervs's implementation of knownRegions
        //

        longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for id, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
        
	var regionTransform survey.Transformer = func(ans interface{}) interface{} {
		switch v := ans.(type) {
		case core.OptionAnswer:
			return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
		case string:
			return strings.SplitN(v, " ", 2)[0]
		}
		return ""
	}

        // TODO(cklokman): We need to verify that us_south is the correct defaultRegion to use
        //
	defaultRegion := "us_south"
	if !IsKnownRegion(defaultRegion) {
		panic(fmt.Sprintf("installer bug: invalid default powervs region %q", defaultRegion))
	}

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	defaultRegionPointer := ssn.region
	if defaultRegionPointer != nil && *defaultRegionPointer != "" {
		if IsKnownRegion(*defaultRegionPointer) {
			defaultRegion = *defaultRegionPointer
		} else {
			logrus.Warnf("Unrecognized powervs region %q, defaulting to %s", *defaultRegionPointer, defaultRegion)
		}
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The powervs region to be used for installation.",
				Default: fmt.Sprintf("%s (%s)", defaultRegion, regions[defaultRegion]),
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
		return nil, err
	}

	zones := knownZones(region)
	defaultZone := zones[0]

	longZones := make([]string, 0, len(regions))
        shortZones := make([]string, 0, len(regions))
        for id, location := range regions {
                longZones = append(longZones, fmt.Sprintf("%s (%s)", id, location))
                shortZones = append(shortZones, id)
        }
	
        var zoneTransform survey.Transformer = func(ans interface{}) interface{} {
                switch v := ans.(type) {
                case core.OptionAnswer:
                        return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
                case string:
                        return strings.SplitN(v, " ", 2)[0]
                }
                return ""
        }
        
        var zone string
        err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Zone",
				Help:    "The powervs zone within the region to be used for installation.",
				Default: fmt.Sprintf("%s", defaultZone),
				Options: longZones,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := zoneTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(shortZones, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return errors.Errorf("invalid zone %q", choice)
				}
				return nil
			}),
			Transform: zoneTransform,
		},
	}, &zone)
	if err != nil {
		return nil, err
	}

	return &powervs.Platform{
		Region: region,
                Zone:   zone,
	}, nil
}
