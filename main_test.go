package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
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

func TestShowDiff(t *testing.T) {
	// テスト用の一時的なGitリポジトリを作成
	dir := t.TempDir()

	// Gitリポジトリの初期化
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Gitリポジトリの初期化に失敗しました: %v", err)
	}

	// テスト用のファイルを作成してコミット
	testFile := dir + "/test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git addに失敗しました: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("初期コミットに失敗しました: %v", err)
	}

	// タグの作成
	cmd = exec.Command("git", "tag", "v1.0.0")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("タグの作成に失敗しました: %v", err)
	}

	// 新しい変更を追加
	if err := os.WriteFile(testFile, []byte("test2"), 0644); err != nil {
		t.Fatalf("テストファイルの更新に失敗しました: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git addに失敗しました: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "second commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("2回目のコミットに失敗しました: %v", err)
	}

	// テストケース
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		{
			name:    "変更がある場合",
			tag:     "v1.0.0",
			wantErr: false,
		},
		{
			name:    "存在しないタグ",
			tag:     "v0.0.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準出力をキャプチャ
			var buf bytes.Buffer
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			_, err := ShowDiff(tt.tag, dir)

			w.Close()
			os.Stdout = old
			buf.ReadFrom(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("ShowDiff() error = %v, wantErr %v", err, tt.wantErr)
			}

			// エラーが期待されない場合、出力を検証
			if !tt.wantErr {
				output := buf.String()
				if !strings.Contains(output, tt.tag+"からの変更") {
					t.Errorf("ShowDiff() output = %v, want contains %v", output, tt.tag+"からの変更")
				}
			}
		})
	}
}
