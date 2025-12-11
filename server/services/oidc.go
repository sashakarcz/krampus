package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"krampus/server/config"
	"krampus/server/database"
	"krampus/server/models"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var (
	OIDCProvider *oidc.Provider
	OAuth2Config *oauth2.Config
)

// IDToken is an interface for ID tokens (both standard OIDC and HS256)
type IDToken interface {
	Claims(v interface{}) error
}

// InitializeOIDC sets up the OIDC provider and OAuth2 config
func InitializeOIDC(ctx context.Context) error {
	if config.AppConfig.OIDCProviderURL == "" {
		return fmt.Errorf("OIDC provider URL not configured")
	}

	var err error
	OIDCProvider, err = oidc.NewProvider(ctx, config.AppConfig.OIDCProviderURL)
	if err != nil {
		return fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	OAuth2Config = &oauth2.Config{
		ClientID:     config.AppConfig.OIDCClientID,
		ClientSecret: config.AppConfig.OIDCClientSecret,
		RedirectURL:  config.AppConfig.OIDCRedirectURL,
		Endpoint:     OIDCProvider.Endpoint(),
		Scopes:       config.AppConfig.OIDCScopes,
	}

	return nil
}

// GenerateStateToken creates a random state token for OIDC flow
func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetAuthCodeURL generates the authorization URL for OIDC
func GetAuthCodeURL(state string) string {
	return OAuth2Config.AuthCodeURL(state)
}

// ExchangeCodeForToken exchanges authorization code for token
func ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return OAuth2Config.Exchange(ctx, code)
}

// VerifyIDToken verifies the ID token from OIDC provider
func VerifyIDToken(ctx context.Context, rawIDToken string) (IDToken, error) {
	// Parse the token to check the algorithm
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(rawIDToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if it's using HS256
	if token.Method.Alg() == "HS256" {
		// Verify HS256 token manually using client secret
		parsedToken, err := jwt.Parse(rawIDToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.OIDCClientSecret), nil
		})

		if err != nil || !parsedToken.Valid {
			return nil, fmt.Errorf("HS256 token verification failed: %w", err)
		}

		// Extract claims and create a mock IDToken for compatibility
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("failed to get claims from token")
		}

		// Convert claims to IDToken
		claimsJSON, err := json.Marshal(claims)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal claims: %w", err)
		}

		// Return a mockIDToken that implements the IDToken interface
		return &mockIDToken{
			claims: claimsJSON,
		}, nil
	}

	// Use standard OIDC verification for RS256 and other algorithms
	verifier := OIDCProvider.Verifier(&oidc.Config{
		ClientID: config.AppConfig.OIDCClientID,
	})
	return verifier.Verify(ctx, rawIDToken)
}

// mockIDToken implements IDToken interface for HS256 tokens
type mockIDToken struct {
	claims json.RawMessage
}

// Claims unmarshals the claims into the provided structure
func (m *mockIDToken) Claims(v interface{}) error {
	return json.Unmarshal(m.claims, v)
}

// IDTokenClaims represents the claims from an ID token
type IDTokenClaims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

// GetOrCreateUser retrieves or creates a user based on OIDC claims
func GetOrCreateUser(ctx context.Context, idToken IDToken) (*models.User, error) {
	var claims IDTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Try to find existing user by OIDC subject
	var user models.User
	err := database.DB.QueryRow(
		`SELECT id, username, role, oidc_subject, email, created_at, last_login
		 FROM users WHERE oidc_subject = ?`,
		claims.Subject,
	).Scan(&user.ID, &user.Username, &user.Role, &user.OIDCSubject, &user.Email, &user.CreatedAt, &user.LastLogin)

	if err == nil {
		// User exists, update last login
		now := time.Now()
		_, err = database.DB.Exec(
			`UPDATE users SET last_login = ? WHERE id = ?`,
			now, user.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update last login: %w", err)
		}
		user.LastLogin = &now
		return &user, nil
	}

	// User doesn't exist, create new one
	username := claims.Email
	if username == "" {
		username = claims.Name
	}
	if username == "" {
		username = claims.Subject
	}

	// Determine role based on admin email list
	role := string(models.RoleUser)
	if config.AppConfig.IsAdminEmail(claims.Email) {
		role = string(models.RoleAdmin)
	}

	result, err := database.DB.Exec(
		`INSERT INTO users (username, oidc_subject, email, role) VALUES (?, ?, ?, ?)`,
		username, claims.Subject, claims.Email, role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Fetch the newly created user
	err = database.DB.QueryRow(
		`SELECT id, username, role, oidc_subject, email, created_at, last_login
		 FROM users WHERE id = ?`,
		userID,
	).Scan(&user.ID, &user.Username, &user.Role, &user.OIDCSubject, &user.Email, &user.CreatedAt, &user.LastLogin)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch new user: %w", err)
	}

	return &user, nil
}
