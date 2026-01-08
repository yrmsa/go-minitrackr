package handlers

import (
	"net/http"

	"github.com/yrmsa/go-minitrackr/internal/validation"
)

func (h *Handler) CreateBoardIssue(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title, valid := validation.ValidateTitle(r.FormValue("title"))
	if !valid {
		http.Error(w, "Invalid title", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	if status == "" {
		status = "todo"
	} else if !validation.ValidateStatus(status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	priority := r.FormValue("priority")
	if priority == "" {
		priority = "medium"
	} else if !validation.ValidatePriority(priority) {
		http.Error(w, "Invalid priority", http.StatusBadRequest)
		return
	}

	issue, err := h.db.CreateIssue(title, status, priority)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.templates.ExecuteTemplate(w, "board-card", issue)
}

func (h *Handler) UpdateBoardIssue(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/board/issues/")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	current, err := h.db.GetIssue(id)
	if err != nil || current == nil {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	if title != "" {
		var valid bool
		title, valid = validation.ValidateTitle(title)
		if !valid {
			http.Error(w, "Invalid title", http.StatusBadRequest)
			return
		}
	} else {
		title = current.Title
	}

	status := r.FormValue("status")
	if status == "" {
		status = current.Status
	} else if !validation.ValidateStatus(status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	priority := r.FormValue("priority")
	if priority == "" {
		priority = current.Priority
	} else if !validation.ValidatePriority(priority) {
		http.Error(w, "Invalid priority", http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateIssue(id, title, status, priority); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updated, _ := h.db.GetIssue(id)
	w.Header().Set("Content-Type", "text/html")

	if status != current.Status {
		// Status changed - use OOB swap to move to new column
		w.Write([]byte(`<div id="card-` + formatID(id) + `" hx-swap-oob="beforeend:#column-` + status + ` .board-issues">`))
		h.templates.ExecuteTemplate(w, "board-card", updated)
		w.Write([]byte(`</div>`))
	} else {
		// Same status - update in place
		h.templates.ExecuteTemplate(w, "board-card", updated)
	}
}

func (h *Handler) DeleteBoardIssue(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/board/issues/")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteIssue(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
