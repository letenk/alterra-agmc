package lib

type ResponseWithData struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ApiResponseWithData(code int, status string, message string, data interface{}) ResponseWithData {
	jsonResponse := ResponseWithData{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}

	return jsonResponse
}
