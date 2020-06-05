package client

import (
	"encoding/json"
	"errors"
	"strings"
)

type PublicPrefixes []string

func (p PublicPrefixes) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}
	return json.Marshal(strings.Join(p, ","))
}

func (p *PublicPrefixes) UnmarshalJSON(data []byte) error {
	if p == nil {
		return errors.New("PublicPrefixes: UnmarshalJSON on nil pointer")
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if len(s) > 0 {
		*p = strings.Split(s, ",")
	} else {
		*p = []string{}
	}

	return nil
}
