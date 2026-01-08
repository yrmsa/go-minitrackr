package db

import (
	"database/sql"
	"time"

	"github.com/yrmsa/go-minitrackr/internal/models"
)

func (db *DB) CreateIssue(title, status, priority string) (*models.Issue, error) {
	now := time.Now().Unix()
	result, err := db.Exec(
		"INSERT INTO issues (title, status, priority, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		title, status, priority, now, now,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Issue{
		ID:        id,
		Title:     title,
		Status:    status,
		Priority:  priority,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (db *DB) GetIssue(id int64) (*models.Issue, error) {
	issue := &models.Issue{}
	err := db.QueryRow(
		"SELECT id, title, status, priority, created_at, updated_at FROM issues WHERE id = ?",
		id,
	).Scan(&issue.ID, &issue.Title, &issue.Status, &issue.Priority, &issue.CreatedAt, &issue.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (db *DB) ListIssues() ([]*models.Issue, error) {
	rows, err := db.Query(
		"SELECT id, title, status, priority, created_at, updated_at FROM issues ORDER BY created_at DESC LIMIT 1000",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := make([]*models.Issue, 0, 100)
	for rows.Next() {
		issue := &models.Issue{}
		if err := rows.Scan(&issue.ID, &issue.Title, &issue.Status, &issue.Priority, &issue.CreatedAt, &issue.UpdatedAt); err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}

	return issues, rows.Err()
}

func (db *DB) UpdateIssue(id int64, title, status, priority string) error {
	now := time.Now().Unix()
	_, err := db.Exec(
		"UPDATE issues SET title = ?, status = ?, priority = ?, updated_at = ? WHERE id = ?",
		title, status, priority, now, id,
	)
	return err
}

func (db *DB) DeleteIssue(id int64) error {
	_, err := db.Exec("DELETE FROM issues WHERE id = ?", id)
	return err
}
