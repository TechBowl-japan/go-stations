package service

import (
	"context"
	"database/sql"
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
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	todo := model.TODO{
		ID: lastID,
	}

	if err := s.db.
		QueryRowContext(ctx, confirm, lastID).
		Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt, &todo.DeletedAt); err != nil {
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

	var rows *sql.Rows
	var err error
	if prevID == 0 {
		if size == 0 {
			size = 3
		}
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}

	if err == sql.ErrNoRows {
		return nil, &model.ErrNotFound{
			Message: fmt.Sprintf("received ID %d does not exist in the database", prevID),
		}
	}

	defer rows.Close()

	var todos []*model.TODO
	for rows.Next() {
		todo := model.TODO{}
		rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		todos = append(todos, &todo)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, &model.ErrNotFound{
			Message: fmt.Sprintf("received ID %d does not exist in the database", id),
		}
	}
	todo := model.TODO{
		ID: id,
	}
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	const deleteQueryTemplate = `DELETE FROM todos WHERE id IN (?%s)`
	query := fmt.Sprintf(deleteQueryTemplate, strings.Repeat(", ?", len(ids)-1))
	contextArgs := make([]interface{}, len(ids))

	for i, id := range ids {
		contextArgs[i] = id
	}

	res, err := s.db.ExecContext(ctx, query, contextArgs...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &model.ErrNotFound{
			Message: fmt.Sprintf("received ID %d does not exist in the database", ids),
		}
	}

	return nil
}
