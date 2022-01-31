package powervs

import (
	"github.com/pkg/errors"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// PlatformStages are the stages to run to provision the infrastructure in PowerVS.
var PlatformStages = []terraform.Stage{
	stages.NewStage("powervs", "cluster"),
	stages.NewStage("powervs", "bootstrap", stages.WithNormalBootstrapDestroy()),
	stages.NewStage("powervs", "bootstrap-routing", stages.WithCustomBootstrapDestroy(removeFromLoadBalancers)),
}

func removeFromLoadBalancers(s stages.SplitStage, directory string, extraArgs []string) error {
	_, err := terraform.Apply(directory, powervstypes.Name, s, append(extraArgs, "-var=powervs_expose_bootstrap=false")...)
	return errors.Wrap(err, "failed disabling bootstrap load balancing")
}
