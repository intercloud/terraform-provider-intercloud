package connector

type CspFamily int

const (
	CspFamilyAws CspFamily = iota
	CspFamilyAzure
	CspFamilyGcp
)

var (
	sliceCspFamilies = []string{"aws", "azure", "gcp"}
)

func (f CspFamily) String() string {
	return sliceCspFamilies[f]
}
