package statUtils

import (
	"github.com/shirou/gopsutil/v3/load"
)

type StatData struct {
	Name              string             `json:"name"`
	Ip                string             `json:"ip"`
	CpuUsage          float64            `json:"cpuUsage"`
	CpuLoad           *load.AvgStat      `json:"cpuLoad"`
	MemInfo           map[string]MemInfo `json:"memInfo"`
	Networks          []Network          `json:"network"`
	DiskInfos         []DiskInfo         `json:"diskInfo"`
	MaxGoroutines     int                `json:"maxGoroutines"`
	CurrentGoroutines int                `json:"currentGoroutines"`
	ServerType        int                `json:"serverType"`
	CreateTime        int64              `json:"createTime"`
}

type MemInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type Network struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
}
