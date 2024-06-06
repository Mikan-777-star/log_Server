# ベースイメージ
FROM golang:1.20-alpine

# 作業ディレクトリを設定
WORKDIR /app


COPY . ./

RUN [ -f go.mod ] && rm go.mod || true
RUN [ -f go.sum ] && rm go.sum || true
# 必要なファイルをコピー
RUN go mod init logserver
RUN go mod tidy
RUN go mod download


# アプリケーションをビルド
RUN go build -o /logserver

# アプリケーションを実行
CMD ["/logserver"]
