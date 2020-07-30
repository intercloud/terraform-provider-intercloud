package tests

import (
	"testing"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector"
)

func TestConnections_Content(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		connection             connector.ConnectionFamily
		expectedCspFamily      string
		expectedConnectionType string
		expectedName           string
	}{
		{
			connection:        connector.ConnectionFamilyAws,
			expectedCspFamily: "aws",
			expectedName:      "aws",
		},
		{
			connection:             connector.ConnectionFamilyAwsHosted,
			expectedCspFamily:      "aws",
			expectedConnectionType: "hosted_connection",
			expectedName:           "awshostedconnection",
		},
		{
			connection:        connector.ConnectionFamilyAzure,
			expectedCspFamily: "azure",
			expectedName:      "azure",
		},
		{
			connection:        connector.ConnectionFamilyGcp,
			expectedCspFamily: "gcp",
			expectedName:      "gcp",
		},
	}

	// check tests cases
	for _, tc := range cases {
		outputFamily := tc.connection.CspFamily()
		outputConnectionType := tc.connection.ConnectionType()
		outputName := tc.connection.String()
		if outputFamily != tc.expectedCspFamily || outputConnectionType != tc.expectedConnectionType || outputName != tc.expectedName {
			t.Fatalf("The connection data are not equal (connection = %q, expected family = %q, current family = %q, expected type = %q, current type = %q, expected name = %q, current name = %q)", tc.connection, tc.expectedCspFamily, outputFamily, tc.expectedConnectionType, outputConnectionType, outputName, tc.expectedName)
		}
	}

}
