package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	defaultFilePrefix = "x"
	suffixFirst       = 'a'
	suffixLast        = 'a'
)

type CLI struct {
	Stdout    io.Writer
	Stderr    io.Writer
	Stdin     io.Reader
	OutputDir string
}

func (cli *CLI) Run(args []string) error {
	commandOpt := &commandOption{}
	commandFS := flag.NewFlagSet("gosplit", flag.ExitOnError)
	commandFS.IntVar(&commandOpt.ByteCount, "b", 0, "byte count")
	commandFS.IntVar(&commandOpt.LineCount, "l", 0, "line count")

	if err := commandFS.Parse(args[1:]); err != nil {
		return err
	}
	if err := commandOpt.validate(); err != nil {
		return err
	}

	nonFlagArgs := commandFS.Args()
	commandArg := &commandArgument{}
	if len(nonFlagArgs) > 2 {
		return fmt.Errorf("too many arguments")
	}
	if len(nonFlagArgs) == 2 {
		commandArg.FilePath = nonFlagArgs[0]
		commandArg.FilePrefix = nonFlagArgs[1]
	}
	if err := commandArg.validate(); err != nil {
		return err
	}

	var filePrefix string
	if commandArg.FilePrefix != "" {
		filePrefix = commandArg.FilePrefix
	} else {
		filePrefix = defaultFilePrefix
	}

	if commandArg.FilePath != "" {
		f, err := os.Open(filepath.Join(cli.OutputDir, commandArg.FilePath))
		if err != nil {
			return err
		}
		defer f.Close()
		cli.Stdin = f
	}

	sc := bufio.NewScanner(cli.Stdin)
	if commandOpt.ByteCount > 0 {
		sc.Split(bufio.ScanBytes)
		var wg sync.WaitGroup
		buf := make([]byte, 0)
		for sc.Scan() {
			buf = append(buf, sc.Bytes()...)

			if len(buf) == commandOpt.ByteCount {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := os.WriteFile(filepath.Join(cli.OutputDir, fmt.Sprintf("%s%c%c", filePrefix, suffixFirst, suffixLast)), buf, 0644); err != nil {
						cli.Stderr.Write([]byte(err.Error()))
						return
					}
				}()
				wg.Wait()

				buf = make([]byte, 0)
				if err := updateSuffix(); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if commandOpt.LineCount > 0 {
		sc.Split(bufio.ScanLines)
		var wg sync.WaitGroup
		buf := make([]string, 0)
		for sc.Scan() {
			buf = append(buf, fmt.Sprintf("%s\n", sc.Text()))
			if len(buf) == commandOpt.LineCount {
				wg.Add(1)
				go func() {
					defer wg.Done()
					f, err := os.Create(filepath.Join(cli.OutputDir, fmt.Sprintf("%s%c%c", filePrefix, suffixFirst, suffixLast)))
					if err != nil {
						cli.Stderr.Write([]byte(err.Error()))
						return
					}
					defer f.Close()
					for _, line := range buf {
						f.WriteString(line)
					}
				}()
				wg.Wait()

				buf = make([]string, 0)
				if err := updateSuffix(); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return nil
}

func updateSuffix() error {
	suffixLast += 1
	if suffixLast > 'z' {
		suffixFirst += 1
		suffixLast = 'a'
	}
	if suffixFirst > 'z' {
		return fmt.Errorf("too many files")
	}
	return nil
}

type commandArgument struct {
	FilePath   string
	FilePrefix string
}

func (arg *commandArgument) validate() error {
	// ファイルが存在するかチェック
	if arg.FilePath != "" {
		_, err := os.Stat(arg.FilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

type commandOption struct {
	ByteCount int
	LineCount int
}

func (opt *commandOption) validate() error {
	if err := opt.hasOneOption(); err != nil {
		return err
	}
	return nil
}

// hasOneOption はオプションが一つだけ指定されていることをチェックする
func (opt *commandOption) hasOneOption() error {
	if opt.ByteCount != 0 && opt.LineCount != 0 {
		return fmt.Errorf("only one of -b, -l can be specified")
	}
	return nil
}
