package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_Run(t *testing.T) {
	t.Parallel()
	const (
		noErr  = false
		hasErr = true
	)

	cases := map[string]struct {
		args    string
		in      string
		want    []string
		wantErr bool
		Err     error
	}{
		"normal/byte count":       {"./gosplit -b 1", "aaa", []string{"a", "a", "a"}, noErr, nil},
		"normal/line count":       {"./gosplit -l 1", "aaa", []string{"aaa\n"}, noErr, nil},
		"abnormal/with prefix":    {"./gosplit -l 1 noexist prefix", "", []string{}, noErr, errors.New("stat noexist: no such file or directory")},
		"abnormal/too many flags": {"./gosplit -l 1 -b 2", "aaa", []string{}, hasErr, errors.New("only one of -b, -l can be specified")},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got bytes.Buffer
			cli := &CLI{
				Stdout:    &got,
				Stderr:    &got,
				Stdin:     strings.NewReader(tt.in),
				OutputDir: t.TempDir(),
			}

			args := strings.Split(tt.args, " ")
			err := cli.Run(args)
			switch {
			case tt.wantErr && err == nil:
				t.Error("expected error did not occur")
			case !tt.wantErr && err != nil && err.Error() != tt.Err.Error():
				t.Error("unexpected error:", err)
			}

			files, err := os.ReadDir(cli.OutputDir)
			if err != nil {
				t.Error(err)
			}
			var buf bytes.Buffer
			for i, file := range files {
				if file.IsDir() {
					continue
				}
				f, err := os.Open(filepath.Join(cli.OutputDir, file.Name()))
				if err != nil {
					t.Error(err)
				}
				defer f.Close()
				io.Copy(&buf, f)
				if buf.String() != tt.want[i] {
					t.Errorf("got %q, want %q", buf.String(), tt.want[i])
				}
				buf.Reset()
			}
		})
	}
}
