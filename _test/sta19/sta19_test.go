package sta19_test

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
)

func TestStation19(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Error(" エラーが発生しました", err)
		return
	}

	d, err := db.NewDB(dbPath)
	if err != nil {
		t.Error(" エラーが発生しました", err)
		return
	}

	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Error(" エラーが発生しました", err)
			return
		}
	})
	t.Cleanup(func() {
		if err := os.Remove(dbPath); err != nil {
			t.Error(" エラーが発生しました", err)
			return
		}
	})

	stmt, err := d.Prepare(`INSERT INTO todos(subject) VALUES(?)`)
	if err != nil {
		t.Error(err)
		return
	}

	t.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			t.Error("エラーが発生しました", err)
			return
		}
	})

	for i := 0; i < 3; i++ {
		_, err = stmt.Exec("sbuject")
		if err != nil {
			t.Error("エラーが発生しました", err)
			return
		}
	}

	stop, err := procStart(t)
	if err != nil {
		t.Error("エラーが発生しました", err)
		return
	}

	t.Cleanup(func() {
		if err := stop(); err != nil {
			t.Error("エラーが発生しました", err)
			return
		}
	})

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
		tc := tc
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todos",
				bytes.NewBufferString(fmt.Sprintf(`{"ids":[%s]}`, strings.Join(tc.IDs, ","))))
			if err != nil {
				t.Error("エラーが発生しました", err)
				return
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error("エラーが発生しました", err)
				return
			}
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Error("エラーが発生しました", err)
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

func procStart(t *testing.T) (func() error, error) {
	t.Helper()

	cmd := exec.Command("go", "run", "../../main.go")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Second)

	stop := func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	return stop, nil
}
