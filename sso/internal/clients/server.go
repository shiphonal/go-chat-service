package clients

import (
	client "ChatService/sso/internal/clients/service"
	"ChatService/sso/internal/config"
	"ChatService/sso/internal/lib/jwt"
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type ClientFabric struct {
	SSO        *client.ClientSSO
	HttpServer *http.Server
}

//go:embed all:front/*
var frontFS embed.FS

func ClientMustLoad(cnf *config.Config, logger *slog.Logger) *ClientFabric {
	ssoClient, err := client.New(
		context.Background(),
		logger,
		cnf.Clients.SSO.Addr,
		cnf.Clients.SSO.Timeout,
		cnf.Clients.SSO.RetriesCount,
	)
	if err != nil {
		logger.Error("failed to initialize SSO client", err)
		os.Exit(1)
	}
	logger.Info("ClientSSO initialized")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: setupRoutes(ssoClient),
	}

	return &ClientFabric{
		HttpServer: httpServer,
		SSO:        ssoClient,
	}
}

func setupRoutes(cli *client.ClientSSO) *http.ServeMux {
	mux := http.NewServeMux()

	// Загрузка HTML шаблонов
	templates := template.Must(template.ParseFS(frontFS,
		"front/templates/register.html",
		"front/templates/login.html",
		"front/templates/profile.html",
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

	// Страница регистрации
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")
			username := r.FormValue("username")

			_, err := cli.Register(r.Context(), username, email, password)
			if err != nil {
				http.Error(w, "Register failed", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			w.Header().Set("Content-Type", "text/html")
			if err := templates.ExecuteTemplate(w, "register.html", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	// Страница входа
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")
			appID := int64(1)

			token, err := cli.Login(r.Context(), email, password, appID)
			if err != nil {
				http.Error(w, "Login failed", http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: token,
			})
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		} else {
			w.Header().Set("Content-Type", "text/html")
			if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	// Страница профиля
	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		TokenInfo := jwt.ValidateToken(tokenCookie.Value, os.Getenv("SECRET"))
		if TokenInfo.Error != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			newName := r.FormValue("new_name")
			_, err := cli.ChangeName(r.Context(), TokenInfo.UserID, newName)
			if err != nil {
				http.Error(w, "Failed to change name", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		} else {
			w.Header().Set("Content-Type", "text/html")
			if err := templates.ExecuteTemplate(w, "profile.html", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	return mux
}
