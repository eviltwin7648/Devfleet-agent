package utils

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
)

type MachineInfo struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
	TotalMem uint64 `json:"totalMem"`
}

// CollectMachineInfo gathers system metadata used during login, verify, heartbeat
func CollectMachineInfo() (MachineInfo, error) {
	hostname, _ := os.Hostname()

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return MachineInfo{}, err
	}

	return MachineInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostname,
		TotalMem: memInfo.Total,
	}, nil
}
