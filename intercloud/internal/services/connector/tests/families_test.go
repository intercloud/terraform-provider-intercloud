package tests

import (
	"testing"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector"
)

func TestFamilies_String(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		family   connector.CspFamily
		expected string
	}{
		{
			family:   connector.CspFamilyAws,
			expected: "aws",
		},
		{
			family:   connector.CspFamilyAzure,
			expected: "azure",
		},
		{
			family:   connector.CspFamilyGcp,
			expected: "gcp",
		},
	}

	// check tests cases
	for _, tc := range cases {
		output := tc.family.String()
		if output != tc.expected {
			t.Fatalf("The string method failed (expected = %q, actual = %q)", tc.expected, output)
		}
	}

}
