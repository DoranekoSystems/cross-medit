package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/DoranekoSystems/cross-medit/cmd"

	prompt "github.com/c-bata/go-prompt"
)

var appPID string
var addrCache []cmd.Found
var specificPID string

func executor(in string) {
	command := strings.Split(in, " ")
	if in == "ps" {
		if err := cmd.Plist(); err != nil {
			log.Fatal(err)
		}

	} else if command[0] == "attach" {
		slice := strings.Split(in, " ")
		var pid int
		if len(slice) > 1 {
			pid, _ = strconv.Atoi(slice[1])
		} else if appPID != "" {
			pid, _ = strconv.Atoi(appPID)
		} else {
			fmt.Println("PID cannot be specified.")
		}
		cmd.Attach(pid)

	} else if command[0] == "attachbyname" {
		slice := strings.Split(in, " ")
		var name string
		if len(slice) > 1 {
			name = slice[1]
		} else {
			fmt.Println("Name not found.")
		}
		cmd.AttachByName(name)
	} else if strings.HasPrefix(in, "find") {
		inputSlice := strings.Split(in, " ")
		dataType := "all"
		targetVal := inputSlice[1]
		if len(inputSlice) < 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}
		if len(inputSlice) == 3 {
			targetVal = inputSlice[2]
			dataType = inputSlice[1]
		}
		foundAddr, _ := cmd.Find(targetVal, dataType)
		addrCache = foundAddr

	} else if strings.HasPrefix(in, "filter") {
		if len(addrCache) == 0 {
			fmt.Println("No previous results. ")
			return
		}
		slice := strings.Split(in, " ")
		if len(slice) == 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}

		foundAddr, err := cmd.Filter(slice[1], addrCache)
		if err != nil {
			fmt.Println(err)
		}
		addrCache = foundAddr

	} else if strings.HasPrefix(in, "patch") {
		slice := strings.Split(in, " ")
		if len(slice) == 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}

		err := cmd.Patch(slice[1], addrCache)
		if err != nil {
			fmt.Println(err)
		}

	} else if in == "detach" {
		if err := cmd.Detach(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else if strings.HasPrefix(in, "dump") {
		inputSlice := strings.Split(in, " ")
		beginAddr, err := parseAddr(inputSlice[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		endAddr, err := parseAddr(inputSlice[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := cmd.Dump(beginAddr, endAddr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else if in == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)

	} else if in == "" {

	} else {
		fmt.Println("Command not found.")
	}
}

func parseAddr(arg string) (int, error) {
	arg = strings.Replace(arg, "0x", "", 1)
	address, err := strconv.ParseInt(arg, 16, 64)
	if err == nil {
		return int(address), nil
	}
	address, err = strconv.ParseInt(arg, 10, 64)
	if err == nil {
		return int(address), nil
	}
	return 0, err
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "find   <int>", Description: "Search the specified integer."},
		{Text: "find   <datatype> <int>", Description: "Types can be specified are string, word, dword, qword."},
		{Text: "filter <int>", Description: "Filter previous search results that match the current search results."},
		{Text: "patch  <int>", Description: "Write the specified value on the address found by search."},
		{Text: "attach", Description: "Attach to the target process by pid."},
		{Text: "attachbyname", Description: "Attach to the target process by name."},
		{Text: "detach", Description: "Detach from the attached process."},
		{Text: "ps", Description: "Find the target process and if there is only one, specify it as the target."},
		{Text: "dump <begin addr> <end addr>", Description: "Display memory dump like hexdump"},
		{Text: "exit"},
	}
}

func main() {
	flag.StringVar(&specificPID, "pid", "", "Attach to a process with this pid")
	flag.Parse()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if specificPID != "" {
		appPID = specificPID
	}
	addrCache = []cmd.Found{}
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("medit: MEmory eDIT tool"),
		prompt.OptionPrefix("> "),
		prompt.OptionInputTextColor(prompt.Cyan),
		prompt.OptionPrefixTextColor(prompt.DarkBlue),
		prompt.OptionPreviewSuggestionTextColor(prompt.Green),
		prompt.OptionDescriptionTextColor(prompt.DarkGray),
	)
	p.Run()
}
