module github.com/mackerelio-labs/mackerel-sql-metric-collector

go 1.24.0

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/config v1.29.9
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.66
	github.com/aws/aws-sdk-go-v2/service/s3 v1.78.2
	github.com/aws/aws-sdk-go-v2/service/ssm v1.57.2
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/stdr v1.2.2
	github.com/go-logr/zapr v1.2.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.0
	github.com/lib/pq v1.9.0
	github.com/mackerelio/mackerel-client-go v0.12.0
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/speee/go-athena v1.0.3
	go.uber.org/zap v1.19.0
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/bigquery v1.0.16
)

require (
	cloud.google.com/go v0.60.0 // indirect
	cloud.google.com/go/bigquery v1.9.0 // indirect
	github.com/aws/aws-sdk-go v1.40.58 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.62 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.17 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/prestodb/presto-go-client v0.0.0-20201204133205-8958eb37e584 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	go.opencensus.io v0.22.3 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.0.0-20200626171337-aa94e735be7f // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/api v0.28.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200626011028-ee7919e894b5 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/gokrb5.v6 v6.1.1 // indirect
	gopkg.in/jcmturner/rpc.v1 v1.1.0 // indirect
)
