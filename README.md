
````markdown
# Go + MySQL 負荷試験用サンプルAPI

## 概要

このプロジェクトは、負荷試験の「的」となる、Go言語製のシンプルなAPIサーバーです。
データベースにはMySQLを使用し、**Docker Compose** を利用することで、コマンド一つでGoサーバーとMySQLデータベースの環境を簡単に構築できます。

負荷試験ツールとして、開発者向けの **k6** と、GUIで操作できる **Apache JMeter** の両方の手順を記載しています。

## 特徴

-   **Go言語**: 高速でリソース効率の良いAPIサーバー
-   **MySQL**: 本番環境を想定したリレーショナルデータベース
-   **Docker**: 使い捨て可能でクリーンな開発・テスト環境
-   **k6 / JMeter**: 目的やスキルに応じて選択できる負荷試験ツール

## 準備 (Prerequisites)

このプロジェクトを実行するには、お使いのPCに以下のソフトウェアがインストールされている必要があります。

-   [**Go** (1.23以上)](https://go.dev/dl/)
-   [**Docker Desktop**](https://www.docker.com/products/docker-desktop/)
-   **負荷試験ツール** (どちらか、または両方)
    -   [**k6**](https://k6.io/docs/getting-started/installation/)
    -   [**Apache JMeter**](https://jmeter.apache.org/download_jmeter.cgi) (+ Java 8以上)

---

## 使い方 (Usage)

### 1. 依存関係の準備
最初に、プロジェクトのルートディレクトリで以下のコマンドを実行し、GoのモジュールとMySQLドライバを準備します。

```bash
# Goモジュールを初期化
go mod init myapp

# MySQLドライバを取得
go get [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
````

### 2\. 環境の起動

以下のコマンドを1行実行するだけで、APIサーバーとMySQLデータベースが起動します。

```bash
docker-compose up --build
```

`API server with MySQL started on :8080` とログに表示されれば成功です。

### 3\. 動作確認

環境が起動したら、別のターミナルを開いて以下の`curl`コマンドでAPIが正しく動作しているか確認できます。

```bash
curl http://localhost:8080/items
```

-----

## 負荷試験の実行

APIサーバーが起動している状態で、以下のいずれかのツールを使って負荷試験を実施します。

### 🧪 オプションA: k6を使う場合 (開発者向け)

`test.js`は、複数のAPIを呼び出す実践的なシナリオが記述されたk6用のスクリプトです。
以下のコマンドで負荷試験を開始します。

```bash
k6 run test.js
```

テストが完了すると、コマンドラインに応答時間やエラーレートなどの結果が表示されます。

-----

### 🧪 オプションB: Apache JMeterを使う場合 (GUIで操作)

1.  **JMeterの起動**:
    JMeterをダウンロード・展開したフォルダ内の`bin`ディレクトリから、`jmeter.sh` (Mac/Linux) または `jmeter.bat` (Windows) を実行してGUIを起動します。

2.  **テスト計画の作成**:

      - **スレッドグループの追加**: `Test Plan` を右クリック → `Add > Threads (Users) > Thread Group`
          - `Number of Threads (users)`: 10
      - **HTTPリクエストの追加**: `Thread Group` を右クリック → `Add > Sampler > HTTP Request`
          - **Server Name or IP**: `localhost`
          - **Port Number**: `8080`
          - **Path**: `/items`
      - **リスナー(結果表示)の追加**: `Thread Group` を右クリック → `Add > Listener > Summary Report`

3.  **テスト実行**:
    画面上部の緑色の再生ボタン ▶ をクリックするとテストが開始されます。`Summary Report`画面で、平均応答時間やエラー率などの統計結果を確認できます。

-----

## データベース負荷の監視

どちらのツールで負荷試験を実行している間も、**実際にRDS(MySQL)の負荷は上昇します**。

AWSマネジメントコンソールなどを開き、**Amazon CloudWatch** や **Performance Insights** で、CPU使用率、DBコネクション数、IOPSなどのメトリクスを監視してください。これにより、どの程度の負荷で性能が劣化するか、どこにボトルネックがあるかを特定できます。

## 環境の停止とクリーンアップ

テストが終わったら、以下のコマンドで起動した全てのコンテナを停止・削除し、クリーンな状態に戻すことができます。

```bash
docker-compose down
```

```
```