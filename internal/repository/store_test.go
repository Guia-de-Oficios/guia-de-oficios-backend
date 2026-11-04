package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPostgresStoreLiveIntegration(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://guia_user:guia_password@localhost:5432/guia_oficios?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("saltando test de integración SQL: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Skipf("PostgreSQL no disponible en %s: %v", dsn, err)
	}

	store := NewPostgresStore(db)

	info, err := store.Ping(ctx)
	if err != nil {
		t.Fatalf("error en store.Ping: %v", err)
	}
	t.Logf("Conectado exitosamente a: %s", info)

	cats, err := store.ListCategories(ctx)
	if err != nil || len(cats) < 8 {
		t.Fatalf("esperaba al menos 8 categorías en PostgreSQL, obtuvo %d (err=%v)", len(cats), err)
	}

	pros, err := store.ListProfessionals(ctx, ProfessionalFilter{
		Lat: -34.6186,
		Lng: -58.4368,
	})
	if err != nil || len(pros) == 0 {
		t.Fatalf("esperaba profesionales desde PostGIS, obtuvo %d (err=%v)", len(pros), err)
	}
	t.Logf("Profesionales encontrados con PostGIS: %d (Primero: %s a %.2f km)", len(pros), pros[0].FullName, pros[0].DistanceKm)

	detail, err := store.GetProfessionalByID(ctx, "aaaa1111-1111-1111-1111-111111111111")
	if err != nil {
		t.Fatalf("error obteniendo perfil de Roberto en PostgreSQL: %v", err)
	}
	if detail.FullName != "Roberto Carlos Gómez" {
		t.Fatalf("esperaba Roberto Carlos Gómez, obtuvo %s", detail.FullName)
	}

	dash, err := store.GetProfessionalDashboard(ctx, "aaaa1111-1111-1111-1111-111111111111")
	if err != nil {
		t.Fatalf("error obteniendo dashboard profesional en PostgreSQL: %v", err)
	}
	t.Logf("Dashboard profesional OK: %d solicitudes activas", len(dash.ActiveRequests))

	// Verificar autenticación real con hash bcrypt en usuarios semilla
	robertoUser, err := store.AuthenticateUser(ctx, "roberto.gomez@guiaoficios.ar", "123456")
	if err != nil || robertoUser.Email != "roberto.gomez@guiaoficios.ar" {
		t.Fatalf("esperaba login exitoso para roberto.gomez@guiaoficios.ar con bcrypt, err=%v", err)
	}
	if _, err := store.AuthenticateUser(ctx, "roberto.gomez@guiaoficios.ar", "clave_incorrecta"); err == nil {
		t.Fatalf("no debía autenticar con clave incorrecta")
	}
}
