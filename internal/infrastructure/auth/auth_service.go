package auth

type AuthService struct {
	JwtService  JWTService
	UserService *UserService
}

func NewAuthService(jwtService JWTService, userService *UserService) *AuthService {
	return &AuthService{jwtService, userService}
}
