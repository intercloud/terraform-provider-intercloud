package client

import (
	"reflect"
	"testing"
)

func TestResourceState_MarshalJSON(t *testing.T) {

	tests := []struct {
		name    string
		u       ResourceState
		want    []byte
		wantErr bool
	}{
		{
			name:    "delivered",
			u:       ResourceStateDelivered,
			want:    []byte("\"delivered\""),
			wantErr: false,
		},
		{
			name:    "deployed",
			u:       ResourceStateDeployed,
			want:    []byte("\"deployed\""),
			wantErr: false,
		},
		{
			name:    "in-deployment",
			u:       ResourceStateInDeployment,
			want:    []byte("\"in-deployment\""),
			wantErr: false,
		},
		{
			name:    "error",
			u:       ResourceStateError,
			want:    []byte("\"error\""),
			wantErr: false,
		},
		{
			name:    "registered",
			u:       ResourceStateRegistered,
			want:    []byte("\"registered\""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceState.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceState.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceState_UnmarshalJSON(t *testing.T) {

	var state ResourceState

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		u       *ResourceState
		state   ResourceState
		args    args
		wantErr bool
	}{
		{
			name:  "delivered",
			u:     &state,
			state: ResourceStateDelivered,
			args: args{
				b: []byte("\"delivered\""),
			},
			wantErr: false,
		},
		{
			name:  "deployed",
			u:     &state,
			state: ResourceStateDeployed,
			args: args{
				b: []byte("\"deployed\""),
			},
			wantErr: false,
		},
		{
			name:  "error",
			u:     &state,
			state: ResourceStateError,
			args: args{
				b: []byte("\"error\""),
			},
			wantErr: false,
		},
		{
			name:  "in-deployment",
			u:     &state,
			state: ResourceStateInDeployment,
			args: args{
				b: []byte("\"in-deployment\""),
			},
			wantErr: false,
		},
		{
			name:  "registered",
			u:     &state,
			state: ResourceStateRegistered,
			args: args{
				b: []byte("\"registered\""),
			},
			wantErr: false,
		},
		{
			name:  "unknown",
			u:     &state,
			state: ResourceStateRegistered,
			args: args{
				b: []byte("\"unknown\""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr || *tt.u != tt.state {
				t.Errorf("ResourceState.UnmarshalJSON() error = %v, wantErr %v, state %v != %v", err, tt.wantErr, tt.state, *tt.u)
			}
		})
	}
}
