module github.com/chuan-fu/Common

go 1.16

// github.com/zeromicro/go-zero

replace (
	github.com/tidwall/gjson => github.com/tidwall/gjson v1.9.3
	github.com/tidwall/match => github.com/tidwall/match v1.0.3
)

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/ahmetb/go-linq v2.0.0-rc0+incompatible
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129
	github.com/apache/rocketmq-client-go/v2 v2.0.0
	github.com/bytedance/sonic v1.4.0
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang/mock v1.6.0 // indirect
	github.com/jinzhu/configor v1.2.1
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/panjf2000/ants/v2 v2.5.0
	github.com/pkg/errors v0.9.1
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f
	github.com/robfig/cron v1.2.0
	github.com/rs/zerolog v1.26.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.7.4
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20220421235706-1d1ef9303861 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/time v0.0.0-20220411224347-583f2d630306
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.3.2
	gorm.io/gorm v1.23.1
)
