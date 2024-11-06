package pun

import (
	"context"
	"io"
	"sync"
)

func newStream(file io.ReadWriteCloser, parentCtx context.Context) *Stream {
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
		return make([]byte, 1600)
	},
}

func (stream *Stream) readFromTunnel() {
	for {
		select {
		case <-(*stream._life).Done():
			return
		default:
			packet := bufferPool.Get().([]byte)
			num, err := stream._io.Read(packet)
			if !stream._reading {
				break
			}
			if err != nil {
				continue
			}

			stream._outputStream <- packet[:num]
		}
	}
}

func (stream *Stream) writeToTunnel() {
	for {
		select {
		case <-(*stream._life).Done():
			return
		case packet := <-stream._inputStream:
			if !stream._reading {
				break
			}
			if num, _ := stream._io.Write(packet); num == 1600 {
				bufferPool.Put(packet)
			}
		}
	}
}
