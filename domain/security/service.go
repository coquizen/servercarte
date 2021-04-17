package security

type Service interface {
	ConfirmationChecker(password string, confirmPassword string) error
	VerifyPasswordMatches(hashedPW string, password string) error
	Hash(password string) string
	IsValid(password string) error
}

type security struct {
	Service
}

// NewService returns a new receiver for Security
func NewService(secSvc Service) *security {
	return &security{secSvc}
}

