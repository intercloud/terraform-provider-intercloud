package client

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

var (
	onceFamilies sync.Once
	families     = map[string]uuid.UUID{}
	errFamilies  error
)

type ResponseFamilies struct {
	Families []struct {
		ID        uuid.UUID `json:"id"`
		ShortName string    `json:"shortName"`
	} `json:"cloudDestinationFamilies"`
}

func (c *client) ReadDestinationFamilies() (output *ReadDestinationFamiliesOutput, err error) {

	onceFamilies.Do(func() {
		var resp *http.Response
		resp, errFamilies = c.apiClient.DoRequest(
			context.Background(),
			http.MethodGet,
			"/cloud-destination-families",
			nil,
			"",
			nil,
		)
		if errFamilies != nil {
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		var respFamilies *ResponseFamilies

		errFamilies = json.NewDecoder(resp.Body).Decode(&respFamilies)

		if errFamilies != nil {
			return
		}

		for idx := range respFamilies.Families {

			families[respFamilies.Families[idx].ShortName] = respFamilies.Families[idx].ID
		}

	})

	if errFamilies != nil {
		return nil, errFamilies
	}

	output = &ReadDestinationFamiliesOutput{
		Families: families,
	}

	return
}

func (c *client) GetFamilyFromID(id uuid.UUID) (name *string) {
	if len(families) == 0 {
		_, _ = c.ReadDestinationFamilies()
	}

	for idx := range families {
		if families[idx] == id {
			idx := idx
			name = &idx
			break
		}
	}
	return
}
