package auth

import (
	"strings"
	"testing"
)

func TestJWTAndPasswordHashing(t *testing.T) {
	hash, err := HashPassword("Secreto2026!")
	if err != nil {
		t.Fatalf("error hasheando password: %v", err)
	}
	if !strings.HasPrefix(hash, "$2a$10$") && !strings.HasPrefix(hash, "$2b$10$") {
		t.Fatalf("esperaba hash bcrypt con prefijo $2a$10$, obtuvo: %s", hash)
	}
	if !VerifyPassword("Secreto2026!", hash) {
		t.Fatalf("esperaba verificación exitosa del password")
	}
	if VerifyPassword("Incorrecta", hash) {
		t.Fatalf("no debía verificar un password incorrecto")
	}
	if VerifyPassword("123456", "") {
		t.Fatalf("no debe existir puerta trasera para storedHash vacío")
	}

	seedHashes := []string{
		"$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS",
		"$2a$10$X9j0kHu0e7UmDvcemNqE5ujNhhKppb6coqyKnnkjJKJAMJ2Y8O1D.",
		"$2a$10$DYXz4hs.721eXHdkzox5je3eEtaaxRWe2B3Oq6XfMH4HXbmKr.9fO",
	}
	for i, sh := range seedHashes {
		if !VerifyPassword("123456", sh) {
			t.Fatalf("el hash semilla #%d no verificó la clave '123456'", i+1)
		}
		if VerifyPassword("654321", sh) {
			t.Fatalf("el hash semilla #%d no debía verificar una clave incorrecta", i+1)
		}
	}

	token, err := GenerateToken("user-123", "vecino@correo.com", "client", "")
	if err != nil {
		t.Fatalf("error generando JWT: %v", err)
	}
	claims, err := ParseAndVerifyToken(token)
	if err != nil {
		t.Fatalf("error verificando JWT: %v", err)
	}
	if claims.UserID != "user-123" || claims.Email != "vecino@correo.com" || claims.Role != "client" {
		t.Fatalf("claims inesperados: %+v", claims)
	}
}
