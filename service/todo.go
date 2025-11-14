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
	
	// size==0 は空配列（テスト「Zero read」）を返す
	if size <= 0 {
		return []*model.TODO{}, nil
	}

	var (
		rows *sql.Rows
		err  error
	)
	if prevID > 0 {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
		rows, err = s.db.QueryContext(ctx, read, size)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	todos := make([]*model.TODO, 0)
	
	for rows.Next() {
		t := new(model.TODO)
		if err := rows.Scan(&t.ID, &t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
    const (
        update  = `UPDATE todos SET subject = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
        confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
    )

    if id == 0 {
        return nil, &model.ErrNotFound{Resource: "todo", ID: id}
    }

    // ★ ここでの空チェックや Trim はしない（subject="" は DB 制約に委ねる）
    res, err := s.db.ExecContext(ctx, update, subject, description, id)
    if err != nil {
        // ★ ラップしないでそのまま返す（sqlite3.Error 型が維持され、テスト期待と一致）
        return nil, err
    }

    n, err := res.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("failed to get affected rows: %w", err)
    }
    if n == 0 {
        return nil, &model.ErrNotFound{Resource: "todo", ID: id}
    }

    var t model.TODO
    if err := s.db.QueryRowContext(ctx, confirm, id).Scan(
        &t.ID, &t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt,
    ); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, &model.ErrNotFound{Resource: "todo", ID: id}
        }
        return nil, fmt.Errorf("failed to confirm updated todo: %w", err)
    }
    return &t, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	// ids が空なら何もしないで終了（要件どおり）
    if len(ids) == 0 {
        return nil
    }

    // プレースホルダ生成: "?", "?,?", "?,?,?" ... の部分
    placeholders := ""
    if len(ids) > 1 {
        placeholders = strings.Repeat(",?", len(ids)-1)
    }
    query := fmt.Sprintf(deleteFmt, placeholders)
    // 例: len(ids)=3 → "DELETE FROM todos WHERE id IN (?,?,?)"

    // []int64 → []interface{} に詰め替え（ExecContext に渡すため）
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        args[i] = id
    }

    // 実行
    res, err := s.db.ExecContext(ctx, query, args...)
    if err != nil {
        // sqlite3.Error の型を保ったまま返しておく
        return err
    }

    // 何件削除されたか確認
    n, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %w", err)
    }
    if n == 0 {
        // 一件も消せなかった → NotFound 扱い
        return &model.ErrNotFound{Resource: "todo"}
    }

	return nil
}
