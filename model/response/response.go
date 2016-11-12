package response

type Response struct {
	Message string
	Code    int
	Data    interface{}
	Next    string `json:"Next,omitempty"`
}
