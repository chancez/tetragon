// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package syscallinfo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

// SyscallArgInfo is the name and the type (as string) of a syscall argument
type SyscallArgInfo struct {
	Name string
	Type string
}

// SyscallArgs is the arguments for a given syscall
type SyscallArgs []SyscallArgInfo

// syscall table: name -> []SyscallArgs
var sysargsInfo map[string]SyscallArgs

func init() {
	// parse syscall table
	err := json.Unmarshal(syscalls_, &sysargsInfo)
	if err != nil {
		panic(err)
	}
}

// GetSyscallName returns the name of a syscall based on its i d
func GetSyscallName(sysID int) string {
	if name, ok := syscallNames[sysID]; ok {
		// NB: remove the sys_ prefix, so name can be used with GetSyscallInfo
		return name[4:]
	}
	return ""
}

// GetSyscallArgs returns the arguments of a system call
func GetSyscallArgs(name string) (SyscallArgs, bool) {
	if args, ok := sysargsInfo[name]; ok {
		ret := make([]SyscallArgInfo, len(args))
		copy(ret, args)
		return SyscallArgs(ret), true
	}
	return nil, false
}

// Proto returns a string representing a  prototype for the system call
func (sai SyscallArgs) Proto(name string) string {
	args := make([]string, 0, len(sai))
	for i := range sai {
		args = append(args, fmt.Sprintf("%s %s", sai[i].Type, sai[i].Name))
	}
	return fmt.Sprintf("long %s(%s)", name, strings.Join(args, ", "))
}
