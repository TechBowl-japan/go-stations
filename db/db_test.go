package db_test

import (
	"errors"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/mattn/go-sqlite3"
)

func TestNewDB(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path string
		err  error
	}{
		"Normal":                {path: "../.sqlite3/db_test.db", err: nil},
		"Directory not created": {path: "../.nothing/db_test.db", err: sqlite3.ErrCantOpenFullPath},
	}

	t.Cleanup(func() {
		for _, c := range cases {
			if c.err == nil {
				if err := os.Remove(c.path); err != nil {
					t.Error("failed to cleanup testdata, err =", err)
				}
			}
		}
	})

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := db.NewDB(c.path)
			if err != nil {
				if errors.Is(c.err, err) {
					t.Errorf("unexpected value, given = %s, expected = %s\n", err, c.err)
				}
			}
		})
	}
}
