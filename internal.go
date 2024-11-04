package pun

import (
	"context"
	"os"
	"sync"
)

func opened(_ctx context.Context) bool {
	select {
	case <-_ctx.Done():
		return false
	default:
		return true
	}
}

func newStream(file *os.File, parentCtx context.Context) *Stream {
	stream := new(Stream)
	stream._inputStream = make(chan []byte)
	stream._outputStream = make(chan []byte)
	stream._io = file
	stream.InputStream = stream._inputStream
	stream.OutputStream = stream._outputStream

	ctx, cancelFunc := context.WithCancel(parentCtx)
	stream._life = &ctx
	stream._cancelFunc = &cancelFunc
	return stream
}

func (stream *Stream) open() (<-chan []byte, chan<- []byte) {
	if !stream._reading {
		stream._reading = true
		go stream.readFromTunnel()
		go stream.writeToTunnel()
	}

	return stream.OutputStream, stream.InputStream
}

// todo: move Pool to Device
var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 64*1024)
	},
}

func (stream *Stream) readFromTunnel() {
	for opened(*stream._life) {
		packet := bufferPool.Get().([]byte)
		num, err := stream._io.Read(packet)
		if err != nil {
			continue
		}
		if !stream._reading {
			break
		}
		stream._outputStream <- packet[:num]
	}
}

func (stream *Stream) writeToTunnel() {
	for opened(*stream._life) {
		packet := <-stream._inputStream
		if !stream._reading {
			break
		}
		_, _ = stream._io.Write(packet)
		bufferPool.Put(packet)
	}
}
