package hello

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (s *Repository) GetMessage() Model {
	return Model{Message: "Hello World"}
}
