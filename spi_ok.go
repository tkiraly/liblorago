package liblorago

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

//#define DEBUG_MSG(str)                fprintf(stderr, str)
//#define DEBUG_PRINTF(fmt, args...)    fprintf(stderr,"%s:%d: "fmt, __FUNCTION__, __LINE__, args)
const (
	READ_ACCESS               = 0x00
	WRITE_ACCESS              = 0x80
	SPI_SPEED                 = 8000000
	LGW_BURST_CHUNK           = 1024
	LGW_SPI_MUX_MODE0         = 0x0 /* No FPGA */
	LGW_SPI_MUX_MODE1         = 0x1 /* FPGA, with spi mux header */
	LGW_SPI_MUX_TARGET_SX1301 = 0x0
	LGW_SPI_MUX_TARGET_FPGA   = 0x1
	LGW_SPI_MUX_TARGET_EEPROM = 0x2
	LGW_SPI_MUX_TARGET_SX127X = 0x3

	spiIOCWrMode        = 0x40016B01
	spiIOCWrBitsPerWord = 0x40016B03
	spiIOCWrMaxSpeedHz  = 0x40046B04
	spiIOCRdMode        = 0x80016B01
	spiIOCRdBitsPerWord = 0x80016B03
	spiIOCRdMaxSpeedHz  = 0x80046B04
	spiIOCMessage0      = 0x40006B00
	spiIOCIncrementor   = 0x200000
)

var lock sync.Mutex = sync.Mutex{}
var locks map[*os.File]*sync.Mutex = make(map[*os.File]*sync.Mutex, 0)

type spiIOCTransfer struct {
	txBuf       uint64
	rxBuf       uint64
	length      uint32
	speedHz     uint32
	delayus     uint16
	bitsPerWord byte
	csChange    byte
	pad         uint32
}

func spiIOCMessageN(n uint32) uint32 {
	return (spiIOCMessage0 + (n * spiIOCIncrementor))
}

func Lgw_spi_open(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	mode := 0
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), spiIOCWrMode, uintptr(unsafe.Pointer(&mode)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return nil, err
	}
	speed := SPI_SPEED
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), spiIOCWrMaxSpeedHz, uintptr(unsafe.Pointer(&speed)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return nil, err
	}
	bpw := 8
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), spiIOCWrBitsPerWord, uintptr(unsafe.Pointer(&bpw)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return nil, err
	}
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), spiIOCRdMode, uintptr(unsafe.Pointer(&mode)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return nil, err
	}
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), spiIOCRdMaxSpeedHz, uintptr(unsafe.Pointer(&speed)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return nil, err
	}
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), spiIOCRdBitsPerWord, uintptr(unsafe.Pointer(&bpw)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return nil, err
	}
	locks[file] = &sync.Mutex{}

	log.Print("Note: SPI port opened and configured ok\n")
	return file, nil
}

func Lgw_spi_close(f *os.File) error {
	lock.Lock()
	defer lock.Unlock()
	locks[f].Lock()
	defer locks[f].Unlock()
	err := f.Close()
	if err != nil {
		return err
	}
	return nil
}

func Lgw_spi_w(file *os.File, spi_mux_mode, spi_mux_target, address, data byte) error {
	write := make([]byte, 0)

	if (address & 0x80) != 0 {
		fmt.Print("WARNING: SPI address > 127\n")
	}

	if spi_mux_mode == LGW_SPI_MUX_MODE1 {
		write = append(write, []byte{0, 0, 0}...)
		write[0] = spi_mux_target
		write[1] = WRITE_ACCESS | (address & 0x7F)
		write[2] = data
	} else {
		write = append(write, []byte{0, 0}...)
		write[0] = WRITE_ACCESS | (address & 0x7F)
		write[1] = data
	}

	k := spiIOCTransfer{}
	k.length = uint32(len(write))
	k.txBuf = uint64(uintptr(unsafe.Pointer(&write[0])))
	k.csChange = 0
	//k.speedHz = SPI_SPEED
	//k.bitsPerWord = 8

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(spiIOCMessageN(1)), uintptr(unsafe.Pointer(&k)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return err
	}
	return nil
}

func Lgw_spi_r(file *os.File, spi_mux_mode, spi_mux_target, address byte) (byte, error) {
	write := make([]byte, 0)

	if (address & 0x80) != 0 {
		fmt.Print("WARNING: SPI address > 127\n")
	}

	if spi_mux_mode == LGW_SPI_MUX_MODE1 {
		write = append(write, []byte{0, 0, 0}...)
		write[0] = spi_mux_target
		write[1] = READ_ACCESS | (address & 0x7F)
	} else {
		write = append(write, []byte{0, 0}...)
		write[0] = READ_ACCESS | (address & 0x7F)
	}

	read := make([]byte, len(write))
	k := spiIOCTransfer{}
	k.length = uint32(len(write))
	k.txBuf = uint64(uintptr(unsafe.Pointer(&write[0])))
	k.rxBuf = uint64(uintptr(unsafe.Pointer(&read[0])))
	k.csChange = 0
	k.speedHz = SPI_SPEED
	k.bitsPerWord = 8

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(spiIOCMessageN(1)), uintptr(unsafe.Pointer(&k)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return 0, err
	}
	return read[1], nil
}

func Lgw_spi_wb(file *os.File, spi_mux_mode, spi_mux_target, address byte, data []byte) error {
	write := make([]byte, 0)
	if (address & 0x80) != 0 {
		fmt.Print("WARNING: SPI address > 127\n")
	}
	//if spi_mux_mode == LGW_SPI_MUX_MODE1 {
	//	write = append(write, []byte{0, 0}...)
	//	write[0] = spi_mux_target
	//	write[1] = WRITE_ACCESS | (address & 0x7F)
	//} else {
	write = append(write, []byte{0}...)
	write[0] = WRITE_ACCESS | (address & 0x7F)
	//}

	size_to_do := uint32(len(data))
	k := make([]spiIOCTransfer, 2)

	k[0].txBuf = uint64(uintptr(unsafe.Pointer(&write[0])))
	k[0].length = uint32(len(write))
	byte_transfered := uint64(0)
	for i := 0; size_to_do > 0; i++ {
		chunk_size := uint32(LGW_BURST_CHUNK)
		if size_to_do < LGW_BURST_CHUNK {
			chunk_size = size_to_do
		}
		offset := uint32(i) * LGW_BURST_CHUNK
		k[1].txBuf = uint64(uintptr(unsafe.Pointer(&data[0+offset])))
		k[1].length = chunk_size
		I, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(spiIOCMessageN(2)), uintptr(unsafe.Pointer(&k[0])))
		if errno != 0 {
			err := syscall.Errno(errno)
			return err
		}
		byte_transfered += uint64(I) - uint64(k[0].length)
		//fmt.Printf("BURST WRITE: to trans %d # chunk %d # transferred %d \n", size_to_do, chunk_size, byte_transfered)
		size_to_do -= chunk_size /* subtract the quantity of data already transferred */
	}
	return nil
}

func Lgw_spi_rb(file *os.File, spi_mux_mode, spi_mux_target, address byte, size uint16) ([]byte, error) {
	read := make([]byte, size)
	write := make([]byte, 0)

	if (address & 0x80) != 0 {
		fmt.Print("WARNING: SPI address > 127\n")
	}

	if spi_mux_mode == LGW_SPI_MUX_MODE1 {
		write = append(write, []byte{0, 0}...)
		write[0] = spi_mux_target
		write[1] = READ_ACCESS | (address & 0x7F)
	} else {
		write = append(write, []byte{0}...)
		write[0] = READ_ACCESS | (address & 0x7F)
	}

	size_to_do := uint32(size)
	k := make([]spiIOCTransfer, 2)

	k[0].txBuf = uint64(uintptr(unsafe.Pointer(&write[0])))
	k[0].length = uint32(len(write))

	byte_transfered := uint64(0)
	for i := 0; size_to_do > 0; i++ {
		chunk_size := uint32(LGW_BURST_CHUNK)
		if size_to_do < LGW_BURST_CHUNK {
			chunk_size = size_to_do
		}
		offset := uint32(i * LGW_BURST_CHUNK)
		k[1].rxBuf = uint64(uintptr(unsafe.Pointer(&read[0+offset])))
		k[1].length = chunk_size
		I, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(spiIOCMessageN(2)), uintptr(unsafe.Pointer(&k[0])))
		if errno != 0 {
			err := syscall.Errno(errno)
			return nil, err
		}
		byte_transfered += uint64(I) - uint64(k[0].length)
		//fmt.Printf("BURST READ: to trans %d # chunk %d # transferred %d \n", size_to_do, chunk_size, byte_transfered)
		size_to_do -= chunk_size
	}
	return read, nil
}
