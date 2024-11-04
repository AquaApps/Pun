package pun

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestLinuxPlatform(t *testing.T) {
	appCtx := context.Background()
	config := Config{
		Name: "TestPun",
		MTU:  1500,
	}
	if IPv4, CIDRv4, err := net.ParseCIDR("10.10.10.14/24"); err != nil {
		log.Fatal(err)
	} else {
		config.CIDRv4 = *CIDRv4
		config.CIDRv4.IP = IPv4
	}

	device, err := New(&config, appCtx)
	if err != nil {
		log.Fatal(err)
	}
	out, in := device.OpenStream()
	go func() {
		for {
			data := <-out
			log.Println("origin", data)
		}
	}()
	go func() {
		for {
			packet := make([]byte, 4*1024)
			length := fillRandomBytes(packet)
			in <- packet[:length]
		}

	}()

	stream, err := device.OpenExtraStream()
	if err != nil {
		log.Fatal(err)
	}
	data := <-stream.OutputStream
	log.Println("extra", data)
	stream.Close()
	log.Println("waiting for signal")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func fillRandomBytes(packet []byte) int {
	for i := range packet {
		packet[i] = byte(rand.Intn(256))
	}
	return len(packet)
}
