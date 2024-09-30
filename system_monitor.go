package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// Struct to hold network stats
type NetworkStats struct {
	BytesSent     uint64
	BytesRecv     uint64
	LastBytesSent uint64
	LastBytesRecv uint64
	LastUpdate    time.Time
}

func getProgressBar(percent float64, width int) string {
	filled := int(float64(width) * percent / 100)
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "-"
		}
	}
	return fmt.Sprintf("|%s| %.1f%%", bar, percent)
}

func centreText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	leftPad := (width - len(text)) / 2
	rightPad := width - len(text) - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

func formatBytes(bytes float64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%.2f B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", bytes/float64(div), "KMGTPE"[exp])
}

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Initialize CPU measurement
	cpu.Percent(0, false)

	// Initialize network stats
	var netStats NetworkStats
	initialIOCounters, _ := net.IOCounters(false)
	if len(initialIOCounters) > 0 {
		netStats = NetworkStats{
			LastBytesSent: initialIOCounters[0].BytesSent,
			LastBytesRecv: initialIOCounters[0].BytesRecv,
			LastUpdate:    time.Now(),
		}
	}

	// Print initial static content
	fmt.Print("\033[2J") // Clear screen
	fmt.Print("\033[H")  // Move cursor to home position
	fmt.Println(centreText("===== System Information =====", 80))
	fmt.Println(strings.Repeat("=", 80))

	hostInfo, _ := host.Info()
	fmt.Printf("%-40s%-40s\n", fmt.Sprintf("Hostname: %s", hostInfo.Hostname), fmt.Sprintf("OS: %s", hostInfo.OS))
	fmt.Printf("%-40s%-40s\n", fmt.Sprintf("Kernel: %s", hostInfo.KernelVersion), fmt.Sprintf("Uptime: %s", time.Duration(hostInfo.Uptime)*time.Second))
	fmt.Printf("%-40s%-40s\n", "CPU Usage:", "Memory Usage:")
	fmt.Printf("%-40s%-40s\n", "Disk Usage:", "Total Memory: ")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-40s%-40s\n", "Network Monitor", "Temperatures")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-40s%-40s\n", "Upload Speed:", "CPU Temperature:")
	fmt.Printf("%-40s%-40s\n", "Download Speed:", "GPU Temperature:")
	fmt.Printf("%-40s%-40s\n", "Total Uploaded:", "Disk Temperature:")
	fmt.Printf("%-40s%-40s\n", "Total Downloaded:", "Battery Temperature:")

	for {
		select {
		case <-ticker.C:
			// Move cursor to update dynamic content
			fmt.Print("\033[5;1H")
			fmt.Printf("%-40s%-40s\n", "CPU Usage:", "Memory Usage:")

			cpuPercent, _ := cpu.Percent(0, false)
			memory, _ := mem.VirtualMemory()
			disk, _ := disk.Usage("/")
			ioCounters, _ := net.IOCounters(false)
			temperatures, _ := host.SensorsTemperatures()

			if len(cpuPercent) > 0 {
				fmt.Printf("%-40s", getProgressBar(cpuPercent[0], 15))
			} else {
				fmt.Printf("%-40s", "N/A")
			}
			fmt.Printf("%-40s\n", getProgressBar(memory.UsedPercent, 15))
			fmt.Printf("%-40s%-40s\n", "Disk Usage:", "Total Memory: ")

			fmt.Printf("%-40s", getProgressBar(disk.UsedPercent, 15))
			fmt.Printf("%-40s\n", fmt.Sprintf("%.2f GB", float64(memory.Total)/(1024*1024*1024)))

			fmt.Println(strings.Repeat("=", 80))
			fmt.Printf("%-40s%-40s\n", "Network Monitor", "Temperatures")
			fmt.Println(strings.Repeat("-", 80))

			// Network and Temperature info
			if len(ioCounters) > 0 {
				now := time.Now()
				duration := now.Sub(netStats.LastUpdate).Seconds()

				uploadSpeed := float64(ioCounters[0].BytesSent-netStats.LastBytesSent) / duration
				downloadSpeed := float64(ioCounters[0].BytesRecv-netStats.LastBytesRecv) / duration

				fmt.Printf("%-40s", fmt.Sprintf("Upload Speed: %s/s", formatBytes(uploadSpeed)))
				if len(temperatures) > 0 {
					fmt.Printf("%-40s\n", fmt.Sprintf("CPU: %.1f°C", temperatures[0].Temperature))
				} else {
					fmt.Printf("%-40s\n", "CPU: N/A")
				}

				fmt.Printf("%-40s", fmt.Sprintf("Download Speed: %s/s", formatBytes(downloadSpeed)))
				if len(temperatures) > 1 {
					fmt.Printf("%-40s\n", fmt.Sprintf("GPU: %.1f°C", temperatures[1].Temperature))
				} else {
					fmt.Printf("%-40s\n", "GPU: N/A")
				}

				fmt.Printf("%-40s", fmt.Sprintf("Total Uploaded: %s", formatBytes(float64(ioCounters[0].BytesSent))))
				if len(temperatures) > 2 {
					fmt.Printf("%-40s\n", fmt.Sprintf("Disk: %.1f°C", temperatures[2].Temperature))
				} else {
					fmt.Printf("%-40s\n", "Disk: N/A")
				}

				fmt.Printf("%-40s", fmt.Sprintf("Total Downloaded: %s", formatBytes(float64(ioCounters[0].BytesRecv))))
				if len(temperatures) > 3 {
					fmt.Printf("%-40s\n", fmt.Sprintf("Battery: %.1f°C", temperatures[3].Temperature))
				} else {
					fmt.Printf("%-40s\n", "Battery: N/A")
				}

				// Update network stats for next iteration
				netStats.LastBytesSent = ioCounters[0].BytesSent
				netStats.LastBytesRecv = ioCounters[0].BytesRecv
				netStats.LastUpdate = now
			}

			// Update uptime
			fmt.Print("\033[3;41H")
			fmt.Printf("%-40s", fmt.Sprintf("Uptime: %s", time.Duration(hostInfo.Uptime)*time.Second))

			// Move cursor to the bottom
			fmt.Print("\033[15;1H")
			fmt.Println(centreText("Press Ctrl+C to exit", 80))

		}
	}
}
