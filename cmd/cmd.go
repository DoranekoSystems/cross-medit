package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"github.com/DoranekoSystems/cross-medit/pkg/converter"
	"github.com/DoranekoSystems/cross-medit/pkg/memory"
)

var AttachPid = 0

type Found struct {
	addrs     []int
	converter func(string) ([]byte, error)
	dataType  string
}

func Plist() error {
	processes, err := memory.EnumProcess()
	if err != nil {
		return err
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Pid < processes[j].Pid
	})

	for _, process := range processes {
		fmt.Printf("PID: %d, Process Name: %s\n", process.Pid, process.ProcessName)
	}

	return nil
}

func Attach(pid int) error {
	if AttachPid != 0 && pid == AttachPid {
		fmt.Println("Already attached.")
		return nil
	}
	AttachPid = pid
	return nil
}

func AttachByName(name string) error {
	processes, err := memory.EnumProcess()
	if err != nil {
		return err
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Pid < processes[j].Pid
	})

	for _, process := range processes {
		if name == process.ProcessName {
			if AttachPid != 0 && AttachPid == process.Pid {
				fmt.Println("Already attached.")
				return nil
			} else {
				fmt.Printf("Attach Success:%s %d\n", name, process.Pid)
				AttachPid = process.Pid
				return nil
			}
		}
	}
	fmt.Printf("Process not found.")
	return nil
}

func Find(targetVal string, dataType string) ([]Found, error) {
	founds := []Found{}
	addrRanges, err := memory.GetWritableAddrRanges(AttachPid)
	if err != nil {
		return nil, err
	}

	if dataType == "all" {
		// search string
		foundAddrs, err := memory.FindString(AttachPid, targetVal, addrRanges)
		if err == nil && len(foundAddrs) > 0 {
			founds = append(founds, Found{
				addrs:     foundAddrs,
				converter: converter.StringToBytes,
				dataType:  "UTF-8 string",
			})
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
		fmt.Println("------------------------")

		// search int
		foundAddrs, err = memory.FindWord(AttachPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.WordToBytes,
					dataType:  "word",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
		fmt.Println("------------------------")
		foundAddrs, err = memory.FindDword(AttachPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.DwordToBytes,
					dataType:  "dword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
		fmt.Println("------------------------")
		foundAddrs, err = memory.FindQword(AttachPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.QwordToBytes,
					dataType:  "qword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "string" {
		foundAddrs, _ := memory.FindString(AttachPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.StringToBytes,
					dataType:  "UTF-8 string",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "word" {
		foundAddrs, err := memory.FindWord(AttachPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.WordToBytes,
					dataType:  "word",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "dword" {
		foundAddrs, err := memory.FindDword(AttachPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.DwordToBytes,
					dataType:  "dword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "qword" {
		foundAddrs, err := memory.FindQword(AttachPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.QwordToBytes,
					dataType:  "qword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
	}

	return nil, errors.New("Error: specified datatype does not exist")
}

func Filter(targetVal string, prevFounds []Found) ([]Found, error) {
	founds := []Found{}
	writableAddrRanges, err := memory.GetWritableAddrRanges(AttachPid)
	if err != nil {
		return nil, err
	}
	addrRanges := [][2]int{}

	// check if previous result address exists in current memory map
	for i, prevFound := range prevFounds {
		targetBytes, _ := prevFound.converter(targetVal)
		targetLength := len(targetBytes)
		fmt.Printf("Check previous results of searching %s...\n", prevFound.dataType)
		fmt.Printf("Target Value: %s(%v)\n", targetVal, targetBytes)
		for _, prevAddr := range prevFound.addrs {
			for _, writable := range writableAddrRanges {
				if writable[0] < prevAddr && prevAddr < writable[1] {
					addrRanges = append(addrRanges, [2]int{prevAddr, prevAddr + targetLength})
				}
			}
		}
		foundAddrs, _ := memory.FindDataInAddrRanges(AttachPid, targetBytes, addrRanges)
		fmt.Printf("Found: %d!!\n", len(foundAddrs))
		if len(foundAddrs) < 10 {
			for _, v := range foundAddrs {
				fmt.Printf("Address: 0x%x\n", v)
			}
		}
		founds = append(founds, Found{
			addrs:     foundAddrs,
			converter: prevFound.converter,
			dataType:  prevFound.dataType,
		})
		if i != len(prevFounds)-1 {
			fmt.Println("------------------------")
		}
	}
	return founds, nil
}

func Patch(targetVal string, targetAddrs []Found) error {
	for _, found := range targetAddrs {
		targetBytes, _ := found.converter(targetVal)
		for _, targetAddr := range found.addrs {
			err := memory.WriteMemory(AttachPid, targetAddr, targetBytes)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Successfully patched!")
	return nil
}

func Detach() error {
	if AttachPid == 0 {
		fmt.Println("Already detached.")
		return nil
	}
	AttachPid = 0
	return nil
}

func Dump(beginAddress int, endAddress int) error {
	memSize := endAddress - beginAddress
	buf := make([]byte, memSize)
	memory := memory.ReadMemory(AttachPid, buf, beginAddress, endAddress)
	fmt.Printf("Address range: 0x%x - 0x%x\n", beginAddress, endAddress)
	fmt.Println("--------------------------------------------")
	fmt.Printf("%s", hex.Dump(memory))
	return nil
}
