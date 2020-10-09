package base

type ResponseCode int

const (
	ResponseOK                 ResponseCode = 1000
	ResponseValidationError    ResponseCode = 2000
	ResponseRequestParamsError ResponseCode = 2100
	ResponseInnerServerError   ResponseCode = 5000
	ResponseBusyError          ResponseCode = 6000
)

type Code struct {
	Val int
	Msg string
}

type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
}
