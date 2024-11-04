package pun

import (
	"context"
	"fmt"
	"github.com/vishvananda/netlink"
	"io"
	"os"
	"syscall"
	"unsafe"
)

const _IFF_MULTI_QUEUE = 0x0100

var _req *_ifReq
var _index = 0

type _ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func _initIfReq(name string) *_ifReq {
	req := new(_ifReq)
	copy(req.Name[:], name)
	req.Flags = syscall.IFF_TUN | syscall.IFF_NO_PI | _IFF_MULTI_QUEUE
	return req
}

func _openTunDevices() (int, error) {
	fd, err := syscall.Open("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return -1, fmt.Errorf("open /dev/net/tun %v", err)
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(_req)))
	if errno != 0 {
		return -1, fmt.Errorf("ioctl tunsetiff %v", errno)
	}

	return fd, nil
}

func openStreamWithIP(config *Config, parentCtx context.Context) (*Stream, error) {
	if _req != nil {
		return nil, fmt.Errorf("you have already opened a tun device")
	}
	_req = _initIfReq(config.Name)

	fd, err := _openTunDevices()
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
	file := os.NewFile(uintptr(fd), fmt.Sprintf("pun%d", _index))
	_index = _index + 1
	return newStream(file, parentCtx), nil
}

func openStream(parentCtx context.Context) (*Stream, error) {
	fd, err := _openTunDevices()
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), fmt.Sprintf("pun%d", _index))
	_index = _index + 1
	return newStream(file, parentCtx), nil
}

func closeStream(closer io.ReadWriteCloser) {
	_ = closer.Close()
}
