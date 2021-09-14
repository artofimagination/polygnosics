package http

type ResponseData struct {
	Error string      `json:"error" validation:"required"`
	Data  interface{} `json:"data" validation:"required"`
}
