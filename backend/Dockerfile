# ベースイメージ
FROM golang:latest

# 作業ディレクトリ
WORKDIR /app/backend

# モジュールのキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# ソースのコピー
COPY . .

# airをインストール
RUN go install github.com/air-verse/air@latest

# 初回のビルド
RUN go build -o main .

# ポート
EXPOSE 8080

# airを起動
CMD ["air"]