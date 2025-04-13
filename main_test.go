package main

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    Version
		wantErr bool
	}{
		{
			name:    "正常なバージョン",
			tag:     "v1.2.3",
			want:    Version{Major: 1, Minor: 2, Patch: 3},
			wantErr: false,
		},
		{
			name:    "不正なバージョン形式",
			tag:     "v1.2",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "数字以外の文字を含む",
			tag:     "v1.a.3",
			want:    Version{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch) {
				t.Errorf("ParseVersion() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		bumpType string
		want     string
		wantErr  bool
	}{
		{
			name:     "patchバージョンを上げる",
			tag:      "v1.2.3",
			bumpType: "patch",
			want:     "v1.2.4",
			wantErr:  false,
		},
		{
			name:     "minorバージョンを上げる",
			tag:      "v1.2.3",
			bumpType: "minor",
			want:     "v1.3.0",
			wantErr:  false,
		},
		{
			name:     "majorバージョンを上げる",
			tag:      "v1.2.3",
			bumpType: "major",
			want:     "v2.0.0",
			wantErr:  false,
		},
		{
			name:     "不正なバンプタイプ",
			tag:      "v1.2.3",
			bumpType: "invalid",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "不正なバージョン形式",
			tag:      "v1.2",
			bumpType: "patch",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpVersion(tt.tag, tt.bumpType)
			if (err != nil) != tt.wantErr {
				t.Errorf("BumpVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BumpVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
