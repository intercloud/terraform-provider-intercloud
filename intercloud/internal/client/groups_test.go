package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/api"
	"github.com/stretchr/testify/assert"
)

func Test_client_ReadGroup(t *testing.T) {

	inputID := uuid.MustParse("de6324b3-bd24-462e-bf82-c0586c8bbc92")
	inputOrgID := uuid.MustParse("e731c2f4-0ff2-42fd-b681-f02a2da09530")

	type args struct {
		in *ReadGroupInput
	}
	tests := []struct {
		name    string
		c       Client
		args    args
		wantOut *ReadGroupOutput
		wantErr bool
	}{
		{
			name:    "nil input",
			wantErr: true,
			c: makeTestClient(
				func(req *http.Request) *http.Response {
					assert.NotNil(t, req)
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
						Header:     make(http.Header),
					}
				},
			),
		},
		{
			name:    "ok",
			wantErr: false,
			args: args{
				in: &ReadGroupInput{
					ID:             inputID,
					OrganizationID: inputOrgID,
				},
			},
			c: makeTestClient(
				func(req *http.Request) *http.Response {
					// Test request parameters
					assert.NotNil(t, req)
					assert.Equal(t, req.URL.String(), "http://test.test/groups/de6324b3-bd24-462e-bf82-c0586c8bbc92")
					assert.Equal(t, req.Method, http.MethodGet)

					return &http.Response{
						StatusCode: 200,
						// Send response to be tested
						Body: ioutil.NopCloser(bytes.NewBufferString(
							makeGroupResponseBody(
								inputID.String(),
								"group_1",
								"desc",
								"IRN:group_1",
								inputOrgID.String(),
							),
						)),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}
				},
			),
			wantOut: &ReadGroupOutput{
				ID:             inputID,
				OrganizationID: inputOrgID,
				Description:    "desc",
				Irn:            "IRN:group_1",
				Name:           "group_1",
			},
		},
		{
			name:    "not found",
			wantErr: true,
			args: args{
				in: &ReadGroupInput{
					ID:             inputID,
					OrganizationID: inputOrgID,
				},
			},
			c: makeTestClient(
				func(req *http.Request) *http.Response {
					assert.NotNil(t, req)
					return &http.Response{
						StatusCode: 404,
						Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
						Header:     make(http.Header),
					}
				},
			),
		},
		{
			name:    "forbidden",
			wantErr: true,
			args: args{
				in: &ReadGroupInput{
					ID:             inputID,
					OrganizationID: inputOrgID,
				},
			},
			c: makeTestClient(
				func(req *http.Request) *http.Response {
					assert.NotNil(t, req)
					return &http.Response{
						StatusCode: 403,
						Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
						Header:     make(http.Header),
					}
				},
			),
		},
		{
			name:    "Unauthorized",
			wantErr: true,
			args: args{
				in: &ReadGroupInput{
					ID:             inputID,
					OrganizationID: inputOrgID,
				},
			},
			c: makeTestClient(
				func(req *http.Request) *http.Response {
					assert.NotNil(t, req)
					return &http.Response{
						StatusCode: 401,
						Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
						Header:     make(http.Header),
					}
				},
			),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			gotOut, err := tt.c.ReadGroup(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.ReadGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("client.ReadGroup() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func makeGroupResponseBody(
	id string,
	name string,
	desc string,
	irn string,
	organizationID string,
) string {
	return fmt.Sprintf(`
{
	"id": "%s",
	"name": "%s",
	"description": "%s",
	"irn": "%s",
	"organisationId": "%s"
}`,
		id,
		name,
		desc,
		irn,
		organizationID,
	)
}

func makeTestClient(
	roundTrip func(*http.Request) *http.Response,
) Client {

	return newTestClient(
		&Config{
			Endpoint:           "http://test.test",
			PrivateAccessToken: "access_token",
			UserAgentProducts:  make(api.UserAgentProducts, 0),
		},
		roundTrip,
	)
}
