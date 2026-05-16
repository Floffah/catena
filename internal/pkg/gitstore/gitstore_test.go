package gitstore

import (
	"testing"

	"github.com/zeebo/assert"
)

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
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.That(t, got == tt.want)
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
	assert.That(t, len(got) == len(want))

	for i := range want {
		assert.That(t, got[i] == want[i])
	}
}
