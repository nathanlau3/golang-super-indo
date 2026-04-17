package common

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type Meta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func Success(status int, message string, data interface{}) Response {
	return Response{Status: status, Message: message, Data: data}
}

func Error(status int, message string) Response {
	return Response{Status: status, Message: message}
}

func Paginated(status int, message string, data interface{}, meta Meta) PaginatedResponse {
	return PaginatedResponse{Status: status, Message: message, Data: data, Meta: meta}
}
