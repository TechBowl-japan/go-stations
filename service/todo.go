package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// Insert the TODO
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Retrieve the TODO from the database
	row := s.db.QueryRowContext(ctx, confirm, id)
	var todo model.TODO
	err = row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		updateBoth = `UPDATE todos SET subject = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		updateOnly = `UPDATE todos SET subject = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		confirm    = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// idバリデーション
	if id == 0 {
		return nil, &model.ErrNotFound{Resource: "todo", ID: int(id)}
	}

	// subjectバリデーション
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	// 更新クエリの選択
	var (
		query string
		args  []interface{}
	)
	if strings.TrimSpace(description) == "" {
		query = updateOnly
		args = []interface{}{subject, id}
	} else {
		query = updateBoth
		args = []interface{}{subject, description, id}
	}

	// 更新実行
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	// 更新件数チェック
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return nil, &model.ErrNotFound{Resource: "todo", ID: int(id)}
	}

	// 更新後データを再取得
	var t model.TODO
	if err := s.db.QueryRowContext(ctx, confirm, id).Scan(
    	&t.ID, &t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
    	if errors.Is(err, sql.ErrNoRows) {
        	return nil, &model.ErrNotFound{Resource: "todo", ID: int(id)}
    }
    return nil, fmt.Errorf("failed to confirm updated todo: %w", err)
}

	return &t, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
