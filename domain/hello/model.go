package hello

type Model struct {
	Message string `json:"message"`
}

type HelloResponse struct {
	Body Model
}
