package tests

import (
	"testing"

	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/services/group"
)

func TestConvertCapacityWithUnits(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		capacityKbps int
		expected     string
	}{
		{
			capacityKbps: 0,
			expected:     "",
		},
		{
			capacityKbps: 5000,
			expected:     "5Mbps",
		},
		{
			capacityKbps: 999000,
			expected:     "999Mbps",
		},
		{
			capacityKbps: 1000000,
			expected:     "1Gbps",
		},
		{
			capacityKbps: 3500000,
			expected:     "3.5Gbps",
		},
		{
			capacityKbps: 345000000,
			expected:     "345Gbps",
		},
		{
			capacityKbps: 1000000000000,
			expected:     "1000000Gbps",
		},
	}

	// check tests cases
	for _, tc := range cases {
		output := group.ConvertCapacityWithUnits(tc.capacityKbps)
		if output != tc.expected {
			t.Fatalf("The converted capacity does not match (expected = %s, actual = %s)", tc.expected, output)
		}
	}

}

func TestRevertCapacityWithUnitsMbps(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		humanReadableCapacity string
		expected              int
	}{
		{
			humanReadableCapacity: "",
			expected:              0,
		},
		{
			humanReadableCapacity: "5.5Mbps",
			expected:              5500,
		},
		{
			humanReadableCapacity: "66Mbps",
			expected:              66000,
		},
		{
			humanReadableCapacity: "999Mbps",
			expected:              999000,
		},
		{
			humanReadableCapacity: "1Gbps",
			expected:              1000000,
		},
		{
			humanReadableCapacity: "175.5Gbps",
			expected:              175500000,
		},
	}

	// check tests cases
	for _, tc := range cases {
		output := group.RevertCapacityWithUnits(tc.humanReadableCapacity)
		if output != tc.expected {
			t.Fatalf("The converted capacity does not match (expected = %d, actual = %d)", tc.expected, output)
		}
	}

}
