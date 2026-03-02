# ローカル開発用 Dockerfile
# Air によるホットリロードで、コード変更時に自動でサーバーが再起動する
# go.mod の go ディレクティブ（1.25.6）に合わせてベースイメージを選択する
FROM golang:1.25-alpine

WORKDIR /app

# Air（ホットリロードツール）をインストール
RUN go install github.com/air-verse/air@latest

# 依存関係を先にコピーしてキャッシュを効かせる
# go.mod / go.sum が変わらない限り、go mod download の層はキャッシュされる
# go.sum がまだ存在しない場合も考慮し、個別にコピーする
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# ソースコード全体をコピー
COPY . .

# Air の設定ファイルがあればそれを使い、なければデフォルト設定で起動
CMD ["air", "-c", ".air.toml"]
