package auth

import "testing"

func TestNewSignature(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    Signature
		wantErr bool
	}{
		{
			name: "署名が作成できる",
			args: args{
				value: "signature",
			},
			want:    Signature("signature"),
			wantErr: false,
		},
		{
			name: "署名が作成できない",
			args: args{
				value: "",
			},
			want:    Signature(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewSignature(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
