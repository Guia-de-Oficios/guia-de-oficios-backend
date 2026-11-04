package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Claims representa el payload del token JWT de Guía de Oficios.
type Claims struct {
	UserID         string `json:"sub"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	ProfessionalID string `json:"professional_id,omitempty"`
	IssuedAt       int64  `json:"iat"`
	ExpiresAt      int64  `json:"exp"`
}

func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "guia-de-oficios-production-secret-key-2026"
	}
	return []byte(secret)
}

// GenerateToken firma un JWT HS256 estándar con validez de 30 días.
func GenerateToken(userID, email, role, professionalID string) (string, error) {
	header := `{"alg":"HS256","typ":"JWT"}`
	now := time.Now().Unix()
	claims := Claims{
		UserID:         userID,
		Email:          email,
		Role:           role,
		ProfessionalID: professionalID,
		IssuedAt:       now,
		ExpiresAt:      now + int64(30*24*3600),
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	hEnc := base64.RawURLEncoding.EncodeToString([]byte(header))
	pEnc := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := hEnc + "." + pEnc

	mac := hmac.New(sha256.New, getSecret())
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig, nil
}

// ParseAndVerifyToken valida la firma HS256 y la expiración del JWT.
func ParseAndVerifyToken(tokenStr string) (*Claims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("formato de token JWT inválido")
	}

	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, getSecret())
	mac.Write([]byte(signingInput))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, errors.New("firma de token JWT inválida")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("error decodificando payload JWT")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("error parseando claims JWT")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("el token JWT ha expirado")
	}

	return &claims, nil
}

// BcryptCost define el costo estándar (>= 10) para el hashing bcrypt de contraseñas.
const BcryptCost = bcrypt.DefaultCost

// HashPassword genera un hash seguro de la contraseña utilizando bcrypt (costo >= 10).
func HashPassword(password string) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", errors.New("la contraseña no puede estar vacía")
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// VerifyPassword comprueba si una contraseña en texto plano coincide con su hash bcrypt almacenado.
// No admite contraseñas sin hash ni puertas traseras para hashes vacíos.
func VerifyPassword(password, storedHash string) bool {
	if strings.TrimSpace(storedHash) == "" || strings.TrimSpace(password) == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}
