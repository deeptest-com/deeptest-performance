package statUtils

import (
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	gonet "net"
	"runtime"
	"strings"
	"time"
)

func GetAll() *StatData {
	data := new(StatData)

	data.Name = GetHostName()
	data.Ip = GetIp()

	data.CpuUsage = GetCpuUsed()
	data.CpuLoad = GetCpuLoad()
	data.MemInfo = GetMemInfo()
	data.Networks = GetNetwork()
	data.DiskInfos = GetDiskInfo()
	data.CreateTime = time.Now().Unix()
	data.CurrentGoroutines = runtime.NumGoroutine()

	return data
}

func GetCpuUsed() float64 {
	percent, _ := cpu.Percent(time.Second, false) // false表示总使用率，true为单核
	return percent[0]
}

func GetCpuLoad() (info *load.AvgStat) {
	info, _ = load.Avg()
	return
}

func GetMemInfo() (memInfo map[string]MemInfo) {
	memInfo = map[string]MemInfo{}

	// virtual
	memVir := MemInfo{}

	memInfoVir, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	memVir.Total = memInfoVir.Total
	memVir.Free = memInfoVir.Free
	memVir.Used = memInfoVir.Used
	memVir.UsedPercent = memInfoVir.UsedPercent
	memInfo["virtual"] = memVir

	// swap
	memSwap := MemInfo{}

	memInfoSwap, err := mem.SwapMemory()
	if err != nil {
		return
	}

	memSwap.Total = memInfoSwap.Total
	memSwap.Free = memInfoSwap.Free
	memSwap.Used = memInfoSwap.Used
	memSwap.UsedPercent = memInfoSwap.UsedPercent
	memInfo["swap"] = memSwap

	return
}

func GetHostName() string {
	hostInfo, _ := host.Info()
	return hostInfo.Hostname
}

func GetDiskInfo() (diskInfoList []DiskInfo) {
	disks, err := disk.Partitions(false)
	if err != nil {
		return
	}

	for _, v := range disks {
		diskInfo := DiskInfo{}
		info, err := disk.Usage(v.Device)
		if err != nil {
			continue
		}
		diskInfo.Total = info.Total
		diskInfo.Free = info.Free
		diskInfo.Used = info.Used
		diskInfo.UsedPercent = info.UsedPercent
		diskInfoList = append(diskInfoList, diskInfo)
	}

	return
}

func GetNetwork() (networkList []Network) {
	netIOs, _ := net.IOCounters(true)
	if netIOs == nil {
		return
	}

	for _, netIO := range netIOs {
		network := Network{}
		network.Name = netIO.Name
		network.BytesSent = netIO.BytesSent
		network.BytesRecv = netIO.BytesRecv
		network.PacketsSent = netIO.PacketsSent
		network.PacketsRecv = netIO.PacketsRecv
		networkList = append(networkList, network)
	}

	return
}

func GetIp() (ret string) {
	conn, err := gonet.Dial("udp", "8.8.8.8:53")
	if err != nil {
		logUtils.Errorf("udp error：%s", err.Error())
		return
	}

	localAddr := conn.LocalAddr().(*gonet.UDPAddr)
	ret = strings.Split(localAddr.String(), ":")[0]

	return
}
