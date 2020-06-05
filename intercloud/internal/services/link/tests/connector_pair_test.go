package tests

import (
	"testing"

	"github.com/google/uuid"
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud/internal/services/link"
)

func TestEqualConnectorPairs(t *testing.T) {
	// create two distinct uuid
	firstUUID, _ := uuid.NewUUID()
	secondUUID, _ := uuid.NewUUID()
	thirdUUID, _ := uuid.NewUUID()

	// tests cases declaration
	cases := []struct {
		first            link.ConnectorPair
		second           link.ConnectorPair
		expectedEquality bool
	}{
		{
			first:            link.ConnectorPair{},
			second:           link.ConnectorPair{},
			expectedEquality: true,
		},
		{
			// same uuids, same order
			first:            link.ConnectorPair{First: firstUUID, Second: secondUUID},
			second:           link.ConnectorPair{First: firstUUID, Second: secondUUID},
			expectedEquality: true,
		},
		{
			// same uuids, different order
			first:            link.ConnectorPair{First: firstUUID, Second: secondUUID},
			second:           link.ConnectorPair{First: secondUUID, Second: firstUUID},
			expectedEquality: true,
		},
		{
			// only one uuid in second connector pair
			first:            link.ConnectorPair{First: firstUUID, Second: secondUUID},
			second:           link.ConnectorPair{First: firstUUID, Second: firstUUID},
			expectedEquality: false,
		},
		{
			// three different uuids present
			first:            link.ConnectorPair{First: firstUUID, Second: secondUUID},
			second:           link.ConnectorPair{First: firstUUID, Second: thirdUUID},
			expectedEquality: false,
		},
	}

	// check tests cases
	for _, tc := range cases {
		output := tc.first.Equal(&tc.second)
		if output != tc.expectedEquality {
			t.Fatalf("The connector pairs equality does not match (expected = %t, actual = %t, first = %+v, second = %+v)", tc.expectedEquality, output, tc.first, tc.second)
		}
	}

}
