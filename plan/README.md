# mdcheck 実装計画

このドキュメントは、今後 `mdcheck` を別のエージェントや開発者が実装・拡張するときの引き継ぎ用メモです。

`mdcheck` は、ブログ記事の公開前チェックを目的とした Go 製 CLI です。Markdown 記事に対して、見出し構造、空リンク、画像パス、文字数、front matter のメタ情報を確認します。

## 現状

- Go 標準ライブラリのみで初期 CLI を実装済み。
- エントリーポイント: `cmd/mdcheck/main.go`
- CLI 全体の制御: `internal/app/run.go`
- Markdown 解析: `internal/markdown`
- チェックルール: `internal/rules`
- 出力処理: `internal/report`
- 設定読み込み: `internal/config`
- Markdown 解析と主要ルールのテストを追加済み。
- module path は `github.com/Hero-exe/mdcheck`。

動作確認:

```sh
go test ./...
go run ./cmd/mdcheck README.md
go run ./cmd/mdcheck --format json README.md
```

## 公開リポジトリとしての注意点

- `.cache/`、`outputs/`、`.DS_Store`、ビルド済みバイナリはコミットしない。
- 秘密情報、個人用設定、トークン類を追加しない。
- まだ `LICENSE` がない。公開OSSとして見せる前に追加する。

## 設計方針

- ルールは独立させ、追加しやすくする。
- まずは分かりやすい検出結果を優先する。
- 外部依存は、設定読み込みや Markdown 解析の品質を明確に上げる場合に追加する。
- CI で使えるように、終了コードと JSON 出力は安定させる。
- text 出力は人間向け、JSON 出力は自動化向けとして扱う。

## ルール追加の基本

新しいチェックは `internal/rules` 配下に実装する。

各ルールでやること:

- `Name() string` を実装する。
- `Check(ctx Context, doc markdown.Document) []Finding` を実装する。
- 可能な限り正確な行番号を返す。
- ルール内で `os.Exit` や標準出力への直接出力をしない。
- `internal/rules/rule.go` の `DefaultRules()` に追加する。
- ルール単位のテストを追加する。

## 優先タスク

### 1. LICENSE を追加する

目的:
public リポジトリとして再利用しやすい状態にする。

実装案:

- `LICENSE` を追加する。
- 小さな CLI なので、特に希望がなければ MIT License が扱いやすい。
- `README.md` にライセンス欄を追加する。

受け入れ条件:

- `LICENSE` が存在する。
- `README.md` からライセンスが分かる。

### 2. public 向けに README を拡充する

目的:
GitHub を見た人が、インストール、実行、設定まで迷わず試せるようにする。

実装案:

- インストール手順を追加する。
  - `go install github.com/Hero-exe/mdcheck/cmd/mdcheck@latest`
  - ローカル開発時の実行方法
- 出力例を追加する。
- 設定ファイル例を増やす。
- 終了コードの説明を補足する。
- 短いロードマップを追加する。

受け入れ条件:

- README だけを読めば、インストール、実行、設定ができる。

### 3. GitHub Actions CI を追加する

目的:
push と pull request で自動テストを走らせる。

実装案:

- `.github/workflows/ci.yml` を追加する。
- `go test ./...` を実行する。
- 安定版の Go を使う。
- 可能なら `gofmt` チェックも入れる。

受け入れ条件:

- `push` と `pull_request` で CI が走る。
- テスト失敗時に CI が失敗する。
- gofmt されていない場合に CI が失敗する。

### 4. 設定ファイル読み込みを本物の YAML パーサに置き換える

目的:
`mdcheck.yaml` をより自然に書けるようにする。

現状の制約:
`internal/config` は小さな行ベースのパーサで、現在のサンプルは読めるが完全な YAML パーサではない。

実装案:

- `gopkg.in/yaml.v3` など、メンテナンスされている YAML パーサを追加する。
- 現在の設定形式は維持する。
- 設定ファイルがない場合はデフォルト値を使う。
- 次のテストを追加する。
  - 設定ファイルなし
  - ルール重要度の上書き
  - 必須メタ情報の上書き
  - 文字数の min/max
  - ignore pattern

受け入れ条件:

- 既存の設定例がそのまま動く。
- コメントやリストを含む、より自然な YAML を扱える。

### 5. Markdown 解析を強化する

目的:
誤検出を減らし、より多くの Markdown 記法に対応する。

現状の制約:
`internal/markdown` は正規表現と行走査で実装されている。

実装案:

- `github.com/yuin/goldmark` の AST 利用を検討する。
- 可能であれば、現在の `markdown.Document` を内部契約として維持する。
- 次の記法に対応する。
  - reference-style link
  - autolink
  - title 付き Markdown link
  - title 付き image
  - ネストした emphasis を含む link

受け入れ条件:

- 既存テストが通る。
- reference-style link と image title のテストが追加されている。
- fenced code block 内は引き続き無視される。

### 6. 画像 alt テキストルールを追加する

目的:
alt のない画像を検出し、ブログ記事の読みやすさとアクセシビリティを上げる。

実装案:

- `image_alt` ルールを追加する。
- `![](path.png)` と `![   ](path.png)` を検出する。
- デフォルト重要度は `warn` にする。
- 既存の severity 設定で `off` にできるようにする。

受け入れ条件:

- 空の alt text が finding になる。
- alt text がある画像は finding にならない。
- `image_alt: off` で無効化できる。

### 7. タイトル長チェックを追加する

目的:
front matter の `title` が長すぎないか確認できるようにする。

実装案:

- `title_length` ルールを追加する。
- front matter の `title` をチェックする。
- 設定を追加する。
  - `title_length.min`
  - `title_length.max`
- デフォルトは max 60 文字程度を候補にする。

受け入れ条件:

- `title` がない場合は既存の metadata ルールが担当する。
- 長すぎる title が finding になる。
- 設定で最大文字数を変更できる。

### 8. 内部リンクチェックを追加する

目的:
記事内のローカルリンク切れを検出する。

実装案:

- `internal_link` ルールを追加する。
- remote URL は無視する。
- 相対パスは現在の Markdown ファイルのディレクトリ基準で解決する。
- `#section` のようなページ内アンカーを扱う。
- `file.md#section` のようなリンクは、まずファイルの存在を確認する。

受け入れ条件:

- 存在しないローカルファイルリンクが finding になる。
- remote URL は無視される。
- 存在するファイルへのリンクは許可される。

### 9. 終了コードまわりのテストを追加する

目的:
CI での利用に必要な挙動を安定させる。

実装案:

- `internal/app.Run` のテストを追加する。
- 次を確認する。
  - error finding がなければ nil を返す。
  - error finding があれば `ExitCodeError{Code: 1}` を返す。
  - 不正な format は通常の error を返す。
  - JSON 出力が valid JSON になる。

受け入れ条件:

- `os.Exit` を呼ばずに app-level の挙動をテストできる。

### 10. リリース準備をする

目的:
CLI をバイナリとして導入しやすくする。

実装案:

- GoReleaser 設定、またはシンプルな GitHub Actions の release workflow を追加する。
- 次のバイナリを作る。
  - macOS arm64
  - macOS amd64
  - Linux amd64
  - Linux arm64
- checksums を出力する。

受け入れ条件:

- tag を打つとダウンロード可能なバイナリが作られる。
- README にインストール方法が書かれている。

## あると便利なタスク

- `--quiet` を追加する。
- `--fail-on warn` を追加する。
- `--no-color` と色付き text 出力を追加する。
- `--ignore` CLI flag を追加する。
- `--rule` CLI flag で一部ルールだけ実行できるようにする。
- GitHub code scanning 向けに SARIF 出力を追加する。
- `mdcheck init` で `mdcheck.yaml` を生成する。

## 最初に切るとよさそうな Issue

最初の別エージェント向けタスクとしては、次の順番が扱いやすい。

1. `LICENSE` を追加する。
2. `README.md` を public 向けに拡充する。
3. GitHub Actions CI を追加する。
4. `image_alt` ルールを追加する。

これらは互いに独立していて、リスクが低く、public リポジトリとしての見栄えと実用性をすぐに上げられる。
