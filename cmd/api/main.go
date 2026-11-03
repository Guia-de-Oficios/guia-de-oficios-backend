package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"

	"guiaoficios/backend/internal/handlers"
	"guiaoficios/backend/internal/repository"
	"guiaoficios/backend/internal/telemetry"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://guia_user:guia_password@localhost:5432/guia_oficios?sslmode=disable"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error al inicializar el driver PostgreSQL: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error fatal: No se pudo conectar a PostgreSQL/PostGIS en DATABASE_URL (%v). Asegurate de levantar 'docker compose up -d'.", err)
	}

	if err := runMigrations(ctx, db); err != nil {
		log.Printf("Aviso al ejecutar migraciones automáticas: %v", err)
	}

	store := repository.NewPostgresStore(db)
	storageMode := "postgresql+postgis"

	// Inicializar recolector de métricas clave para Datadog (DogStatsD UDP)
	_ = telemetry.GetMetrics()

	mux := http.NewServeMux()
	apiHandler := handlers.NewAPIHandler(store, storageMode, uploadDir)
	apiHandler.RegisterRoutes(mux)

	log.Printf("=======================================================================")
	log.Printf("🚀 API REST 'Guía de Oficios' escuchando en http://localhost:%s", port)
	log.Printf("📦 Motor de Base de Datos:   %s (100%% persistencia real, sin mocks)", storageMode)
	log.Printf("🔍 Trazabilidad Activa:      Header X-Trace-ID / X-Datadog-Trace-ID en cada request")
	log.Printf("📊 Métricas Clave Datadog:   DogStatsD UDP + http://localhost:%s/api/v1/metrics", port)
	log.Printf("📂 Directorio de Uploads:    %s (/uploads/)", uploadDir)
	log.Printf("🔗 Health Check:             http://localhost:%s/api/v1/health", port)
	log.Printf("🔗 API Rubros:               http://localhost:%s/api/v1/categories", port)
	log.Printf("🔗 API Profesionales:        http://localhost:%s/api/v1/professionals", port)
	log.Printf("=======================================================================")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      telemetry.ObservabilityMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error en el servidor HTTP: %v", err)
	}
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	migrationFiles := []string{
		"migrations/001_init.sql",
		"migrations/002_seed.sql",
	}
	for _, relPath := range migrationFiles {
		path := filepath.Clean(relPath)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Clean(filepath.Join("..", relPath))
		}
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if _, err := db.ExecContext(ctx, string(sqlBytes)); err != nil {
			return err
		}
	}
	return nil
}
