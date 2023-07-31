// +build darwin

package memory

// #include "../../native/native_darwin.h"
import "C"
import (
	"fmt"
	"unsafe"
)

const bufferSize = 65536 * 4

type ProcessInfo struct {
	Pid         int
	ProcessName string
}

func ReadMemory(pid int, buffer []byte, beginAddr int, endAddr int) []byte {
	size := endAddr - beginAddr
	C.read_memory_native(C.int(pid), C.mach_vm_address_t(beginAddr), C.mach_vm_size_t(size), (*C.uchar)(&buffer[0]))
	return buffer
}

func WriteMemory(pid int, targetAddr int, targetVal []byte) error {
	size := len(targetVal)
	result := C.write_memory_native(C.int(pid), C.mach_vm_address_t(targetAddr), C.mach_vm_size_t(size), (*C.uchar)(&targetVal[0]))
	if result == -1 {
		return fmt.Errorf("Failed to write memory")
	}
	return nil
}

func EnumRegion(pid int) (string, error) {
	const bufferSize = 65536

	var buffer [bufferSize]C.char

	C.enumerate_regions_to_buffer(C.pid_t(pid), (*C.char)(&buffer[0]), C.size_t(bufferSize))

	goString := C.GoString(&buffer[0])

	return goString, nil
}

func EnumProcess() ([]ProcessInfo, error) {
	var count C.size_t
	procInfos := C.enumprocess_native(&count)
	if procInfos == nil {
		return nil, fmt.Errorf("Failed to enumerate processes")
	}
	defer C.free(unsafe.Pointer(procInfos))

	processes := make([]ProcessInfo, count)
	for i := 0; i < int(count); i++ {
		cProcInfo := (*C.ProcessInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(procInfos)) + uintptr(i)*unsafe.Sizeof(C.ProcessInfo{})))
		processes[i] = ProcessInfo{
			Pid:         int(cProcInfo.pid),
			ProcessName: C.GoString(cProcInfo.processname),
		}
		C.free(unsafe.Pointer(cProcInfo.processname)) // C関数内でstrdupが使われているため
	}

	return processes, nil
}
