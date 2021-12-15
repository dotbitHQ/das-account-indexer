package toolib

import "encoding/json"

func JsonString(o interface{}) string {
	s, _ := json.Marshal(o)
	return string(s)
}
