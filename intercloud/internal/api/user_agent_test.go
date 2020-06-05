package api

import (
	"testing"
)

func TestUserAgentProduct_String(t *testing.T) {
	type fields struct {
		Name    string
		Version string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Test string agent",
			fields: fields{Name: "Agent 1", Version: "1.0.1"},
			want:   "Agent 1/1.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserAgentProduct{
				Name:    tt.fields.Name,
				Version: tt.fields.Version,
			}
			if got := u.String(); got != tt.want {
				t.Errorf("UserAgentProduct.String() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}

func TestUserAgentProducts_String(t *testing.T) {
	tests := []struct {
		name    string
		u       UserAgentProducts
		wantRet string
	}{

		{
			name: "Test string agents",
			u: []UserAgentProduct{
				{Name: "Agent 1", Version: "1.0.0"},
				{Name: "Agent 2", Version: "2.0.0"},
			},
			wantRet: "Agent 1/1.0.0 Agent 2/2.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := tt.u.String(); gotRet != tt.wantRet {
				t.Errorf("UserAgentProducts.String() = \"%v\", want \"%v\"", gotRet, tt.wantRet)
			}
		})
	}
}
