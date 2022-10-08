package response

type Response struct {
	ErrorCode    int         `json:"err_code,omitempty"`
	ErrorMessage string      `json:"err_message,omitempty"`
	Data         interface{} `json:"data"`
}
