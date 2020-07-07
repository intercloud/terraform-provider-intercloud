package connector

type Family int

const (
	FamilyAws Family = iota
	FamilyAzure
	FamilyGcp
)

var (
	sliceFamilies = []string{"aws", "azure", "gcp"}
)

func (f Family) String() string {
	return sliceFamilies[f]
}

func AllFamilies() []string {
	return sliceFamilies
}
