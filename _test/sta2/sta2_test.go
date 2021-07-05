package sta2_test

import (
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestStation2(t *testing.T) {
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

	resp, err := http.Get("http://localhost:8080")
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

	want := http.StatusNotFound
	if resp.StatusCode != want {
		t.Errorf("期待していない HTTP Status Code です, got = %d, want = %d", resp.StatusCode, want)
		return
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
