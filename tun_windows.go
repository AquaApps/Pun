package pun

import (
	"context"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
	"io"
	"net/netip"
)

type winDev struct {
	dev tun.Device
}

func (w *winDev) Close() error {
	return w.dev.Close()
}

func (w *winDev) Write(b []byte) (int, error) {
	return w.dev.Write(b, 0)
}

func (w *winDev) Read(b []byte) (int, error) {
	return w.dev.Read(b, 0)
}

func openStreamWithIP(config *Config, parentCtx context.Context) (*Stream, error) {
	id := &windows.GUID{
		Data2: 0xFFFF,
		Data3: 0xFFFF,
		Data4: [8]byte{0xFF, 0xe9, 0x76, 0xe5, 0x8c, 0x74, 0x06, 0x3e},
	}
	dev, err := tun.CreateTUNWithRequestedGUID(config.Name, id, config.MTU)
	if err != nil {
		return nil, err
	}
	nativeTunDevice := dev.(*tun.NativeTun)
	link := winipcfg.LUID(nativeTunDevice.LUID())

	ipPrefix, err := netip.ParsePrefix(config.CIDRv4.String())
	if err != nil {
		return nil, err
	}

	err = link.AddIPAddress(ipPrefix)
	if err != nil {
		return nil, err
	}

	return newStream(&winDev{dev: dev}, parentCtx), nil
}

func openStream(parentCtx context.Context) (*Stream, error) {
	// todo: not implement
	panic("`openStream` not implement")
	return nil, nil
}

func closeStream(closer io.ReadWriteCloser) {
	// todo: not implement
	panic("`closeStream` not implement")
}
