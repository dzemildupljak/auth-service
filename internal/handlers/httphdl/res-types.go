package httphdl

type ResponsePayload struct {
	Status  bool        `json:"status"`
	Payload interface{} `json:"payload"`
}

type ResponseWithMessage struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
