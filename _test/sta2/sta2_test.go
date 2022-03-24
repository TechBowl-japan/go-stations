package sta2_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func TestStation2(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Error("dbPathのセットに失敗しました。", err)
		return
	}

	t.Cleanup(func() {
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})

	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Error("DBの作成に失敗しました。", err)
		return
	}
	r := router.NewRouter(todoDB)
	srv := httptest.NewServer(r)
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	if err != nil {
		t.Error("リクエストの作成に失敗しました。", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("リクエストの送信に失敗しました。", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Error("レスポンスのクローズに失敗しました。", err)
			return
		}
	}()

	want := http.StatusNotFound
	if resp.StatusCode != want {
		t.Errorf("期待していない HTTP Status Code です, got = %d, want = %d", resp.StatusCode, want)
		return
	}
}
