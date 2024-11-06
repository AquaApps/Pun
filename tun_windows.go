package pun

import (
	"context"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/tun"
	"io"
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

	netInterface, err := netlink.LinkByName(config.Name)
	if err != nil {
		return nil, err
	}

	addrV4 := &netlink.Addr{IPNet: &(config.CIDRv4), Label: ""}

	if err = netlink.LinkSetMTU(netInterface, config.MTU); err != nil {
		return nil, err
	}

	if err = netlink.AddrAdd(netInterface, addrV4); err != nil {
		return nil, err
	}

	if err = netlink.LinkSetUp(netInterface); err != nil {
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
