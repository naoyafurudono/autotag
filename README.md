# autotag

Gitのバージョンタグを自動的に管理するツールです。

## 機能

- 最新のGitタグを取得
- バージョン番号の自動更新（major/minor/patch）
- 新しいタグの作成とプッシュ

## インストール

```bash
go install github.com/naoyafurudono/autotag@latest
```

## 使用方法

```bash
# パッチバージョンを上げる（デフォルト）
autotag

# マイナーバージョンを上げる
autotag minor

# メジャーバージョンを上げる
autotag major
```

## 例

現在のタグが `v1.2.3` の場合：

```bash
# v1.2.4 が作成される
autotag

# v1.3.0 が作成される
autotag minor

# v2.0.0 が作成される
autotag major
```

## 注意事項

- タグが存在しない場合は `v0.0.0` から開始します
- バージョン番号は `vX.Y.Z` の形式である必要があります
- バージョン番号は数字のみで構成されている必要があります 
