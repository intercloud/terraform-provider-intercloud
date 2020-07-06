package tests

import (
	"testing"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/services/connector"
)

func TestConnections_Content(t *testing.T) {
	// tests cases declaration
	cases := []struct {
		connection             connector.Connection
		expectedFamily         string
		expectedConnectionType string
	}{
		{
			connection:     connector.ConnectionAws,
			expectedFamily: "aws",
		},
		{
			connection:             connector.ConnectionAwsHosted,
			expectedFamily:         "aws",
			expectedConnectionType: "hosted_connection",
		},
		{
			connection:     connector.ConnectionAzure,
			expectedFamily: "azure",
		},
		{
			connection:     connector.ConnectionGcp,
			expectedFamily: "gcp",
		},
	}

	// check tests cases
	for _, tc := range cases {
		outputFamily := tc.connection.Family()
		outputConnectionType := tc.connection.ConnectionType()
		if outputFamily != tc.expectedFamily || outputConnectionType != tc.expectedConnectionType {
			t.Fatalf("The connection data are not equal (connection = %q, expected family = %q, actual family = %q, expected type = %q, actual type = %q)", tc.connection, tc.expectedFamily, outputFamily, tc.expectedConnectionType, outputConnectionType)
		}
	}

}
