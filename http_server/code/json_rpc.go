package code

import "encoding/json"

type JsonRequest struct {
	ID      interface{}     `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type JsonResponse struct {
	ID      interface{} `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

func (j *JsonResponse) ResultData(data interface{}) {
	j.Result = data
}
