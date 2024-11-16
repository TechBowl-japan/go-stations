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

	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if description == "" {
		description = ""
	}

	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	resTODO := &model.TODO{}
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(
		&resTODO.Subject,
		&resTODO.Description,
		&resTODO.CreatedAt,
		&resTODO.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	resTODO.ID = id

	return resTODO, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	if size == 0 {
		return []*model.TODO{}, nil
	}

	if prevID == 1 {
		return []*model.TODO{}, nil
	}

	var (
		rows *sql.Rows
		err  error
	)

	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []*model.TODO

	for rows.Next() {
		todo := &model.TODO{}
		if err := rows.Scan(
			&todo.ID,
			&todo.Subject,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowAffected == 0 {
		return nil, &model.ErrNotFound{}
	}

	resTODO := &model.TODO{}
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(
		&resTODO.Subject,
		&resTODO.Description,
		&resTODO.CreatedAt,
		&resTODO.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	resTODO.ID = id

	return resTODO, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) == 0 {
		return nil
	}

	placeholders := strings.Repeat(", ?", len(ids)-1)
	query := fmt.Sprintf(deleteFmt, placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
