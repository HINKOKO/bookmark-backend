package main

import (
	"bookmarks/internal/models"
	"context"
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
	ID       string `json:"id"`
	Username string `json:"username"`
}

// TokenPairs - gather token and refresh token
type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refreshy_token_of_devil"`
}

// Claims - wrapper type around the jwt registered claims
type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

// GenerateTokenPair - generate the token pair
func (j *Auth) GenerateTokenPair(userID int) (TokenPairs, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			Audience:  jwt.ClaimStrings{j.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	// Create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			Audience:  jwt.ClaimStrings{j.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	// Create a refresh token - set claims - kinda parallel/similar methods here
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}
	// log.Println("successfull signed token ? => \t", signedAccessToken)
	// log.Println("successfull refresh signed token ? => \t", signedRefreshToken)

	// Finally Return TokenPairs
	return TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil

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
		MaxAge:   60,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true, // No javascript access at all to this cookie
		Secure:   false,
	}
}

func (j *Auth) GetTokenFromCookieAndVerify(tokenString string) (string, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})
	if err != nil || !token.Valid {
		return "", nil, err
	}
	return tokenString, claims, nil
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

// Middleware to add userID to context
func (app *application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No Authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// token := headerParts[1]
		_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
