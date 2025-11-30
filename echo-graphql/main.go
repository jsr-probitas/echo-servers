package main

import (
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"

	"github.com/jsr-probitas/dockerfiles/echo-graphql/graph"
)

func main() {
	cfg := LoadConfig()

	resolver := graph.NewResolver()
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	// HTTP transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// WebSocket transport for subscriptions
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	// Enable introspection
	srv.Use(extension.Introspection{})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// GraphQL playground
	http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))

	// GraphQL endpoint
	http.Handle("/graphql", srv)

	log.Printf("Starting server on %s", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), nil); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
