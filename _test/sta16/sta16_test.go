package sta16_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
)

func TestStation16(t *testing.T) {
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
		t.Error(" エラーが発生しました", err)
		return
	}

	t.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			t.Error(" エラーが発生しました", err)
			return
		}
	})

	for i := 0; i < 3; i++ {
		_, err = stmt.Exec(fmt.Sprintf("todo subject %d", i+1))
		if err != nil {
			t.Error(" エラーが発生しました", err)
			return
		}
	}

	stop, err := procStart(t)
	if err != nil {
		t.Error(" エラーが発生しました", err)
		return
	}

	t.Cleanup(func() {
		if err := stop(); err != nil {
			t.Error(" エラーが発生しました", err)
			return
		}
	})

	testcases := map[string]struct {
		PrevID   int
		Size     int
		TODOsLen int
	}{
		"Default read": {
			PrevID:   0,
			Size:     0,
			TODOsLen: 3,
		},
		"All read": {
			PrevID:   0,
			Size:     5,
			TODOsLen: 3,
		},
		"One read": {
			PrevID:   0,
			Size:     1,
			TODOsLen: 1,
		},
		"All read with prev id = 3": {
			PrevID:   3,
			Size:     5,
			TODOsLen: 2,
		},
		"All read with prev id = 1": {
			PrevID:   1,
			Size:     5,
			TODOsLen: 0,
		},
		"One read with prev id = 3": {
			PrevID:   3,
			Size:     1,
			TODOsLen: 1,
		},
		"One read with prev id = 1": {
			PrevID:   1,
			Size:     1,
			TODOsLen: 0,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var q string
			if tc.PrevID != 0 && tc.Size != 0 {
				q = fmt.Sprintf("?prev_id=%d&size=%d", tc.PrevID, tc.Size)
			} else if tc.PrevID != 0 {
				q = fmt.Sprintf("?prev_id=%d", tc.PrevID)
			} else if tc.Size != 0 {
				q = fmt.Sprintf("?size=%d", tc.Size)
			}

			resp, err := http.Get("http://localhost:8080/todos" + q)
			if err != nil {
				t.Error(" エラーが発生しました", err)
				return
			}
			defer resp.Body.Close()

			body := map[string]interface{}{
				"todos": []map[string]interface{}{},
			}
			err = json.NewDecoder(resp.Body).Decode(&body)
			if err != nil {
				t.Error(" エラーが発生しました", err)
				return
			}

			got, ok := body["todos"].([]interface{})
			if !ok {
				t.Error("todos field が見つかりません")
				return
			}

			if len(got) != tc.TODOsLen {
				t.Errorf("期待していない todos の長さです, got = %d, want = %d", len(got), tc.TODOsLen)
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
