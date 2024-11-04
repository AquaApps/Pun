package main

import (
	"context"
	"io"
	"net"
	"time"
)

type Config struct {
	Name   string
	CIDRv4 net.IPNet
	CIDRv6 net.IPNet // todo: not available now...
	MTU    int
}

// Device provide OpenStream(), OpenExtraStream(), Close()
type Device struct {
	_config     Config
	_life       *context.Context
	_cancelFunc *context.CancelFunc
	_stream     *Stream // default stream
}

func New(config *Config, parentCtx context.Context) (*Device, error) {
	device := new(Device)
	deviceCtx, cancelFunc := context.WithCancel(parentCtx)
	device._life = &deviceCtx
	device._cancelFunc = &cancelFunc

	stream, err := openStreamWithIP(config, deviceCtx)
	if err != nil {
		return nil, err
	}
	device._stream = stream
	return device, nil
}

func (device *Device) OpenStream() (<-chan []byte, chan<- []byte) {
	return device._stream.open()
}

// OpenExtraStream make a parallel stream
// From version 3.8, Linux supports multiqueue tuntap which can uses multiple file descriptors (queues) to parallelize packets sending or receiving.
// The device allocation is the same as before, and if user wants to create multiple queues,
// TUNSETIFF with the same device name must be called many times with IFF_MULTI_QUEUE flag.
// --https://www.kernel.org/doc/html/latest/networking/tuntap.html
func (device *Device) OpenExtraStream() (*Stream, error) {
	stream, err := openStream(*device._life)
	if err != nil {
		return nil, err
	}
	stream.open()
	return stream, nil
}

func (device *Device) Close() {
	(*device._cancelFunc)()
	device._stream.Close()
}

//////////////////////////////////////////////////////////////////////////////////////////////////

// Stream provide InputStream, OutputStream, and Close()
type Stream struct {
	InputStream  chan<- []byte // exposed stream
	OutputStream <-chan []byte

	_life         *context.Context
	_cancelFunc   *context.CancelFunc
	_inputStream  chan []byte // internal stream
	_outputStream chan []byte
	_io           io.ReadWriteCloser // real io file
	_reading      bool               // whether goroutine is running
}

func (stream *Stream) Close() {
	(*stream._cancelFunc)()
	stream._reading = false
	<-time.NewTimer(1 * time.Second).C
	closeStream(stream._io)
	close(stream._inputStream)
	close(stream._outputStream)
}
