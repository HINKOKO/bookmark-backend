package main

import (
	"bookmarks/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

// Data about an user to issue a token
type jwtUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// TokenPairs - gather token and refresh token
type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Claims - wrapper type around the jwt registered claims
type Claims struct {
	jwt.RegisteredClaims
}

// GenerateTokenPair - generate the token pair
func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims - what is this token claims to be ?
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s", user.Username)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["audience"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// set an expiry for jwt token
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create a refresh token - set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set expiry for refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Finally Return TokenPairs
	return tokenPairs, nil
}

// GetRefreshCookie -
func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true, // -> Giving No javascript access at all to this cookie
		Secure:   true,
	}
}

// GetExpiredRefreshCookie - Function when we want the refresh cookie to be deleted from the user's agent
// How you delete cookie -> Set another cookie with same attribute, but you set its max age to minus one
// and expires time Unix zero.
func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true, // No javascript access at all to this cookie
		Secure:   true,
	}
}

func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// Good practice to add header
	w.Header().Add("Vary", "Authorization")

	// get auth header
	authHeader := r.Header.Get("Authorization")
	// Sanity checks
	if authHeader == "" {
		return "", nil, errors.New("no auth header in header")
	}

	// Split the header, to check for 'bearer'
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header format")
	}
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("incorrect authorization format")
	}

	token := headerParts[1]
	// Declare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}
	// Do we actually issue that token ?
	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("invalid user")
	}

	return token, claims, nil
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthenticating struct {
	Username    string
	Email       string
	Password    string
	Verified    bool
	VerifyToken string
}

func (app *application) CheckPassword(u models.User, plainPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPass))
}
