package terminal

import (
	//"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

const (
	SYS_UP = 1000 + iota
	SYS_DOWN
	SYS_LEFT
	SYS_RIGHT
	SYS_PARSE
)

const (
	syscall_IGNBRK = 0x1
	syscall_BRKINT = 0x2
	syscall_PARMRK = 0x8
	syscall_ISTRIP = 0x20
	syscall_INLCR  = 0x40
	syscall_IGNCR  = 0x80
	syscall_ICRNL  = 0x100
	syscall_IXON   = 0x200
	syscall_OPOST  = 0x1
	syscall_ECHO   = 0x8
	syscall_ECHONL = 0x10
	syscall_ICANON = 0x100
	syscall_ISIG   = 0x80
	syscall_IEXTEN = 0x400
	syscall_CSIZE  = 0x300
	syscall_PARENB = 0x1000
	syscall_CS8    = 0x300
	syscall_VMIN   = 0x10
	syscall_VTIME  = 0x11

	syscall_TCGETS = 0x40487413
	syscall_TCSETS = 0x80487414
)

func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	val = int(r)
	if e != 0 {
		panic(e)
	}
	return
}

func tcsetattr(fd int, termios *syscall.Termios) {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall_TCSETS), uintptr(unsafe.Pointer(termios)))
	if r != 0 {
		panic(os.NewSyscallError("SYS_IOCTL", e))
	}
}

func tcgetattr(fd int, termios *syscall.Termios) {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall_TCGETS), uintptr(unsafe.Pointer(termios)))
	if r != 0 {
		panic(os.NewSyscallError("SYS_IOCTL", e))
	}
}

func getch() (int, string) {
	var (
		in        int
		err       error
		sigio     = make(chan os.Signal)
		orig_tios syscall.Termios
	)
	in, err = syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	signal.Notify(sigio, syscall.SIGIO)

	fcntl(in, syscall.F_SETFL, syscall.O_ASYNC|syscall.O_NONBLOCK)
	tcgetattr(in, &orig_tios)
	tios := orig_tios
	tios.Iflag &^= syscall_BRKINT | syscall_IXON
	tios.Lflag &^= syscall_ECHO | syscall_ICANON | syscall_ISIG | syscall_IEXTEN
	tios.Cflag &^= syscall_CSIZE | syscall_PARENB
	tios.Cflag |= syscall_CS8
	tios.Cc[syscall_VMIN] = 1
	tios.Cc[syscall_VTIME] = 0

	tcsetattr(in, &tios)
	defer func() {
		tcsetattr(in, &orig_tios)
		syscall.Close(in)
	}()
	buf := make([]byte, 128)

LOOP:
	<-sigio
	n, err := syscall.Read(in, buf)
	if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
		goto LOOP
	}
	//fmt.Println(n, buf[0:n])
	if n == 1 {
		return int(buf[0]), ""
	} else if n == 3 {
		switch buf[2] {
		case 65:
			return SYS_UP, ""
		case 66:
			return SYS_DOWN, ""
		case 67:
			return SYS_RIGHT, ""
		case 68:
			return SYS_LEFT, ""
		default:
			return SYS_PARSE, string(buf[0:n])
		}
	} else {
		return SYS_PARSE, string(buf[0:n])
	}
	return 0, ""
}
