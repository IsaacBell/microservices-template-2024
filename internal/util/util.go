package util

import (
	"fmt"
	"microservices-template-2024/pkg/influx"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

// defer Benchmark("my-process")()
func Benchmark(processName string) func() {
	start := time.Now()
	done := make(chan bool)
	ticker := time.NewTicker(1500 * time.Millisecond)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Slow query: ", processName, "\n->Tick at", t)
			}
		}
	}()

	return func() {
		fmt.Println(processName, ": ran in "+time.Since(start).String())
		ticker.Stop()
		influx.LogRuntime(processName, time.Since(start))
	}
}

func RecordSystemMetrics() {
	for {
		metrics := make(map[string]map[string]interface{})
		timestamp := time.Now().Local().String()
		metrics[timestamp]["label"] = timestamp
		metrics[timestamp]["timestamp"] = timestamp

		// CPU metrics
		cpuPercent, err := cpu.Percent(0, false)
		if err != nil {
			fmt.Println("Error retrieving CPU percentage:", err)
		} else {
			metrics[timestamp]["cpu.usage"] = cpuPercent[0]
		}

		cpuLoad, err := load.Avg()
		if err != nil {
			fmt.Println("Error retrieving CPU load:", err)
		} else {
			metrics[timestamp]["cpu.load.1m"] = cpuLoad.Load1
			metrics[timestamp]["cpu.load.5m"] = cpuLoad.Load5
			metrics[timestamp]["cpu.load.15m"] = cpuLoad.Load15
		}

		// Memory metrics
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Error retrieving virtual memory stats:", err)
		} else {
			metrics[timestamp]["memory.used"] = vmStat.Used
			metrics[timestamp]["memory.total"] = vmStat.Total
			metrics[timestamp]["memory.used_percent"] = vmStat.UsedPercent
		}

		swapStat, err := mem.SwapMemory()
		if err != nil {
			fmt.Println("Error retrieving swap memory stats:", err)
		} else {
			metrics[timestamp]["swap.used"] = swapStat.Used
			metrics[timestamp]["swap.total"] = swapStat.Total
			metrics[timestamp]["swap.used_percent"] = swapStat.UsedPercent
		}

		// Disk metrics
		diskStat, err := disk.Usage("/")
		if err != nil {
			fmt.Println("Error retrieving disk usage stats:", err)
		} else {
			metrics[timestamp]["disk.used"] = diskStat.Used
			metrics[timestamp]["disk.total"] = diskStat.Total
			metrics[timestamp]["disk.used_percent"] = diskStat.UsedPercent
		}

		diskIOCounters, err := disk.IOCounters()
		if err != nil {
			fmt.Println("Error retrieving disk IO counters:", err)
		} else {
			for idx, counter := range diskIOCounters {
				metrics[timestamp]["disk["+idx+"].read_bytes"] = counter.ReadBytes
				metrics[timestamp]["disk["+idx+"].write_bytes"] = counter.WriteBytes
				metrics[timestamp]["disk["+idx+"].read_count"] = counter.ReadCount
				metrics[timestamp]["disk["+idx+"].write_count"] = counter.WriteCount
			}
		}

		// Network metrics
		netIOCounters, err := net.IOCounters(false)
		if err != nil {
			fmt.Println("Error retrieving network IO counters:", err)
		} else {
			metrics[timestamp]["network.bytes_sent"] = netIOCounters[0].BytesSent
			metrics[timestamp]["network.bytes_recv"] = netIOCounters[0].BytesRecv
			metrics[timestamp]["network.packets_sent"] = netIOCounters[0].PacketsSent
			metrics[timestamp]["network.packets_recv"] = netIOCounters[0].PacketsRecv
			metrics[timestamp]["network.err_in"] = netIOCounters[0].Errin
			metrics[timestamp]["network.err_out"] = netIOCounters[0].Errout
		}

		// Process-level metrics
		processes, err := process.Processes()
		if err != nil {
			fmt.Println("Error retrieving processes:", err)
		} else {
			metrics[timestamp]["process.count"] = len(processes)
			for _, p := range processes {
				cpu, err := p.CPUPercent()
				if err != nil {
					fmt.Printf("Error retrieving CPU percentage for process %d: %v", p.Pid, err)
				} else {
					metrics[timestamp][fmt.Sprintf("process.%d.cpu", p.Pid)] = cpu
				}

				mem, err := p.MemoryInfo()
				if err != nil {
					fmt.Printf("Error retrieving memory info for process %d: %v", p.Pid, err)
				} else {
					metrics[timestamp][fmt.Sprintf("process.%d.memory", p.Pid)] = mem.RSS
				}
			}
		}

		// System metrics
		uptime, err := host.Uptime()
		if err != nil {
			fmt.Println("Error retrieving system uptime:", err)
		} else {
			metrics[timestamp]["system.uptime"] = uptime
		}

		bootTime, err := host.BootTime()
		if err != nil {
			fmt.Println("Error retrieving system boot time:", err)
		} else {
			metrics[timestamp]["system.boot_time"] = bootTime
		}

		// Runtime metrics
		metrics[timestamp]["runtime.goroutines"] = runtime.NumGoroutine()
		var gcStats runtime.MemStats
		runtime.ReadMemStats(&gcStats)
		metrics[timestamp]["runtime.gc.pause_total"] = gcStats.PauseTotalNs
		metrics[timestamp]["runtime.gc.pause_count"] = gcStats.NumGC
		metrics[timestamp]["runtime.gc.heap_alloc"] = gcStats.HeapAlloc
		metrics[timestamp]["runtime.gc.heap_sys"] = gcStats.HeapSys

		times, err := cpu.Times(true)
		if err != nil {
			fmt.Println("Error retrieving CPU times:", err)
		} else {
			cpus := make(map[string]interface{})
			metrics[timestamp]["cpu.count"] = len(times)
			for i, currCpu := range times {
				indexStr := strconv.Itoa(i)
				cpus["times.CPU["+indexStr+"]"] = currCpu.CPU
				cpus["times.User["+indexStr+"]"] = currCpu.User
				cpus["times.System["+indexStr+"]"] = currCpu.System
				cpus["times.Idle["+indexStr+"]"] = currCpu.Idle
				cpus["times.Nice["+indexStr+"]"] = currCpu.Nice
				cpus["times.Iowait["+indexStr+"]"] = currCpu.Iowait
				cpus["times.Irq["+indexStr+"]"] = currCpu.Irq
				cpus["times.Softirq["+indexStr+"]"] = currCpu.Softirq
				cpus["times.Steal["+indexStr+"]"] = currCpu.Steal
				cpus["times.Guest["+indexStr+"]"] = currCpu.Guest
				cpus["times.GuestNice["+indexStr+"]"] = currCpu.GuestNice
			}
			metrics[timestamp]["cpu.times"] = cpus
		}

		err = influx.LogSystemMetrics(metrics)
		if err != nil {
			fmt.Println("Error logging system metrics:", err)
		}

		// Sleep for a specific interval before collecting metrics again
		time.Sleep(10 * time.Second)
	}
}
