package main

import (
	"context"
	"io"
)

func openStreamWithIP(config *Config, parentCtx context.Context) (*Stream, error) {
	// todo: not implement
	panic("`openStreamWithIP` not implement")
	return nil, nil
}

func openStream(parentCtx context.Context) (*Stream, error) {
	// todo: not implement
	panic("`openStream` not implement")
	return nil, nil
}

func closeStream(closer io.ReadWriteCloser) {
	panic("`closeStream` not implement")
	// todo: not implement
}
