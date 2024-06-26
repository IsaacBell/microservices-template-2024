package util

import (
	"core/pkg/influx"
	"fmt"
	"os"
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

// func handle(w http.ResponseWriter, r *http.Request) {
// 	// Your API Logic
// }

// http.Handle("/api", moesifmiddleware.MoesifMiddleware(http.HandlerFunc(handle), moesifOption))

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
		PrintLnInColor(AnsiColorCyan, "Benchmarks: \n->  ", processName, "ran in "+time.Since(start).String())
		ticker.Stop()
		influx.LogRuntime(processName, time.Since(start))
	}
}

func RecordSystemMetrics() {
	os.Stdout.Sync()

	for {
		fmt.Println(AnsiColorYellow)
		defer fmt.Println(AnsiColorReset)
		fmt.Println("-->recording system performance metrics...")

		metrics := make(map[string]map[string]interface{})
		timestamp := time.Now().Local().Format(time.RFC3339)
		metrics[timestamp] = make(map[string]interface{})
		metrics[timestamp]["label"] = timestamp
		metrics[timestamp]["timestamp"] = timestamp

		// CPU metrics
		cpuPercent, err := cpu.Percent(0, false)
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving CPU percentage: ", err, "\n")
		} else {
			metrics[timestamp]["cpu.usage"] = cpuPercent[0]
		}

		cpuLoad, err := load.Avg()
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving CPU load:\n. -> ", err)
		} else {
			metrics[timestamp]["cpu.load.1m"] = cpuLoad.Load1
			metrics[timestamp]["cpu.load.5m"] = cpuLoad.Load5
			metrics[timestamp]["cpu.load.15m"] = cpuLoad.Load15
		}

		// Memory metrics
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving virtual memory stats:\n. -> ", err)
		} else {
			metrics[timestamp]["memory.used"] = vmStat.Used
			metrics[timestamp]["memory.total"] = vmStat.Total
			metrics[timestamp]["memory.used_percent"] = vmStat.UsedPercent
		}

		swapStat, err := mem.SwapMemory()
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving swap memory stats:\n. -> ", err)
		} else {
			metrics[timestamp]["swap.used"] = swapStat.Used
			metrics[timestamp]["swap.total"] = swapStat.Total
			metrics[timestamp]["swap.used_percent"] = swapStat.UsedPercent
		}

		// Disk metrics
		diskStat, err := disk.Usage("/")
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving disk usage stats:\n. -> ", err)
		} else {
			metrics[timestamp]["disk.used"] = diskStat.Used
			metrics[timestamp]["disk.total"] = diskStat.Total
			metrics[timestamp]["disk.used_percent"] = diskStat.UsedPercent
		}

		diskIOCounters, err := disk.IOCounters()
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving disk IO counters:\n. -> ", err)
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
			PrintLnInColor(AnsiColorRed, "Caught error retrieving network IO counters:\n. -> ", err)
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
			PrintLnInColor(AnsiColorRed, "Caught error retrieving processes:\n. -> ", err)
		} else {
			metrics[timestamp]["process.count"] = len(processes)
			for _, p := range processes {
				cpu, err := p.CPUPercent()
				if err != nil {
					PrintLnInColor(AnsiColorRed, "Caught error retrieving CPU percentage for process: ", p.Pid, "\n-->err: ", err)
				} else {
					metrics[timestamp][fmt.Sprintf("process.%d.cpu", p.Pid)] = cpu
				}

				mem, err := p.MemoryInfo()
				if err != nil {
					PrintLnInColor(
						AnsiColorRed, "Caught error retrieving memory info for process: ",
						AnsiColorYellow, p.Pid,
						AnsiColorGray, "\n-->err: ",
						AnsiColorRed, err)
				} else {
					metrics[timestamp][fmt.Sprintf("process.%d.memory", p.Pid)] = mem.RSS
				}
			}
		}

		// System metrics
		uptime, err := host.Uptime()
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving system uptime:\n. -> ", err)
		} else {
			metrics[timestamp]["system.uptime"] = uptime
		}

		bootTime, err := host.BootTime()
		if err != nil {
			PrintLnInColor(AnsiColorRed, "Caught error retrieving system boot time:\n. -> ", err)
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
			PrintLnInColor(AnsiColorRed, "Caught error retrieving CPU times:\n. -> ", err)
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
			PrintLnInColor(AnsiColorRed, "Caught error logging system metrics:\n. -> ", err)
		}

		// Sleep for a specific interval before collecting metrics again
		time.Sleep(30 * time.Second)
	}
}
