package util

import (
	"syscall"

	"golang.org/x/sys/windows"
)

type LazyOrd struct {
	Ord  uintptr
	l    *windows.LazyDLL
	addr *uintptr
}

func NewProcByOrdinal(d *windows.LazyDLL, ord uintptr) *LazyOrd {
	return &LazyOrd{l: d, Ord: ord}
}

func (p *LazyOrd) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	if p.addr == nil {
		proc, err := windows.GetProcAddressByOrdinal(windows.Handle(p.l.Handle()), uintptr(p.Ord))
		if err != nil {
			return 0, 0, err
		}
		p.addr = &proc
	}
	return syscall.SyscallN(*p.addr, a...)
}
