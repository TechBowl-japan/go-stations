package sta19_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/osamingo/go-todo-app/db"
	"github.com/osamingo/go-todo-app/handler/router"
)

func TestStation19(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Errorf("dbPathのセットに失敗しました。%v", err)
		return
	}

	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Errorf("データベースの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := todoDB.Close(); err != nil {
			t.Errorf("データベースのクローズに失敗しました: %v", err)
			return
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})

	stmt, err := todoDB.Prepare(`INSERT INTO todos(subject) VALUES(?)`)
	if err != nil {
		t.Errorf("ステートメントの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			t.Errorf("ステートメントのクローズに失敗しました: %v", err)
			return
		}
	})

	for i := 0; i < 3; i++ {
		if _, err := stmt.Exec("sbuject"); err != nil {
			t.Errorf("todoの追加に失敗しました: %v", err)
			return
		}
	}

	r := router.NewRouter(todoDB)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testcases := map[string]struct {
		IDs                []string
		WantHTTPStatusCode int
	}{
		"Empty Ids": {
			IDs:                []string{},
			WantHTTPStatusCode: http.StatusBadRequest,
		},
		"Not found ID": {
			IDs:                []string{"4"},
			WantHTTPStatusCode: http.StatusNotFound,
		},
		"One delete": {
			IDs:                []string{"1"},
			WantHTTPStatusCode: http.StatusOK,
		},
		"Multiple delete": {
			IDs:                []string{"2", "3"},
			WantHTTPStatusCode: http.StatusOK,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, srv.URL+"/todos",
				bytes.NewBufferString(fmt.Sprintf(`{"ids":[%s]}`, strings.Join(tc.IDs, ","))))
			if err != nil {
				t.Errorf("リクエストの作成に失敗しました: %v", err)
				return
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("リクエストの送信に失敗しました: %v", err)
				return
			}
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("レスポンスのクローズに失敗しました: %v", err)
					return
				}
			})

			if resp.StatusCode != tc.WantHTTPStatusCode {
				t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, tc.WantHTTPStatusCode)
				return
			}
		})
	}
}
