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

func TestSortTreeEntries(t *testing.T) {
	size := int64(10)
	entries := []TreeEntry{
		{Name: "zeta.go", Path: "zeta.go", Type: "blob", Size: &size},
		{Name: "alpha.go", Path: "alpha.go", Type: "blob", Size: &size},
		{Name: "vendor", Path: "vendor", Type: "commit"},
		{Name: "Docs", Path: "Docs", Type: "tree"},
		{Name: "app", Path: "app", Type: "tree"},
		{Name: "README.md", Path: "README.md", Type: "blob", Size: &size},
	}

	sortTreeEntries(entries)

	got := make([]string, 0, len(entries))
	for _, entry := range entries {
		got = append(got, entry.Name)
	}

	want := []string{"app", "Docs", "vendor", "alpha.go", "README.md", "zeta.go"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
