package util

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/pkg/errors"
)

func Command(name string, arg ...string) (resp string, err error) {
	return command(exec.Command(name, arg...))
}

func CommandBash(arg string) (string, error) {
	return command(exec.Command("bash", "-c", arg))
}

func CommandContext(ctx context.Context, name string, arg ...string) (resp string, err error) {
	return command(exec.CommandContext(ctx, name, arg...))
}

func CommandBashContext(ctx context.Context, arg string) (string, error) {
	return command(exec.CommandContext(ctx, "bash", "-c", arg))
}

func command(cmd *exec.Cmd) (resp string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = errors.Wrap(errors.New(stderr.String()), err.Error())
		return
	}
	resp = out.String()
	return
}
