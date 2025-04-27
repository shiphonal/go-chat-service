package app

import (
	client "ChatService/crud/internal/clients/service"
	"ChatService/crud/internal/config"
	"ChatService/crud/internal/lib/jwt"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
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
		Handler: setupRoutes(crudClient),
	}

	return &ClientFabric{
		HttpServer: httpServer,
		CRUD:       crudClient,
	}
}

func setupRoutes(cli *client.ClientCRUD) *http.ServeMux {
	mux := http.NewServeMux()

	/*templates := template.Must(template.ParseGlob("client/front/template/*.html"))*/
	fs := http.FileServer(http.Dir("client/front/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Работа с сообщениями
	mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		TokenInfo := jwt.ValidateToken(tokenCookie.Value, os.Getenv("SECRET"))

		switch r.Method {
		case "POST":
			// Отправка сообщения
			messageType := r.FormValue("test")
			content := r.FormValue("content")
			userID := TokenInfo.UserID

			mid, err := cli.SentMessage(r.Context(), messageType, content, userID)
			if err != nil {
				http.Error(w, "Failed to send message", http.StatusInternalServerError)
				return
			}
			w.Write([]byte(fmt.Sprintf("Message sent with ID: %d", mid)))

		case "GET":
			// Получение сообщения
			midStr := r.URL.Query().Get("mid")
			mid, err := strconv.ParseInt(midStr, 10, 64) // В реальности преобразовать midStr в int64
			if err != nil {
				http.Error(w, "Failed to parse mid", http.StatusBadRequest)
			}

			message, err := cli.GetMessage(r.Context(), tokenCookie.Value, mid)
			if err != nil {
				http.Error(w, "Failed to get message", http.StatusInternalServerError)
				return
			}
			w.Write([]byte(message))
		}
	})

	return mux
}
