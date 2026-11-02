package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"guiaoficios/backend/internal/handlers"
	"guiaoficios/backend/internal/repository"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	var store repository.Store
	storageMode := "in-memory (precargado con datos Stitch)"

	if dsn != "" {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if pingErr := db.PingContext(ctx); pingErr == nil {
				store = repository.NewPostgresStore(db)
				storageMode = "postgresql+postgis"
			} else {
				log.Printf("Aviso: No se pudo conectar a PostgreSQL (%v). Usando modo en memoria.", pingErr)
			}
		}
	}

	if store == nil {
		store = repository.NewInMemoryStore()
	}

	mux := http.NewServeMux()
	apiHandler := handlers.NewAPIHandler(store, storageMode)
	apiHandler.RegisterRoutes(mux)

	// Servir la interfaz Stitch intacta desde ../ui-stitch o ./ui-stitch
	uiDir := filepath.Clean("../ui-stitch")
	if _, err := os.Stat(uiDir); os.IsNotExist(err) {
		uiDir = filepath.Clean("ui-stitch")
	}
	fs := http.FileServer(http.Dir(uiDir))
	mux.Handle("/ui/", http.StripPrefix("/ui/", fs))

	log.Printf("=======================================================================")
	log.Printf("🚀 Servidor 'Guía de Oficios' escuchando en http://localhost:%s", port)
	log.Printf("📦 Modo de almacenamiento: %s", storageMode)
	log.Printf("🎨 Interfaz Stitch Original: http://localhost:%s/ui/app.html", port)
	log.Printf("🔗 API Rubros:               http://localhost:%s/api/v1/categories", port)
	log.Printf("🔗 API Profesionales:        http://localhost:%s/api/v1/professionals", port)
	log.Printf("🔗 API Panel Profesional:    http://localhost:%s/api/v1/professionals/aaaa1111-1111-1111-1111-111111111111/dashboard", port)
	log.Printf("=======================================================================")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error en el servidor HTTP: %v", err)
	}
}
