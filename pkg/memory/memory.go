package memory

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// GetProcess found process based on name
func GetProcess(nameProcess string) (uint32, error) {
	snapshotHandle, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, fmt.Errorf("error return of the package syscall on func CreateToolhelp32Snapshot: %s", err.Error())
	}
	if snapshotHandle == syscall.InvalidHandle {
		return 0, fmt.Errorf("invalid Handle")
	}
	var processEntry syscall.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))
	err = syscall.Process32First(snapshotHandle, &processEntry)
	if err != nil {
		return 0, fmt.Errorf("error return of the package syscall on func Process32First: %s", err.Error())
	}
	for {
		exeFileName := MaxPathToString(processEntry.ExeFile)
		if strings.EqualFold(exeFileName, nameProcess) {
			syscall.CloseHandle(snapshotHandle)
			return processEntry.ProcessID, nil
		}
		if syscall.Process32Next(snapshotHandle, &processEntry) != nil {
			break
		}
	}
	return 0, fmt.Errorf("no found process %s", nameProcess)
}

func GetModuleBaseAddress(processID uint32, module string) (uint32, error) {
	snapshotHandle, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPMODULE|syscall.TH32CS_SNAPMODULE32, processID)
	if err != nil {
		return 0, fmt.Errorf("error return of the package syscall on func CreateToolhelp32Snapshot: %s", err.Error())
	}
	if snapshotHandle == syscall.InvalidHandle {
		return 0, fmt.Errorf("invalid Handle")
	}
	var moduleEntry windows.ModuleEntry32
	moduleEntry.Size = uint32(unsafe.Sizeof(moduleEntry))
	windows.Module32First(windows.Handle(snapshotHandle), &moduleEntry)
	if err != nil {
		return 0, fmt.Errorf("error return on func Module32FirstW: %s", err.Error())
	}
	for {
		maxModule := MaxModuleToString(moduleEntry.Module)
		if strings.EqualFold(maxModule, module) {
			return uint32(moduleEntry.ModBaseAddr), nil
		}
		if windows.Module32Next(windows.Handle(snapshotHandle), &moduleEntry) != nil {
			break
		}
	}
	return 0, fmt.Errorf("no found module %s", module)
}

// uint16ToString Converts a [MAX_PATH]uint16 array to a Go string
func MaxPathToString(array [syscall.MAX_PATH]uint16) string {
	// Finds the index of the first null byte
	var i int
	for i = 0; i < len(array) && array[i] != 0; i++ {
	}
	return syscall.UTF16ToString(array[:i])
}

// uint16ToString Converts an array [MAX_MODULE_NAME32 + 1]uint16 to a Go string
func MaxModuleToString(array [windows.MAX_MODULE_NAME32 + 1]uint16) string {
	return syscall.UTF16ToString(array[:])
}
