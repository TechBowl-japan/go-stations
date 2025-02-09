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
		confirm = `SELECT id,subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	//PrepareContextを使ってSQLステートメントの実行準備を行う
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	//インサートしてあげなければいけない。ExecContext メソッドを使う
	result, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	//リーザルからidを取得する。
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	//取得した情報をQueryRowContext メソッドを使って、返してあげる
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
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

	//使用するクエリの判定
	//SQLの行数変数
	var rows *sql.Rows
	//エラー変数
	var err error
	if prevID == 0 {
		//prevIDがないとき
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		//prevIDがあるとき
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//結果をスライスにまとめて、返す
	todos := []*model.TODO{}
	for rows.Next() {
		var todo model.TODO
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	//SQLステートメントの実行準備
	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	//SQL実行
	result, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	//行数判定
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, &model.ErrNotFound{}
	}

	//保存するTODOを読み取り
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.ErrNotFound{}
		}
		return nil, err
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	//fmt.Sprintf関数と strings.Repeat関数を組み合わせてクエリに?の処理を追加
	add_ids := strings.Repeat("?,", len(ids))
	add_ids = add_ids[:len(add_ids)-1]
	const deleteFmt = `DELETE FROM todos WHERE id IN (%s)`
	query := fmt.Sprintf(deleteFmt, add_ids)

	//[]interface{}にidsを詰め直して、繰り返し処理で引数で展開してargsに受け渡す
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	//クエリを実行準備
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//クエリを実行
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	//リーザルから削除されたTODOを取得
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//削除されたTODOが0であった場合ErrNotFoundを返す
	if rows == 0 {
		return &model.ErrNotFound{}
	}

	//idsがからのスライスの場合、なんも処理しないでnilを返す
	if len(ids) == 0 {
		return nil
	}

	return nil
}
