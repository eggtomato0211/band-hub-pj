# ローカル開発用 Dockerfile
# Air によるホットリロードで、コード変更時に自動でサーバーが再起動する
FROM golang:1.22-alpine

WORKDIR /app

# Air（ホットリロードツール）をインストール
RUN go install github.com/air-verse/air@latest

# 依存関係を先にコピーしてキャッシュを効かせる
# go.mod / go.sum が変わらない限り、go mod download の層はキャッシュされる
COPY go.mod go.sum ./
RUN go mod download

# ソースコード全体をコピー
COPY . .

# Air の設定ファイルがあればそれを使い、なければデフォルト設定で起動
CMD ["air", "-c", ".air.toml"]
