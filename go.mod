module das-account-indexer

go 1.15

require (
	github.com/DeAccountSystems/das-lib v0.0.0-20220531040850-ea68b195348f
	github.com/elazarl/goproxy v0.0.0-20211114080932-d06c3be7c11b // indirect
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/nervosnetwork/ckb-sdk-go v0.101.3
	github.com/parnurzeal/gorequest v0.2.16
	github.com/scorpiotzh/mylog v1.0.9
	github.com/scorpiotzh/toolib v1.1.3
	github.com/urfave/cli/v2 v2.3.0
	gorm.io/gorm v1.22.4
	moul.io/http2curl v1.0.0 // indirect
)

replace github.com/ethereum/go-ethereum v1.9.14 => github.com/ethereum/go-ethereum v1.10.17
