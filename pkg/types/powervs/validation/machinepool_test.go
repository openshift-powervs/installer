package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var minimalPool = &powervs.MachinePool{
	ServiceInstance: "637efe83-0078-46c1-a15d-7d0d3facd651",
	Name:            "powervs-test",
	Memory:          "5",
	Processors:      "1.25",
}

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name     string
		pool     *powervs.MachinePool
		expected string
	}{
		{
			name: "minimal",
			pool: minimalPool,
		},
		{
			name: "invalid serviceinstance",
			pool: createTestPool(&powervs.MachinePool{
				ServiceInstance: "powervs-service-instance",
			}),
			expected: `^test-path\.serviceinstance: Invalid value: "powervs-service-instance": Service Instance provided is not a UUID$`,
		},
		{
			name: "invalid name",
			pool: createTestPool(&powervs.MachinePool{
				Name: "powervs-test!",
			}),
			expected: `^test-path\.name: Invalid value: "powervs-test!": Only letters \(no accents\), numbers, underscores and dashes are allowed$`,
		},
		{
			name: "valid volumeIDs",
			pool: createTestPool(&powervs.MachinePool{
				VolumeIDs: []string{"c8b709c4-93f1-499e-915e-0820bcc51406", "587c5788-107f-4351-aabc-1652c54c4491"},
			}),
		},
		{
			name: "invalid volumeIDs",
			pool: createTestPool(&powervs.MachinePool{
				VolumeIDs: []string{"c8b709c4-93f1-499e-915e-0820bcc51406", "abc123"},
			}),
			expected: `^test-path\.volumeIDs\[1]: Invalid value: "abc123": Volume ID provided is not a UUID$`,
		},
		{
			name: "invalid memory under",
			pool: createTestPool(&powervs.MachinePool{
				Memory: "1",
			}),
			expected: `^test-path\.memory: Invalid value: "1": Memory must be from 2 to 64 GB$`,
		},
		{
			name: "invalid memory over",
			pool: createTestPool(&powervs.MachinePool{
				Memory: "65",
			}),
			expected: `^test-path\.memory: Invalid value: "65": Memory must be from 2 to 64 GB$`,
		},
		{
			name: "invalid memory string",
			pool: createTestPool(&powervs.MachinePool{
				Memory: "all",
			}),
			expected: `^test-path\.memory: Invalid value: "all": Memory must be a valid integer$`,
		},
		{
			name: "invalid processors under",
			pool: createTestPool(&powervs.MachinePool{
				Processors: "0",
			}),
			expected: `^test-path\.processors: Invalid value: "0": Number of processors must be from \.25 to 32 cores$`,
		},
		{
			name: "invalid processors over",
			pool: createTestPool(&powervs.MachinePool{
				Processors: "33",
			}),
			expected: `^test-path\.processors: Invalid value: "33": Number of processors must be from \.25 to 32 cores$`,
		},
		{
			name: "invalid processors string",
			pool: createTestPool(&powervs.MachinePool{
				Processors: "all",
			}),
			expected: `^test-path\.processors: Invalid value: "all": Processors must be a valid floating point number$`,
		},
		{
			name: "invalid processors increment",
			pool: createTestPool(&powervs.MachinePool{
				Processors: "1.33",
			}),
			expected: `^test-path\.processors: Invalid value: "1\.33": Processors must be in increments of \.25$`,
		},
		{
			name: "valid procType",
			pool: createTestPool(&powervs.MachinePool{
				ProcType: "shared",
			}),
		},
		{
			name: "invalid procType",
			pool: createTestPool(&powervs.MachinePool{
				ProcType: "none",
			}),
			expected: `^test-path\.procType: Invalid value: "none": ProcType must be either 'shared' or 'dedicated'$`,
		},
		{
			name: "valid imageID",
			pool: createTestPool(&powervs.MachinePool{
				ImageID: "ce752e31-65b6-48cd-9685-14678688cb6e",
			}),
		},
		{
			name: "invalid imageID",
			pool: createTestPool(&powervs.MachinePool{
				ImageID: "rhel8",
			}),
			expected: `^test-path\.imageID: Invalid value: "rhel8": Image ID provided is not a UUID$`,
		},
		{
			name: "valid networkIDs",
			pool: createTestPool(&powervs.MachinePool{
				NetworkIDs: []string{"b0ef7de9-60a0-463a-8b8a-9407b5b9cc10", "9bb5f47a-e820-4fd5-8feb-07386274eb4d"},
			}),
		},
		{
			name: "invalid networkIDs",
			pool: createTestPool(&powervs.MachinePool{
				NetworkIDs: []string{"b0ef7de9-60a0-463a-8b8a-9407b5b9cc10", "test-net"},
			}),
			expected: `^test-path\.networkIDs\[1]: Invalid value: "test-net": Network ID provided is not a UUID$`,
		},
		{
			name: "valid sysType",
			pool: createTestPool(&powervs.MachinePool{
				SysType: "s922",
			}),
		},
		{
			name: "invalid sysType",
			pool: createTestPool(&powervs.MachinePool{
				SysType: "p922",
			}),
			expected: `^test-path\.sysType: Invalid value: "p922": System type not recognized$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

// Create a MachinePool with minimal defaults set and updated with any fields specified in passed MachinePool
func createTestPool(required *powervs.MachinePool) *powervs.MachinePool {
	testPool := &powervs.MachinePool{}
	testPool.Set(minimalPool)
	testPool.Set(required)
	return testPool
}
