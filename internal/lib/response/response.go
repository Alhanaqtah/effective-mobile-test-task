package response

const (
	StatusOK  = "OK"
	StatusErr = "Error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Ok(msg string) Response {
	return Response{
		Status:  StatusOK,
		Message: msg,
	}
}

func Err(errMsg string) Response {
	return Response{
		Status: StatusErr,
		Error:  errMsg,
	}
}
