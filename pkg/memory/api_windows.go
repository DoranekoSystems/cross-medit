// +build windows

package memory

import (
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

type ProcessInfo struct {
	Pid         int
	ProcessName string
}

func ReadMemory(pid int, buffer []byte, beginAddr int, size int) []byte {
	processHandle, _ := windows.OpenProcess(windows.PROCESS_VM_READ, false, uint32(pid))
	defer windows.CloseHandle(processHandle)

	var bytesRead uintptr
	windows.ReadProcessMemory(processHandle, uintptr(beginAddr), &buffer[0], uintptr(size), &bytesRead)

	return buffer
}

func WriteMemory(pid int, targetAddr int, targetVal []byte) error {
	processHandle, _ := windows.OpenProcess(windows.PROCESS_VM_WRITE|windows.PROCESS_VM_OPERATION, false, uint32(pid))
	defer windows.CloseHandle(processHandle)

	var bytesWritten uintptr
	windows.WriteProcessMemory(processHandle, uintptr(targetAddr), &targetVal[0], uintptr(len(targetVal)), &bytesWritten)

	return nil
}

func EnumRegion(pid int) (string, error) {
	processHandle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(processHandle)

	var buffer string
	var memInfo windows.MemoryBasicInformation
	var addr uintptr

	for {
		err:=windows.VirtualQueryEx(processHandle, addr, &memInfo, unsafe.Sizeof(memInfo))
		if err != nil{
			break
		}

		permissions := "----"
		switch memInfo.Protect {
		case windows.PAGE_EXECUTE_READ:
			permissions = "r-x"
		case windows.PAGE_EXECUTE_READWRITE:
			permissions = "rwx"
		case windows.PAGE_EXECUTE_WRITECOPY:
			permissions = "rwxc"
		case windows.PAGE_READONLY:
			permissions = "r--"
		case windows.PAGE_READWRITE:
			permissions = "rw-"
		}

		buffer += fmt.Sprintf("%x-%x %s\n", uintptr(addr), uintptr(addr)+memInfo.RegionSize, permissions)
		addr += uintptr(memInfo.RegionSize)
	}
	return buffer, nil
}

func EnumProcess() ([]ProcessInfo, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var processEntry windows.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))
	if err := windows.Process32First(snapshot, &processEntry); err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	for {
		processes = append(processes, ProcessInfo{
			Pid:         int(processEntry.ProcessID),
			ProcessName: windows.UTF16ToString(processEntry.ExeFile[:]),
		})

		if err := windows.Process32Next(snapshot, &processEntry); err != nil {
			break
		}
	}

	return processes, nil
}
