package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func main() {
	// 引数の取得（デフォルトは "patch"）
	bumpType := "patch"
	if len(os.Args) > 1 {
		bumpType = strings.ToLower(os.Args[1])
	}

	// 最新タグの取得
	latestTag, err := GetLatestTag()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// バージョンの更新
	newTag, err := BumpVersion(latestTag, bumpType)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// 差分の表示
	if err := ShowDiff(latestTag, "."); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// ユーザーに確認
	fmt.Printf("\n新しいタグ %s を作成しますか？ [y/N]: ", newTag)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input != "y" && input != "yes" {
		fmt.Println("タグの作成をキャンセルしました")
		os.Exit(0)
	}

	// タグの作成
	if err := CreateTag(newTag); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ タグ %s を作成しました。\n", newTag)
}

func GetLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		// タグが存在しない場合は v0.0.0 を返す
		return "v0.0.0", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func ShowDiff(tag string, workDir string) error {
	cmd := exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..HEAD", tag))
	cmd.Dir = workDir
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("差分の取得に失敗しました: %v", err)
	}

	if len(output) == 0 {
		fmt.Printf("%sからの変更はありません\n", tag)
		return nil
	}

	fmt.Printf("%sからの変更:\n", tag)
	fmt.Println(string(output))
	return nil
}

func ParseVersion(tag string) (Version, error) {
	// v を除去
	version := strings.TrimPrefix(tag, "v")
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return Version{}, errors.New("不正なバージョン形式です")
	}

	var v Version
	var err error

	// 各パートが数字のみで構成されているかチェック
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return Version{}, errors.New("バージョンは数字のみで構成されている必要があります")
		}
	}

	v.Major, err = strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, err
	}

	v.Minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, err
	}

	v.Patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, err
	}

	return v, nil
}

func BumpVersion(tag string, bumpType string) (string, error) {
	v, err := ParseVersion(tag)
	if err != nil {
		return "", err
	}

	switch bumpType {
	case "major":
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case "minor":
		v.Minor++
		v.Patch = 0
	case "patch":
		v.Patch++
	default:
		return "", errors.New("引数は 'major', 'minor', 'patch' のいずれかで指定してください")
	}

	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch), nil
}

func CreateTag(tag string) error {
	// タグの作成
	cmd := exec.Command("git", "tag", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("タグの作成に失敗しました: %v\n出力: %s", err, string(output))
	}

	return nil
}
