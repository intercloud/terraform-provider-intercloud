package tests

import (
	"testing"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector/csp/gcp"
)

func TestConvertIntMedToHuman(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		intMed   int
		expected string
	}{
		{
			// unknown med => default "low"
			intMed:   0,
			expected: "low",
		},
		{
			intMed:   gcp.GcpMed["low"],
			expected: "low",
		},
		{
			intMed:   gcp.GcpMed["high"],
			expected: "high",
		},
	}

	// check tests cases
	for _, tc := range cases {
		output := gcp.IntMedToHuman(tc.intMed)
		if output != tc.expected {
			t.Fatalf("The converted human med does not match (expected = %s, actual = %s)", tc.expected, output)
		}
	}
}

func TestConvertHumanMedToInt(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		humanMed string
		expected int
	}{
		{
			// case insensitive
			humanMed: "LoW",
			expected: gcp.GcpMed["low"],
		},
		{
			// case insensitive
			humanMed: "HiGH",
			expected: gcp.GcpMed["high"],
		},
		{
			// unknown med
			humanMed: "fake-med",
			expected: gcp.GcpMed["low"],
		},
		{
			// empty med
			humanMed: "",
			expected: gcp.GcpMed["low"],
		},
		{
			humanMed: "low",
			expected: gcp.GcpMed["low"],
		},
		{
			humanMed: "high",
			expected: gcp.GcpMed["high"],
		},
	}

	// check tests cases
	for _, tc := range cases {
		output := gcp.HumanMedToInt(tc.humanMed)
		if output != tc.expected {
			t.Fatalf("The converted int med does not match (expected = %d, actual = %d)", tc.expected, output)
		}
	}

}
