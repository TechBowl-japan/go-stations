package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the newly created TODO to return
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	todo.ID = id
	log.Println("Scan todo", &todo)

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
	log.Println("prevID", prevID)
	log.Println("size", size)
	// If prevID is 0, fetch the latest 'size' records; otherwise, fetch records before prevID
	if size <= 0 {
		if prevID == 0 {
			return []*model.TODO{}, nil
		}
		return []*model.TODO{}, err
	} else if prevID == 1 {
		return []*model.TODO{}, err
	} else if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
		if err != nil {
			return []*model.TODO{}, err
		}
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			return []*model.TODO{}, err
		}
	}

	defer rows.Close()

	var todos []*model.TODO
	for rows.Next() {
		var todo model.TODO
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return []*model.TODO{}, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	// 3. 更新された行が0件の場合、TODOが存在しないと判断し、ErrNotFoundを返す
	if rowsAffected == 0 {
		return nil, &model.ErrNotFound{} // または適切なカスタムエラー
	}

	// Fetch the newly created TODO to return
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	todo.ID = id

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	result, err := s.db.ExecContext(ctx, fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1)), args...)
	if err != nil {
		log.Println("DeleteTODO err")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("DeleteTODO RowsAffected err")
		return err
	}
	if rowsAffected == 0 {
		return &model.ErrNotFound{} // または適切なカスタムエラー
	}
	return nil
}
