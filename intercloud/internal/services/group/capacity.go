package group

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const gbpsUnits = "Gbps"
const mbpsUnits = "Mbps"

// Converts an float64 capacity (Kbps) to a string representing
// a human readeable capacity with units (Mbps, Gbps)
// The human representation will have a precision of 1
//
// The returned capacity will be in Mbps for capacity < 1Gbps
// The returned capacity will be in Gbps for capacity >= 1Gbps
//
// examples:
//	- 5 000		 => 5Mbps
//	- 35 000 000 => 3.5Gbps (precision of 1)
func ConvertCapacityWithUnits(capacityKbps int) string {
	if capacityKbps == 0 {
		return ""
	}
	// capacity < 1Gbps
	floatCapacityKbps := float64(capacityKbps)
	if floatCapacityKbps < math.Pow10(6) {
		return strings.Replace(fmt.Sprintf("%.1f%s", floatCapacityKbps/math.Pow10(3), mbpsUnits), ".0", "", 1)
	}
	// floatCapacity >= 1Gbps
	return strings.Replace(fmt.Sprintf("%.1f%s", floatCapacityKbps/math.Pow10(6), gbpsUnits), ".0", "", 1)

}

// RevertCapacityWithUnits converts a human readable representation
// of a capacity to an integer capacity (Kbps)
//
// examples:
//	- 5Mbps		=> 5 000
//	- 3.5Gbps	=> 35 000 000
func RevertCapacityWithUnits(capacityWithUnits string) int {
	capacityPattern, _ := regexp.Compile("([0-9.]+)([MG]bps)")
	// unknown format
	if !capacityPattern.Match([]byte(capacityWithUnits)) {
		return 0
	}
	// extract capacity and units
	match := capacityPattern.FindStringSubmatch(capacityWithUnits)
	capacity, _ := strconv.ParseFloat(match[1], 64)
	units := match[2]
	// Mbps => Kbps
	if units == mbpsUnits {
		return int(capacity * math.Pow10(3))
	}
	// Gpbs => Kbps
	return int(capacity * math.Pow10(6))
}
