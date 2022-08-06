package service

import (
	"context"
	"database/sql"

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

	var err error
	var res sql.Result
	var id int64

	res, err = s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	data := s.db.QueryRowContext(ctx, confirm, int(id))
	var todo model.TODO
	err = data.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, err
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	rows, err := s.db.QueryContext(ctx, read, size)
	if err != nil {
		return nil, err
	}
	var data []*model.TODO
	for rows.Next() {
		t := &model.TODO{}
		if err := rows.Scan(&t.ID, &t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		data = append(data, t)
	}
	return data, err
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	var err error

	_, err = s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	data := s.db.QueryRowContext(ctx, confirm, int(id))
	var todo *model.TODO
	err = data.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return todo, err
}

// DeleteTODO deletes TODOs on DB by ids.
// func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
// 	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

// 	return nil
// }

func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?)`
	for id := range ids {
		_, err := s.db.ExecContext(ctx, deleteFmt, id)
		if err != nil {
			return err
		}
	}
	return nil
}
