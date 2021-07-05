package sta9_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStation9(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Error("エラーが発生しました", err)
		return
	}

	t.Cleanup(func() {
		if err := os.Remove(dbPath); err != nil {
			t.Error("エラーが発生しました", err)
			return
		}
	})

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
		Subject            string
		Description        string
		WantHTTPStatusCode int
	}{
		"Subject is empty": {
			WantHTTPStatusCode: http.StatusBadRequest,
		},
		"Description is empty": {
			Subject:            "todo subject",
			WantHTTPStatusCode: http.StatusOK,
		},
		"Subject and Description is not empty": {
			Subject:            "todo subject",
			Description:        "todo description",
			WantHTTPStatusCode: http.StatusOK,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			resp, err := http.Post("http://localhost:8080/todos", "application/json",
				bytes.NewBufferString(fmt.Sprintf(`{"subject":"%s","description":"%s"}`, tc.Subject, tc.Description)))
			if err != nil {
				t.Error("エラーが発生しました", err)
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Error("エラーが発生しました", err)
					return
				}
			}()

			if resp.StatusCode != tc.WantHTTPStatusCode {
				t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, tc.WantHTTPStatusCode)
				return
			}

			if tc.WantHTTPStatusCode != http.StatusOK {
				return
			}

			var m map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
				t.Error("エラーが発生しました", err)
				return
			}

			v, ok := m["todo"]
			if !ok {
				t.Error("todo field not found")
				return
			}

			got, ok := v.(map[string]interface{})
			if !ok {
				t.Error("todo field not object")
				return
			}
			want := map[string]interface{}{
				"subject":     tc.Subject,
				"description": tc.Description,
			}

			now := time.Now().UTC()
			diff := cmp.Diff(got, want, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
				switch k {
				case "id":
					if vv, _ := v.(float64); vv == 0 {
						t.Errorf("id を数値に変換できません, got = %s", k)
					}
					return true
				case "created_at", "updated_at":
					vv, ok := v.(string)
					if !ok {
						t.Errorf("日付が文字列に変換できません, got = %+v", k)
						return true
					}
					if tt, err := time.Parse(time.RFC3339, vv); err != nil {
						t.Errorf("日付が期待しているフォーマットではありません, got = %s", k)
					} else if now.Before(tt) {
						t.Errorf("日付が未来の日付になっています, got = %s", tt)
					}
					return true
				}

				return false
			}))

			if diff != "" {
				t.Error("期待していない値です\n", diff)
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
