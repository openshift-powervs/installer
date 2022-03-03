package validation

import (
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervs"
)

var (
	// Regions is a map of IBM Cloud regions where PowerVS service is supported.
	// The key of the map is the short name of the region. The value
	// of the map is the Description of the region.
	Regions = map[string]string{
		"dal":     "Dallas, USA",
		"eu-de":   "Frankfurt, Germany",
		"lon":     "London, UK.",
		"osa":     "Osaka, Japan",
		"syd":     "Sydney, Australia",
		"sao":     "Sao Paulo, Brazil",
		"tor":     "Toronto, Canada",
		"tok":     "Tokyo, Japan",
		"us-east": "Washington DC, USA",
	}
	regionShortNames = func() []string {
		keys := make([]string, len(Regions))
		i := 0
		for r := range Regions {
			keys[i] = r
			i++
		}
		return keys
	}()
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *powervs.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	//validate ServiceInstanceID
	if p.ServiceInstanceID != "" {
		_, err := uuid.Parse(p.ServiceInstanceID)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("SericeInstanceID"), p.ServiceInstanceID, "ServiceInstanceID provided is not a UUID"))
		}
	}

	//Validate Region
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	} else if _, ok := Regions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, regionShortNames))
	}

	allErrs = append(allErrs, validateVPCConfig(p, fldPath)...)

	//validate DefaultMachinePlatform
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}

func validateVPCConfig(p *powervs.Platform, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.VPC != "" || len(p.Subnets) > 0 {
		if p.VPC == "" {
			allErrs = append(allErrs, field.Required(path.Child("vpc"), "vpc is required when specifying subnets"))
		}
		if len(p.Subnets) == 0 {
			allErrs = append(allErrs, field.Required(path.Child("subnets"), "subnets is required when specifying vpc"))
		}
	}
	return allErrs
}

//var schemeRE = regexp.MustCompile("^([^:]+)://")
