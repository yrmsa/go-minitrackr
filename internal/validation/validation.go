package validation

import "strings"

const MaxTitleLength = 500

func ValidateTitle(title string) (string, bool) {
	title = strings.TrimSpace(title)
	if title == "" || len(title) > MaxTitleLength {
		return "", false
	}
	return title, true
}

func ValidateStatus(status string) bool {
	switch status {
	case "todo", "doing", "done":
		return true
	}
	return false
}

func ValidatePriority(priority string) bool {
	switch priority {
	case "low", "medium", "high":
		return true
	}
	return false
}
