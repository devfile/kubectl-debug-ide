package pkg

import (
	"reflect"
	"testing"

	attributes2 "github.com/devfile/api/v2/pkg/attributes"
)

func Test_projectName(t *testing.T) {
	type args struct {
		remote string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "github https",
			args: args{remote: "https://github.com/l0rd/outyet"},
			want: "outyet",
		},
		{
			name: "github https with trailing slash",
			args: args{remote: "https://github.com/l0rd/outyet/"},
			want: "outyet",
		},
		{
			name: "github https with .git suffix",
			args: args{remote: "https://github.com/l0rd/outyet.git"},
			want: "outyet",
		},
		{
			name: "github ssh with .git suffix",
			args: args{remote: "git@github.com:l0rd/outyet.git"},
			want: "outyet",
		},
		{
			name: "github ssh with no .git suffix",
			args: args{remote: "git@github.com:l0rd/outyet"},
			want: "outyet",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := projectName(tt.args.remote)
			if (err != nil) != tt.wantErr {
				t.Errorf("projectName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("projectName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attributes(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "default attributes generation",
			want: []byte(defaultDevWorkspaceAttributes),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := attributes()
			var wantAttr attributes2.Attributes
			err := wantAttr.UnmarshalJSON(tt.want)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(got, wantAttr) {
				t.Errorf("attributes() = %v, want %s", got, wantAttr)
			}
		})
	}
}
