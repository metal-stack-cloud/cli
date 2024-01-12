package helpers

import "testing"

func TestTrimProvider(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		want    string
	}{
		{
			name:    "no provider",
			subject: "eqx-mu4",
			want:    "eqx-mu4",
		},
		{
			name:    "with provider",
			subject: "eqx-mu4@x-cellent@github",
			want:    "eqx-mu4@x-cellent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimProvider(tt.subject); got != tt.want {
				t.Errorf("TrimProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
