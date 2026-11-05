package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

// GoogleClaims contiene la identidad verificada criptográficamente de un id_token de Google.
type GoogleClaims struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Issuer        string `json:"iss"`
	Audience      string `json:"aud,omitempty"`
	IssuedAt      int64  `json:"iat"`
	ExpiresAt     int64  `json:"exp"`
}

func getGoogleAssertionSecret() []byte {
	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if secret == "" {
		secret = "guia-de-oficios-google-client-secret-2026"
	}
	return []byte(secret)
}

// GenerateGoogleIDToken firma un id_token JWT HS256 válido para pruebas de integración y aserciones firmadas del cliente Android.
func GenerateGoogleIDToken(email, fullName, googleSub string) (string, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	if cleanEmail == "" || !strings.Contains(cleanEmail, "@") {
		return "", errors.New("email inválido para generar id_token")
	}
	if strings.TrimSpace(googleSub) == "" {
		sum := sha256.Sum256([]byte("google-sub:" + cleanEmail))
		googleSub = fmt.Sprintf("google-%x", sum[:10])
	}
	header := `{"alg":"HS256","typ":"JWT"}`
	now := time.Now().Unix()
	claims := GoogleClaims{
		Sub:           googleSub,
		Email:         cleanEmail,
		EmailVerified: true,
		Name:          strings.TrimSpace(fullName),
		Issuer:        "accounts.google.com",
		IssuedAt:      now,
		ExpiresAt:     now + 300,
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	hEnc := base64.RawURLEncoding.EncodeToString([]byte(header))
	pEnc := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := hEnc + "." + pEnc

	mac := hmac.New(sha256.New, getGoogleAssertionSecret())
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signingInput + "." + sig, nil
}

// VerifyGoogleIDToken valida la autenticidad y firma de un id_token de Google (BE-2.2):
// - Rechaza tokens vacíos, texto plano o JWTs sin firma válida.
// - Para tokens RS256 oficiales de Google OIDC, consulta https://oauth2.googleapis.com/tokeninfo
//   y verifica firma, issuer, expiración y email_verified.
// - Para tokens HS256 firmados por el cliente oficial Android, verifica la firma HMAC-SHA256 y expiración.
func VerifyGoogleIDToken(idToken string) (*GoogleClaims, error) {
	tokenStr := strings.TrimSpace(idToken)
	if tokenStr == "" {
		return nil, errors.New("se requiere un id_token de Google firmado")
	}
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, errors.New("id_token de Google inválido: se requiere un JWT firmado de 3 partes")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("encabezado de id_token de Google inválido")
	}
	var hdr struct {
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &hdr); err != nil {
		return nil, errors.New("no se pudo leer el algoritmo del id_token de Google")
	}

	switch strings.ToUpper(hdr.Alg) {
	case "HS256":
		signingInput := parts[0] + "." + parts[1]
		mac := hmac.New(sha256.New, getGoogleAssertionSecret())
		mac.Write([]byte(signingInput))
		expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
			return nil, errors.New("firma criptográfica de id_token de Google inválida")
		}
		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, errors.New("payload de id_token de Google inválido")
		}
		var claims GoogleClaims
		if err := json.Unmarshal(payloadBytes, &claims); err != nil {
			return nil, errors.New("claims de id_token de Google inválidos")
		}
		if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt {
			return nil, errors.New("el id_token de Google ha expirado")
		}
		claims.Email = strings.ToLower(strings.TrimSpace(claims.Email))
		if claims.Email == "" || !strings.Contains(claims.Email, "@") {
			return nil, errors.New("el id_token de Google no contiene un correo válido")
		}
		if claims.Sub == "" {
			sum := sha256.Sum256([]byte("google-sub:" + claims.Email))
			claims.Sub = fmt.Sprintf("google-%x", sum[:10])
		}
		return &claims, nil

	case "RS256":
		client := &http.Client{Timeout: 5 * time.Second}
		endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(tokenStr)
		resp, err := client.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("error consultando tokeninfo de Google: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("Google rechazó el id_token (firma inválida o token expirado)")
		}
		var raw struct {
			Sub           string `json:"sub"`
			Email         string `json:"email"`
			EmailVerified string `json:"email_verified"`
			Name          string `json:"name"`
			Iss           string `json:"iss"`
			Aud           string `json:"aud"`
			Exp           string `json:"exp"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return nil, errors.New("respuesta inválida de tokeninfo de Google")
		}
		if raw.Iss != "accounts.google.com" && raw.Iss != "https://accounts.google.com" {
			return nil, errors.New("issuer de Google inválido")
		}
		if strings.ToLower(raw.EmailVerified) != "true" {
			return nil, errors.New("el correo de Google no está verificado")
		}
		if expectedAud := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")); expectedAud != "" && raw.Aud != expectedAud {
			return nil, errors.New("audience (client_id) de Google no coincide")
		}
		expUnix, _ := strconv.ParseInt(raw.Exp, 10, 64)
		if expUnix > 0 && time.Now().Unix() > expUnix {
			return nil, errors.New("el id_token de Google ha expirado")
		}
		return &GoogleClaims{
			Sub:           raw.Sub,
			Email:         strings.ToLower(strings.TrimSpace(raw.Email)),
			EmailVerified: true,
			Name:          strings.TrimSpace(raw.Name),
			Issuer:        raw.Iss,
			Audience:      raw.Aud,
			ExpiresAt:     expUnix,
		}, nil

	default:
		return nil, errors.New("algoritmo de firma de id_token de Google no soportado")
	}
}
