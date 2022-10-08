package response

type Response struct {
	ErrorCode    int         `json:"errCode"`
	ErrorMessage string      `json:"errMessage,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}
