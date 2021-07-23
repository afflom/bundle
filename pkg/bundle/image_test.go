package bundle

import (
	"testing"
)

func Test_ImageBlocking(t *testing.T) {
	type fields struct {
		blockedImages []string
	}
	tests := []struct {
		name   string
		fields fields
		ref    string
		want   bool
	}{
		{
			name: "testing want to block",
			fields: fields{
				blockedImages: []string{"alpine"},
			},
			ref:  "alpine",
			want: true,
		},
		{
			name: "testing do not want to block",
			fields: fields{
				blockedImages: []string{"alpine"},
			},
			ref:  "not",
			want: false,
		},
	}
	for _, tt := range tests {
		m := &BundleSpec{
			Mirror: Mirror{
				BlockedImages: tt.fields.blockedImages,
			},
		}

		actual := IsBlocked(m, tt.ref)

		if actual != tt.want {
			t.Errorf("Test %s: Expected '%v', got '%v'", tt.name, tt.want, actual)
		}

	}
}
