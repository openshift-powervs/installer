package validation

import (
	"github.com/google/uuid"
	"github.com/openshift/installer/pkg/types/powervs"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"math"
	"regexp"
	"strconv"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *powervs.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate ServiceInstance
	if !isUUID(p.ServiceInstance) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceinstance"), p.ServiceInstance, "Service Instance provided is not a UUID"))
	}

	// Validate Name
	// Check restrictions on character set and length imposed by Power VS
	if !regexp.MustCompile(`^[a-zA-Z0-9-_]{1,}$`).MatchString(p.Name) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), p.Name, "Only letters (no accents), numbers, underscores and dashes are allowed"))
	}
	if len(p.Name) > 47 {
		allErrs = append(allErrs, field.TooLong(fldPath, p.Name, 47))
	}

	// Validate VolumeIDs
	for i, volumeID := range p.VolumeIDs {
		if !isUUID(volumeID) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("volumeIDs").Index(i), volumeID, "Volume ID provided is not a UUID"))
		}
	}

	// Validate Memory
	memory, err := strconv.ParseInt(p.Memory, 10, 64)
	if err == nil {
		if memory < 2 || memory > 64 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.Memory, "Memory must be from 2 to 64 GB"))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.Memory, "Memory must be a valid integer"))
	}

	// Validqte Processors
	processors, err := strconv.ParseFloat(p.Processors, 64)
	if err == nil {
		if processors < 0.25 || processors > 32 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), p.Processors, "Number of processors must be from .25 to 32 cores"))
		}
		if math.Mod(processors, 0.25) != 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), p.Processors, "Processors must be in increments of .25"))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), p.Processors, "Processors must be a valid floating point number"))
	}

	// Validate optional fields
	// Validate ProcType
	if p.ProcType != "" {
		if p.ProcType != "shared" && p.ProcType != "dedicated" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("procType"), p.ProcType, "ProcType must be either 'shared' or 'dedicated'"))
		}
	}

	// Validate ImageID
	if p.ImageID != "" {
		if !isUUID(p.ImageID) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("imageID"), p.ImageID, "Image ID provided is not a UUID"))
		}
	}

	// Validate NetworkIDs
	for i, networkID := range p.NetworkIDs {
		if !isUUID(networkID) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("networkIDs").Index(i), networkID, "Network ID provided is not a UUID"))
		}
	}

	// Validate SysType
	if p.SysType != "" {
		sysTypeRegex := `^(?:e980|s922(-.*|))$`
		if !regexp.MustCompile(sysTypeRegex).MatchString(p.SysType) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("sysType"), p.SysType, "System type not recognized"))
		}
	}
	return allErrs
}

func isUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
