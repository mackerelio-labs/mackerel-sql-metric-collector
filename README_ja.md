# mackerel-sql-metric-collector


DB からメトリックを収集して Mackerel のサービスメトリックとして投稿する CLI です。

## 使い方

```console
./bin/mackerel-sql-metric-collector \
  --dsn="ssm://PARAMETER_NAME?withDecryption=true" \
  --mackerel-apikey="ssm://PARAMETER_NAME?withDecryption=true" \
  --default-service="myapp" \
  --query-file "s3://BUCKET/KEY"
```

- `--dsn` は環境変数 `DSN` でも設定可能です
- `--mackerel-apikey` は環境変数 `MACKEREL_APIKEY` でも設定可能です
- `--default-service` は環境変数 `DEFAULT_SERVICE` でも設定可能です
- `--query-file` を指定しない場合は標準入力からクエリ設定 (YAML) を読み込みます

### オプションのデータソース

`--query-file` では **YAML のデータソースとして** 以下の形式をサポートします。
`--query-file` 以外では **オプション値のデータソースとして** 以下の形式をサポートします。

- `STRING`: 文字列
- `file:///PATH/TO/FILE`: ファイル
- `s3://BUCKET/KEY?regisonHint=REGION`: s3
- `ssm://NAME?withDecryption=(true|false)`: パラメータストア

### DSN 例

- PostgreSQL
  - `postgres://<CONNECTION_STRING>`
  - `postgres://host=localhost port=5432 user=myuser password=mypassword dbname=myapp sslmode=disable`
  - see. <https://github.com/lib/pq/blob/master/README.md>
- MySQL
  - `mysql://<CONNECTION_STRING>`
  - `mysql://user:password@/dbname`
  - see. <https://github.com/go-sql-driver/mysql#dsn-data-source-name>
- SQLite3
  - `sqlite3://<CONNECTION_STRING>`
  - `sqlite3://file:test.db?cache=shared&mode=memory`
  - see. <https://github.com/mattn/go-sqlite3#connection-string>
- Athena
  - `athena://<CONNECTION_STRING>`
  - `athena://db=alb&region=ap-northeast-1&output_location=s3://athena-results`
  - see. <https://github.com/speee/go-athena/blob/master/README.md>
- BigQuery
  - `bigquery://<CONNECTION_STRING>`
  - `bigquery://project/location/dataset`
  - see <https://github.com/go-gorm/bigquery#readme>

## クエリ設定例

```yaml
---
- keyPrefix: "users"
  valueKey:
    "count": "user_num" # users.count メトリックとして user_num の値を使用します
  sql: |-
    SELECT
      COUNT(id) AS user_num
    FROM
      users
    WHERE
      created_at >= current_timestamp - INTERVAL '30 DAYS'
- keyPrefix: "users"
  service: "other_service" # 投稿先のサービスを --default-service とは別にしたい場合に定義します
  valueKey:
    "status.#{status}": "user_num" # クエリ結果をメトリック名に使用します
  sql: |-
    SELECT
      status,
      COUNT(id) AS user_num,
    FROM
      users
    WHERE
      is_admin = $1
      AND users.updated_at >= current_timestamp - INTERVAL '30 DAYS'
    GROUP BY status
  params: # プレースホルダのパラメータ値を指定します
    - false
- keyPrefix: "users"
  valueKey:
    "count.#{status}": "user_num"
  defaultValue:
    "count.pending": 0 # 値がなかったときに投稿するメトリックを指定します
  sql: |-
    SELECT
      COUNT(id) AS user_num,
      status
    FROM
      users
    WHERE
      created_at >= current_timestamp - INTERVAL '30 DAYS'
    GROUP BY status
```

## コンテナイメージの取得方法

Docker Hub、Amazon ECR Public Gallery、GitHub Packages Container registry にて公開しております。以下のようなコマンドでコンテナイメージを取得することができます。

```
docker pull mackerel/mackerel-sql-metric-collector:latest
```

```
docker pull public.ecr.aws/mackerel/mackerel-sql-metric-collector:latest
```

```
docker pull ghcr.io/mackerelio-labs/mackerel-sql-metric-collector:latest
```

