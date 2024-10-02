# System Monitor

## Overview

This Go program is a simple system monitor that provides real-time information about your computer's performance and resources. It displays various system metrics including CPU usage, memory usage, disk usage, network statistics, and temperature readings.

## Features

- Real-time updates of system metrics
- CPU usage monitoring
- Memory usage and total memory display
- Disk usage information
- Network monitoring (upload/download speeds and total data transferred)
- Temperature readings (CPU, GPU, Disk, Battery)
- Visual progress bars for easy interpretation of usage percentages
- Continuous updates with a clean, console-based interface


## Dependencies

go-fetch relies on the following external libraries:

- github.com/shirou/gopsutil/v3

Ensure you have these dependencies installed before building the project.

## Building

To build the project, navigate to the project directory and run:

```
go build
```

This will create an executable named `system-monitor` in the current directory.

## Notes

- The program updates every 500 milliseconds by default. You can adjust this by modifying the `ticker` interval in the `main()` function.
- Temperature readings may not be available on all systems. The program will display "N/A" for unavailable metrics.
