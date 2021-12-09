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

	// Validate VolumeIDs
	for i, volumeID := range p.VolumeIDs {
		_, err := uuid.Parse(volumeID)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("volumeIDs").Index(i), volumeID, "Volume ID provided is not a UUID"))
		}
	}

	// Validate Memory
	if p.Memory != "" {
		memory, err := strconv.ParseInt(p.Memory, 10, 64)
		if err == nil {
			if memory < 2 || memory > 64 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.Memory, "Memory must be from 2 to 64 GB"))
			}
		} else {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.Memory, "Memory must be a valid integer"))
		}
	}

	// Validqte Processors
	processors, err := strconv.ParseFloat(p.Processors, 64)
	if p.Processors != "" {
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
	}
	// Validate ProcType
	if p.ProcType != "" {
		if p.ProcType != "shared" && p.ProcType != "dedicated" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("procType"), p.ProcType, "ProcType must be either 'shared' or 'dedicated'"))
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
