module opTools/getES

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/google/gops v0.3.11
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/olivere/elastic/v7 v7.0.20
	go.uber.org/zap v1.16.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace common => ../common
