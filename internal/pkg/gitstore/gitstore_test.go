package gitstore

import "testing"

func TestNormalizeGitDirectory(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		want      string
		wantErr   bool
	}{
		{
			name:      "empty path is root",
			directory: "",
			want:      "",
		},
		{
			name:      "slash path is root",
			directory: "/",
			want:      "",
		},
		{
			name:      "leading slash is repository relative",
			directory: "/docs",
			want:      "docs",
		},
		{
			name:      "nested path is preserved",
			directory: "docs/guides",
			want:      "docs/guides",
		},
		{
			name:      "dot segments are ignored",
			directory: "./docs/./guides",
			want:      "docs/guides",
		},
		{
			name:      "parent segment is rejected",
			directory: "../docs",
			wantErr:   true,
		},
		{
			name:      "nested parent segment is rejected",
			directory: "docs/../guides",
			wantErr:   true,
		},
		{
			name:      "nul byte is rejected",
			directory: "docs\x00guides",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeGitDirectory(tt.directory)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
