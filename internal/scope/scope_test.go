package scope

import "testing"

func TestInScope(t *testing.T) {
	tests := []struct {
		name            string
		seedHost        string
		candidateURL    string
		allowSubdomains bool
		allowExternal   bool
		want            bool
		wantErr         bool
	}{
		{
			name:         "exact match",
			seedHost:     "example.com",
			candidateURL: "https://example.com/page",
			want:         true,
		},
		{
			name:            "subdomain allowed",
			seedHost:        "example.com",
			candidateURL:    "https://blog.example.com/page",
			allowSubdomains: true,
			want:            true,
		},
		{
			name:            "subdomain disallowed",
			seedHost:        "example.com",
			candidateURL:    "https://blog.example.com/page",
			allowSubdomains: false,
			want:            false,
		},
		{
			name:          "external allowed",
			seedHost:      "example.com",
			candidateURL:  "https://other.com/page",
			allowExternal: true,
			want:          true,
		},
		{
			name:          "external disallowed",
			seedHost:      "example.com",
			candidateURL:  "https://other.com/page",
			allowExternal: false,
			want:          false,
		},
		{
			name:         "malformed URL",
			seedHost:     "example.com",
			candidateURL: "://bad-url",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InScope(tt.seedHost, tt.candidateURL, tt.allowSubdomains, tt.allowExternal)

			if (err != nil) != tt.wantErr {
				t.Fatalf("InScope() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("InScope() = %v, want %v", got, tt.want)
			}
		})
	}
}
