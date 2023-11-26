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
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		log.Println("Failed to get todo with id")
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Failed to get id")
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, confirm, id)
	todo := &model.TODO{}
	err = row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		log.Println("Failed to scan todo")
		return nil, err
	}

	return todo, nil
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
		if size <= 0 {
			const query = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC`
			rows, err = s.db.QueryContext(ctx, query, size)
		} else {
			rows, err = s.db.QueryContext(ctx, read, size)
		}
	} else {
		if size <= 0 {
			const queryWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC`
			rows, err = s.db.QueryContext(ctx, queryWithID, prevID, size)
		} else {
			rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
		}
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	todos := []*model.TODO{}

	for rows.Next() {
		todo := &model.TODO{}

		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		todos = append(todos, todo)
	} 
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}


	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if rows == 0 {
		log.Println("Expected to update more than 1 row, but no rows were updated.")
		return nil, model.NewErrNotFound(id)
	}

	row := s.db.QueryRowContext(ctx, confirm, id)
	todo := &model.TODO{}
	err = row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) == 0 {
		return nil
	}
	placeholders := strings.Repeat(", ?", len(ids) - 1)
	query := fmt.Sprintf(deleteFmt, placeholders)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	id := make([]interface{}, len(ids))
	for i, ID := range ids {
		id[i] = ID
	}
	result, err := stmt.ExecContext(ctx, id...)
	if err != nil {
		log.Println(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == 0 {
		log.Println("No rows were deleted.")
    	return model.NewErrNotFound(ids[0])
	}
	return nil
}
