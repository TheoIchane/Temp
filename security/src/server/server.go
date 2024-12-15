package server

// http server for forum with graceful shutdown

import (
	"context"
	"errors"
	"fmt"
	"forum/src/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Server() {

	handlers.InitTemplates()
	fs := http.FileServer(http.Dir("src/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/profil", handlers.ProfileHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/registeroauth", handlers.RegisterOauth)
	http.HandleFunc("/topics/", handlers.TopicHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/oauth/", handlers.OauthHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler) // Logout handler when you click a button that is routed through this handler
	http.HandleFunc("/admin", handlers.AdminHandler)
	http.HandleFunc("/uploads/posts", handlers.UploadImageHandler)
	http.HandleFunc("/uploads/avatars", handlers.UploadAvatarHandler)
	http.HandleFunc("/comments", handlers.CommentHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/post/create", handlers.CreatePostHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/dislike", handlers.DislikeHandler)
	server := &http.Server{
		Addr: ":8080",
	}

	// Goroutine for graceful shutdown
	fmt.Printf("Server started on: http://localhost%s\n", server.Addr)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println(("Stopped serving new connections"))
	}()

	SignalChan := make(chan os.Signal, 1)
	signal.Notify(SignalChan, syscall.SIGINT, syscall.SIGTERM)
	<-SignalChan

	shutdownContext, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete")

}
