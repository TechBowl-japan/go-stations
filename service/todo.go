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
	//TODOを挿入
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}
	//resultの新しく作成されたTODOのIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	//IDを使用してTODOを取得
	row := s.db.QueryRowContext(ctx, confirm, id)
	todo := &model.TODO{
		ID: id,
	}
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	//成功した場合、新しいTODOを返す
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
	log.Printf("Received prevID=%d, size=%d", prevID, size)

	if prevID > 0 {
		log.Printf("Executing query with prevID=%d and size=%d", prevID, size)
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
		log.Printf("Executing query: %s with size=%d", read, size)
		rows, err = s.db.QueryContext(ctx, read, size)
	}

	if err != nil {
		//クエリ実行中にエラーが発生した場合
		log.Printf("Query execution failed: %v", err)
		return nil, err
	}
	defer rows.Close() //rowsを必ず閉じる

	//TODOリストを格納するスライス
	todos := []*model.TODO{}
	for rows.Next() {
		todo := &model.TODO{}
		//結果セットの行をスキャン
		if err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Printf("Row scanning failed: %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}

	//繰り返し処理後のエラーを確認
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration failed: %v", err)
		return nil, err
	}

	//結果を返す
	log.Printf("Retrieved todos count: %d, details: %v", len(todos), todos)
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	//TODOを更新
	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		//更新処理中にエラーが発生すれば、そのエラーを返す
		return nil, err
	}

	//影響を受けた行数を確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		//行数取得中にエラーが発生すれば、そのエラーを返す
		return nil, err
	}
	//もし更新された行が0のとき
	if rowsAffected == 0 {
		//エラーとして、「対象のTODOが見つかりませんでした」と返す。
		return nil, &model.ErrNotFound{Resource: "TODO"}
	}
	//更新されたTODOを取得
	row := s.db.QueryRowContext(ctx, confirm, id)
	todo := &model.TODO{
		ID: id,
	}
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		//データ取得中にエラーが発生すれば、そのエラーを返す
		return nil, err
	}
	//更新されたTODOを返す
	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (%s)`
	if len(ids) == 0 {
		return nil
	}

	Formats := strings.Repeat("?,", len(ids))
	Formats = Formats[:len(Formats)-1]
	query := fmt.Sprintf(deleteFmt, Formats)

	inter := make([]interface{}, len(ids))
	for i, id := range ids {
		inter[i] = id
	}

	result, err := s.db.ExecContext(ctx, query, inter...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &model.ErrNotFound{Resource: "TODO"}
	}

	return nil
}
