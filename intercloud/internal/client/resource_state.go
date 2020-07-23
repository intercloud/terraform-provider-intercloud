package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ResourceState int

const (
	ResourceStateRegistered ResourceState = iota
	ResourceStateInDeployment
	ResourceStateDeployed
	ResourceStateDelivered
	ResourceStateError
)

var (
	sliceResourceStates = [...]string{
		"registered",
		"in-deployment",
		"deployed",
		"delivered",
		"error",
	}
)

func GetResourceState(s string) (ResourceState, error) {
	for idx := range sliceResourceStates {
		if s == sliceResourceStates[idx] {
			cs := ResourceState(idx)
			return cs, nil
		}
	}
	return 0, fmt.Errorf("resource state is invalid (state = %q)", s)
}

func (s ResourceState) String() string {
	return sliceResourceStates[s]
}

func (u *ResourceState) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	*u, err = GetResourceState(s)
	return
}

// MarshalText implements TextMarshaler, invoked when encoding JSON.
func (u ResourceState) MarshalJSON() ([]byte, error) {
	if int(u) < len(sliceResourceStates) {
		return json.Marshal(u.String())
	}

	return nil, errors.New("resource state in invalid")
}
