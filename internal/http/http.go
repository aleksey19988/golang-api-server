package http

type Response struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   string `json:"data"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)
