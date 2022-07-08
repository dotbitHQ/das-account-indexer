package code

type Code = int

const (
	ApiCodeSuccess        Code = 0
	ApiCodeError500       Code = 500
	ApiCodeParamsInvalid  Code = 10000
	ApiCodeMethodNotExist Code = 10001
	ApiCodeDbError        Code = 10002

	ApiCodeAccountFormatInvalid Code = 20006
	ApiCodeAccountNotExist      Code = 20007
	ApiCodeAccountOnLock        Code = 20008
)
