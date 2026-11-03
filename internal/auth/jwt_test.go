package auth

import (
	"testing"
)

func TestJWTAndPasswordHashing(t *testing.T) {
	hash, err := HashPassword("Secreto2026!")
	if err != nil {
		t.Fatalf("error hasheando password: %v", err)
	}
	if !VerifyPassword("Secreto2026!", hash) {
		t.Fatalf("esperaba verificación exitosa del password")
	}
	if VerifyPassword("Incorrecta", hash) {
		t.Fatalf("no debía verificar un password incorrecto")
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
