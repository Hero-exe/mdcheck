# mdcheck

`mdcheck` は、ブログ記事向けの Markdown チェッカーです。

公開前の記事に対して、見出し構造、空リンク、画像パス、文字数、YAML front matter のメタ情報を確認できます。

## チェック内容

- 見出し構造
  - H1 がない
  - H1 が複数ある
  - 空の見出しがある
  - 見出しレベルが飛んでいる
- 空リンク
  - リンクテキストが空
  - リンク先 URL が空
- 画像パス
  - ローカル画像ファイルが存在しない
  - 画像パスが空
- 文字数
  - 本文の文字数を表示
  - 設定された最小/最大文字数の範囲外を検出
- メタ情報
  - YAML front matter の必須項目を確認

## 使い方

単一ファイルをチェックします。

```sh
go run ./cmd/mdcheck article.md
```

ディレクトリ配下の Markdown ファイルをまとめてチェックします。

```sh
go run ./cmd/mdcheck posts/
```

JSON 形式で出力します。

```sh
go run ./cmd/mdcheck --format json posts/
```

設定ファイルを指定します。

```sh
go run ./cmd/mdcheck --config mdcheck.yaml posts/
```

## 設定

作業ディレクトリに `mdcheck.yaml` を作成すると、ルールの重要度や必須メタ情報を変更できます。

```yaml
rules:
  heading_structure: error
  empty_link: error
  image_path: error
  word_count: warn
  metadata: warn

metadata:
  required:
    - title
    - description
    - date
    - tags

word_count:
  min: 800
  max: 5000

ignore:
  - drafts/
  - node_modules/
```

ルールの重要度には、次の値を指定できます。

- `error`
- `warn`
- `info`
- `off`

## 終了コード

- `0`: error の検出なし
- `1`: error の検出あり
- `2`: CLI の使い方、または実行時エラー

## 開発

テストを実行します。

```sh
go test ./...
```

CLI を直接実行します。

```sh
go run ./cmd/mdcheck README.md
```

## 今後の予定

今後の実装候補は [plan/README.md](plan/README.md) にまとめています。
