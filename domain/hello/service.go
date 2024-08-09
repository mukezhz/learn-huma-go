package hello

// Service handles the business logic of the module
type Service struct {
	// Add any dependencies here
	repo *Repository
}

// NewService creates a new instance of TestService
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetMessage returns a greeting message
func (s *Service) GetMessage() Model {
	return s.repo.GetMessage()
}
