module app

go 1.13

require (
	github.com/aws/aws-sdk-go v1.38.51
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/ezamriy/gorpm v0.0.0-20160905202458-25f7273cbf51
	github.com/getkin/kin-openapi v0.9.0
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/validator/v10 v10.8.0 // indirect
	github.com/go-redis/redis/v8 v8.11.4
	github.com/gorilla/websocket v1.4.2
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/lestrrat-go/backoff v1.0.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.0
	github.com/redhatinsights/app-common-go v1.5.1
	github.com/redhatinsights/platform-go-middlewares v0.8.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/ugorji/go v1.2.6 // indirect
	github.com/zsais/go-gin-prometheus v0.1.0
	golang.org/x/net v0.0.0-20211104170005-ce137452f963 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/driver/postgres v1.3.1
	gorm.io/gorm v1.23.1
)

replace github.com/ezamriy/gorpm v0.0.0-20160905202458-25f7273cbf51 => github.com/MichaelMraka/gorpm v0.0.0-20210923131407-e21b5950f175
