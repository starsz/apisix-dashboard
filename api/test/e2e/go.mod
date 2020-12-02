module e2e

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/gavv/httpexpect/v2 v2.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/stretchr/testify v1.4.0
	github.com/tidwall/gjson v1.6.1
	go.etcd.io/etcd v3.3.25+incompatible
)
