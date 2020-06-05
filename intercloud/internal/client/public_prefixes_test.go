package client

import (
	"reflect"
	"testing"
)

func TestPublicPrefixes_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		p       PublicPrefixes
		want    []byte
		wantErr bool
	}{
		{
			name:    "1 entry",
			p:       PublicPrefixes{"192.168.142.2/30"},
			want:    []byte("\"192.168.142.2/30\""),
			wantErr: false,
		},
		{
			name:    "2 entry",
			p:       PublicPrefixes{"192.168.142.2/30", "192.168.142.2/16"},
			want:    []byte("\"192.168.142.2/30,192.168.142.2/16\""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicPrefixes.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicPrefixes.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicPrefixes_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		p       *PublicPrefixes
		len     int
		args    args
		wantErr bool
	}{
		{
			name: "2 entries",
			p:    &PublicPrefixes{},
			len:  2,
			args: args{
				data: []byte("\"192.168.142.2/30,192.168.142.2/16\""),
			},
			wantErr: false,
		},
		{
			name: "1 entry",
			p:    &PublicPrefixes{},
			len:  1,
			args: args{
				data: []byte("\"192.168.142.2/30\""),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr || len(*tt.p) != tt.len {
				t.Errorf("PublicPrefixes.UnmarshalJSON() error = %v, wantErr %v len %v != %v, ", err, tt.wantErr, len(*tt.p), tt.len)
			}
		})
	}
}
