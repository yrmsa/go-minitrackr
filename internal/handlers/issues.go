package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/yrmsa/go-minitrackr/internal/db"
)

type Handler struct {
	db        *db.DB
	templates *template.Template
}

func New(database *db.DB, templates *template.Template) *Handler {
	return &Handler{
		db:        database,
		templates: templates,
	}
}

func (h *Handler) ListIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := h.db.ListIssues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issues)
}

func (h *Handler) GetIssue(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/api/issues/")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	issue, err := h.db.GetIssue(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if issue == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

// CreateIssue is deprecated - use CreateBoardIssue or CreateBacklogIssue
func (h *Handler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Use /board/issues or /backlog/issues", http.StatusGone)
}

// UpdateIssue is deprecated - use UpdateBoardIssue or UpdateBacklogIssue
func (h *Handler) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Use /board/issues or /backlog/issues", http.StatusGone)
}

// DeleteIssue is deprecated - use DeleteBoardIssue or DeleteBacklogIssue
func (h *Handler) DeleteIssue(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Use /board/issues or /backlog/issues", http.StatusGone)
}

func parseID(path, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.TrimSuffix(idStr, "/")
	return strconv.ParseInt(idStr, 10, 64)
}

func formatID(id int64) string {
	return strconv.FormatInt(id, 10)
}
