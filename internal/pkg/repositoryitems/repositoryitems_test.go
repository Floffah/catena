package repositoryitems

import (
	"testing"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/zeebo/assert"
)

func TestReferenceFromParts(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		number int64
		want   string
	}{
		{
			name:   "repository item reference",
			prefix: "I",
			number: 1,
			want:   "I-1",
		},
		{
			name:   "custom prefix is preserved",
			prefix: "CAT",
			number: 42,
			want:   "CAT-42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.That(t, ReferenceFromParts(tt.prefix, tt.number) == tt.want)
		})
	}
}

func TestReference(t *testing.T) {
	repository := db.Repository{ItemPrefix: "CAT"}
	item := db.RepositoryItem{Number: 12}

	assert.That(t, Reference(repository, item) == "CAT-12")
}
