package sta1_test

import (
	"bytes"
	"os/exec"
	"regexp"
	"testing"
)

func TestStation1(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("go", "version")
	w := bytes.NewBuffer(nil)
	cmd.Stdout = w

	if err := cmd.Run(); err != nil {
		t.Error("エラーが発生しました", err)
		return
	}

	if !regexp.MustCompile(`go\sversion\sgo1\.\d+(\.\d+)?\s.+/.+`).Match(w.Bytes()) {
		t.Error("go version の実行結果に問題があります")
		return
	}
}
