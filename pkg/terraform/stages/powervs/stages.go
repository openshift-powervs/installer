package powervs

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in PowerVS.
var PlatformStages = []terraform.Stage{
	stages.NewStage("powervs", "cluster"),
	stages.NewStage("powervs", "post-install"),
}
