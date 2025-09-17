# Go言語の公式イメージを新しいバージョンに変更
FROM golang:1.23-alpine

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# 最初に依存関係のファイルだけコピーして、効率的にライブラリをダウンロード
COPY go.mod ./
COPY go.sum ./

# PROXYのエラーを回避するため、証明書の検証を無効化
ENV GOINSECURE=proxy.golang.org
# さらに、HTTPでの通信を強制する
ENV GOPROXY=http://proxy.golang.org,direct

RUN go mod download

# プロジェクトの全ファイルをコピー
COPY . .

# Goアプリをビルド（コンパイル）
RUN go build -o /app/main .

# コンテナ起動時に実行するコマンド
CMD [ "/app/main" ]