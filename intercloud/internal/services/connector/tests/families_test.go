package tests

import (
	"testing"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector"
)

func TestFamilies_String(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		family   connector.Family
		expected string
	}{
		{
			family:   connector.FamilyAws,
			expected: "aws",
		},
		{
			family:   connector.FamilyAzure,
			expected: "azure",
		},
		{
			family:   connector.FamilyGcp,
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
