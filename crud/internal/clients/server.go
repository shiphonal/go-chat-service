package clients

import (
	client "ChatService/crud/internal/clients/service"
	"ChatService/crud/internal/config"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type ClientFabric struct {
	CRUD       *client.ClientCRUD
	HttpServer *http.Server
}

func ClientMustLoad(cnf *config.Config, logger *slog.Logger) *ClientFabric {
	crudClient, err := client.New(
		context.Background(),
		logger,
		cnf.Clients.CRUD.Addr,
		cnf.Clients.CRUD.Timeout,
		cnf.Clients.CRUD.RetriesCount,
	)
	if err != nil {
		logger.Error("failed to initialize CRUD client", err)
		os.Exit(1)
	}
	logger.Info("ClientCRUD initialized")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: setupRoutes(crudClient, logger),
	}

	return &ClientFabric{
		HttpServer: httpServer,
		CRUD:       crudClient,
	}
}

//go:embed all:front/*
var frontFS embed.FS

func setupRoutes(cli *client.ClientCRUD, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	templates := template.Must(template.ParseFS(frontFS,
		"front/templates/index.html",
	))

	// Обработчик статических файлов (CSS, JS)
	staticFS, _ := fs.Sub(frontFS, "front/static")
	fileServer := http.FileServer(http.FS(staticFS))

	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем правильные MIME-типы
		switch {
		case strings.HasSuffix(r.URL.Path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(r.URL.Path, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(r.URL.Path, ".png"):
			w.Header().Set("Content-Type", "image/png")
		}
		fileServer.ServeHTTP(w, r)
	})))

	// Serve index.html for root path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// API endpoints for messages
	mux.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			messageType := r.FormValue("type")
			content := r.FormValue("content")
			userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64) // In real app, get from session/token

			mid, err := cli.SentMessage(r.Context(), messageType, content, userID)
			if err != nil {
				logger.Error("failed to send message", "error", err)
				http.Error(w, "Failed to send message", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status": "success", "message_id": %d}`, mid)

		case "GET":
			midStr := r.URL.Query().Get("id")
			mid, err := strconv.ParseInt(midStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid message ID", http.StatusBadRequest)
				return
			}

			// In real app, validate user has access to this message
			message, err := cli.GetMessage(r.Context(), "", mid) // Empty token for demo
			if err != nil {
				logger.Error("failed to get message", "error", err)
				http.Error(w, "Failed to get message", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status": "success", "message": %q}`, message)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
