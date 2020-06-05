package link

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ConnectorPair represents a pair of two connectors uuids
type ConnectorPair struct {
	First  uuid.UUID
	Second uuid.UUID
}

// Equal tests the equality of between two ConnectorPair
// The two pairs are equals if they contain the same elements no matter the order
func (cp *ConnectorPair) Equal(other *ConnectorPair) bool {
	return cp.toSet().Equal(other.toSet())
}

// toSet transforms the connector pair into a set of two elements
func (cp *ConnectorPair) toSet() *schema.Set {
	connectors := []string{cp.First.String(), cp.Second.String()}
	ret := make([]interface{}, 0, 2)
	for _, a := range connectors {
		ret = append(ret, a)
	}
	return schema.NewSet(schema.HashString, ret)
}
