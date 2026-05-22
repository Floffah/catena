package gitserver

import (
	"net/http"
	"testing"

	"github.com/zeebo/assert"
)

func TestParseRequestPath(t *testing.T) {
	tests := []struct {
		name    string
		rawPath string
		want    requestPath
		wantOK  bool
	}{
		{
			name:    "upload pack discovery",
			rawPath: "/floffah/catena/info/refs",
			want: requestPath{
				Owner:      "floffah",
				Repository: "catena",
				GitPath:    "/info/refs",
			},
			wantOK: true,
		},
		{
			name:    "repository dot git suffix is removed",
			rawPath: "/floffah/catena.git/git-upload-pack",
			want: requestPath{
				Owner:      "floffah",
				Repository: "catena",
				GitPath:    "/git-upload-pack",
			},
			wantOK: true,
		},
		{
			name:    "receive pack route",
			rawPath: "/floffah/catena/git-receive-pack",
			want: requestPath{
				Owner:      "floffah",
				Repository: "catena",
				GitPath:    "/git-receive-pack",
			},
			wantOK: true,
		},
		{
			name:    "non git route is ignored",
			rawPath: "/floffah/catena/issues",
		},
		{
			name:    "path without git command is ignored",
			rawPath: "/floffah/catena",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseRequestPath(tt.rawPath)

			assert.That(t, ok == tt.wantOK)
			if !tt.wantOK {
				return
			}

			assert.That(t, got == tt.want)
		})
	}
}

func TestIsValidGitMethod(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		gitPath string
		want    bool
	}{
		{
			name:    "info refs uses get",
			method:  http.MethodGet,
			gitPath: "/info/refs",
			want:    true,
		},
		{
			name:    "upload pack uses post",
			method:  http.MethodPost,
			gitPath: "/git-upload-pack",
			want:    true,
		},
		{
			name:    "receive pack uses post",
			method:  http.MethodPost,
			gitPath: "/git-receive-pack",
			want:    true,
		},
		{
			name:    "info refs rejects post",
			method:  http.MethodPost,
			gitPath: "/info/refs",
		},
		{
			name:    "upload pack rejects get",
			method:  http.MethodGet,
			gitPath: "/git-upload-pack",
		},
		{
			name:    "unknown path rejects method",
			method:  http.MethodGet,
			gitPath: "/HEAD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.That(t, isValidGitMethod(tt.method, tt.gitPath) == tt.want)
		})
	}
}

func TestGetOperation(t *testing.T) {
	tests := []struct {
		name      string
		gitPath   string
		service   string
		want      gitOperation
		wantFound bool
	}{
		{
			name:      "upload pack post route",
			gitPath:   "/git-upload-pack",
			want:      gitOperationUploadPack,
			wantFound: true,
		},
		{
			name:      "receive pack post route",
			gitPath:   "/git-receive-pack",
			want:      gitOperationReceivePack,
			wantFound: true,
		},
		{
			name:      "upload pack discovery",
			gitPath:   "/info/refs",
			service:   "git-upload-pack",
			want:      gitOperationUploadPack,
			wantFound: true,
		},
		{
			name:      "receive pack discovery",
			gitPath:   "/info/refs",
			service:   "git-receive-pack",
			want:      gitOperationReceivePack,
			wantFound: true,
		},
		{
			name:    "unknown discovery service",
			gitPath: "/info/refs",
			service: "git-archive",
		},
		{
			name:    "unknown git path",
			gitPath: "/HEAD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := getOperation(tt.gitPath, tt.service)

			assert.That(t, found == tt.wantFound)
			if !tt.wantFound {
				return
			}

			assert.That(t, got == tt.want)
		})
	}
}
