package handlers

import (
	"net/http"

	"github.com/yrmsa/go-minitrackr/internal/models"
)

type ViewData struct {
	Title        string
	ActiveView   string
	Issues       interface{}
	TodoIssues   interface{}
	DoingIssues  interface{}
	DoneIssues   interface{}
}

func (h *Handler) BacklogView(w http.ResponseWriter, r *http.Request) {
	issues, err := h.db.ListIssues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := ViewData{
		Title:      "Backlog",
		ActiveView: "backlog",
		Issues:     issues,
	}
	if err := h.templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) BoardView(w http.ResponseWriter, r *http.Request) {
	issues, err := h.db.ListIssues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Group issues by status
	var todoIssues, doingIssues, doneIssues []*models.Issue
	for _, issue := range issues {
		switch issue.Status {
		case "todo":
			todoIssues = append(todoIssues, issue)
		case "doing":
			doingIssues = append(doingIssues, issue)
		case "done":
			doneIssues = append(doneIssues, issue)
		}
	}

	data := ViewData{
		Title:        "Board",
		ActiveView:   "board",
		TodoIssues:   todoIssues,
		DoingIssues:  doingIssues,
		DoneIssues:   doneIssues,
	}
	if err := h.templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) InboxView(w http.ResponseWriter, r *http.Request) {
	data := ViewData{
		Title:      "Inbox",
		ActiveView: "inbox",
	}
	if err := h.templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) SettingsView(w http.ResponseWriter, r *http.Request) {
	data := ViewData{
		Title:      "Settings",
		ActiveView: "settings",
	}
	if err := h.templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
