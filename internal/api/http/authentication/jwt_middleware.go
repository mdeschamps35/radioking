package authentication

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// KeycloakJWKS represents the structure of Keycloak public keys
type KeycloakJWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a single JSON Web Key
type JWK struct {
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
	X5t string   `json:"x5t"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// KeycloakClaims represents JWT claims from Keycloak
type KeycloakClaims struct {
	jwt.RegisteredClaims
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	Name              string `json:"name"`
}

// JWTMiddleware handles JWT authentication
type JWTMiddleware struct {
	keycloakURL string
	realm       string
	publicKeys  map[string]*rsa.PublicKey
}

const (
	userClaimsContextKey = "user_claims"
	bearerPrefix         = "Bearer "
)

// NewJWTMiddleware creates a new JWT middleware instance
func NewJWTMiddleware(keycloakURL, realm string) *JWTMiddleware {
	return &JWTMiddleware{
		keycloakURL: keycloakURL,
		realm:       realm,
		publicKeys:  make(map[string]*rsa.PublicKey),
	}
}

// LoadPublicKeys loads public keys from Keycloak
func (j *JWTMiddleware) LoadPublicKeys() error {
	jwks, err := j.fetchJWKS()
	if err != nil {
		return err
	}

	signingKeys := j.filterSigningKeys(jwks.Keys)
	if len(signingKeys) == 0 {
		return fmt.Errorf("no valid signing keys found")
	}

	for _, key := range signingKeys {
		publicKey, err := j.parsePublicKey(key)
		if err != nil {
			log.Printf("Error parsing key %s: %v", key.Kid, err)
			continue
		}
		j.publicKeys[key.Kid] = publicKey
	}

	if len(j.publicKeys) == 0 {
		return fmt.Errorf("no valid public keys loaded")
	}

	log.Printf("Public keys loaded: %d", len(j.publicKeys))
	return nil
}

// fetchJWKS retrieves JWKS from Keycloak
func (j *JWTMiddleware) fetchJWKS() (*KeycloakJWKS, error) {
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", j.keycloakURL, j.realm)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks KeycloakJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	return &jwks, nil
}

// filterSigningKeys filters keys that can be used for signature verification
func (j *JWTMiddleware) filterSigningKeys(keys []JWK) []JWK {
	var signingKeys []JWK
	for _, key := range keys {
		if key.Use == "sig" && key.Kty == "RSA" && len(key.X5c) > 0 {
			signingKeys = append(signingKeys, key)
		}
	}
	return signingKeys
}

// parsePublicKey converts a JWK to RSA public key
func (j *JWTMiddleware) parsePublicKey(key JWK) (*rsa.PublicKey, error) {
	certPEM := fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", key.X5c[0])
	return jwt.ParseRSAPublicKeyFromPEM([]byte(certPEM))
}

// Middleware returns JWT authentication middleware
func (j *JWTMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := j.extractBearerToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := j.parseAndValidateToken(tokenString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractBearerToken extracts and validates the Bearer token from Authorization header
func (j *JWTMiddleware) extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", fmt.Errorf("empty bearer token")
	}

	return token, nil
}

// parseAndValidateToken parses and validates the JWT token
func (j *JWTMiddleware) parseAndValidateToken(tokenString string) (*KeycloakClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, j.getSigningKey)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*KeycloakClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// getSigningKey returns the appropriate signing key for token validation
func (j *JWTMiddleware) getSigningKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing kid in token header")
	}

	publicKey, exists := j.publicKeys[kid]
	if !exists {
		return nil, fmt.Errorf("public key not found for kid: %s", kid)
	}

	return publicKey, nil
}

// GetUserClaims extracts user claims from request context
func GetUserClaims(r *http.Request) (*KeycloakClaims, bool) {
	claims, ok := r.Context().Value(userClaimsContextKey).(*KeycloakClaims)
	return claims, ok
}
