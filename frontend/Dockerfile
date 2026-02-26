# ローカル開発用 Dockerfile
# next dev によるホットリロードで、コード変更が即座にブラウザに反映される
FROM node:20-alpine

WORKDIR /app

# 依存関係を先にコピーしてキャッシュを効かせる
COPY package.json package-lock.json ./
RUN npm ci

# ソースコード全体をコピー
COPY . .

# Next.js の開発サーバーを起動
# ポート 3000 で待ち受ける
CMD ["npm", "run", "dev"]
