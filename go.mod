module das-account-indexer

go 1.16

require (
	github.com/dotbitHQ/das-lib v1.0.1-0.20230222083231-42a3107914f5
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.8.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/nervosnetwork/ckb-sdk-go v0.101.3
	github.com/parnurzeal/gorequest v0.2.16
	github.com/scorpiotzh/mylog v1.0.10
	github.com/scorpiotzh/toolib v1.1.4
	github.com/urfave/cli/v2 v2.3.0
	gorm.io/gorm v1.23.6
)

replace github.com/ethereum/go-ethereum v1.9.14 => github.com/ethereum/go-ethereum v1.10.17
