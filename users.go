package amtli

type Credentials struct {
	Name     string
	Password []byte
}

type DecodedToken struct {
	UserName string
}

type AuthenticationResult struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type PasswordPolicy struct {
	RequiredCharsets [][]byte
}

func NewPasswordPolicy(requireDigits bool, requireUppercase bool) {
}

type User struct {
	Name          string
	EmailAddress  string
	EmailVerified bool
	IsDisabled    bool
	PasswordHash  []byte `json:"-"`
}

type RegistrationRequest struct {
	Name                 string
	Email                string
	Password             []byte
	PasswordConfirmation []byte
}
