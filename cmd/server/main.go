package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/yrmsa/go-minitrackr/internal/config"
	"github.com/yrmsa/go-minitrackr/internal/db"
	"github.com/yrmsa/go-minitrackr/internal/handlers"
	"github.com/yrmsa/go-minitrackr/internal/middleware"
	"github.com/yrmsa/go-minitrackr/internal/static"
	"github.com/yrmsa/go-minitrackr/internal/templates"
)

func main() {
	cfg := config.Load()

	// Set memory limit
	debug.SetMemoryLimit(parseMemLimit(cfg.MemLimit))
	debug.SetGCPercent(50)

	// Initialize database
	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Load templates
	tmpl, err := templates.Load()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Initialize handlers
	h := handlers.New(database, tmpl)

	// Setup routes
	mux := http.NewServeMux()
	
	// UI routes
	mux.HandleFunc("/", h.BacklogView)
	mux.HandleFunc("/backlog", h.BacklogView)
	mux.HandleFunc("/board", h.BoardView)
	mux.HandleFunc("/inbox", h.InboxView)
	mux.HandleFunc("/settings", h.SettingsView)
	
	// API routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/issues", issuesRouter(h))
	mux.HandleFunc("/api/issues/", issueRouter(h))
	
	// Board HTMX routes
	mux.HandleFunc("/board/issues", boardIssuesRouter(h))
	mux.HandleFunc("/board/issues/", boardIssueRouter(h))
	
	// Backlog HTMX routes
	mux.HandleFunc("/backlog/issues", backlogIssuesRouter(h))
	mux.HandleFunc("/backlog/issues/", backlogIssueRouter(h))
	
	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS()))))
	mux.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(static.FS(), "swagger.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Starting go-minitrackr on %s (mem limit: %s)", addr, cfg.MemLimit)
	printMemStats()
	
	handler := middleware.Recovery(middleware.Logging(mux))
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

func issuesRouter(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ListIssues(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func issueRouter(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetIssue(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func boardIssuesRouter(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateBoardIssue(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func boardIssueRouter(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			h.UpdateBoardIssue(w, r)
		case http.MethodDelete:
			h.DeleteBoardIssue(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func backlogIssuesRouter(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateBacklogIssue(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func backlogIssueRouter(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			h.UpdateBacklogIssue(w, r)
		case http.MethodDelete:
			h.DeleteBacklogIssue(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func parseMemLimit(limit string) int64 {
	var val int64
	fmt.Sscanf(limit, "%dMiB", &val)
	return val * 1024 * 1024
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc=%vMB Sys=%vMB NumGC=%v", m.Alloc/1024/1024, m.Sys/1024/1024, m.NumGC)
}
