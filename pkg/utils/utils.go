package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/uuid"
)

func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown", err
	}
	return hostname, nil
}

func GetProcessID() int {
	return os.Getpid()
}

func GetSessionID() string {
	return uuid.New().String()
}

// 获取本地ip地址
func GetLocalIP() (string, error) {
	interfaces, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown", err
	}
	for _, inter := range interfaces {
		if ipnet, ok := inter.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "unknown", nil
}

func GetMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "unknown", err
	}
	for _, iface := range interfaces {
		// 跳过回环接口和非激活接口
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 {
			mac := iface.HardwareAddr.String()
			if mac != "" {
				return mac, nil
			}
		}
	}
	return "unknown", err
}

func GenerateClientID() string {
	// hostname+ip+sessionID
	hostname, err := GetHostname()
	if err != nil {
		return "unknown"
	}
	ip, err := GetLocalIP()
	if err != nil {
		return "unknown"
	}
	mac, err := GetMACAddress()
	if err != nil {
		return "unknown"
	}
	sessionID := GetSessionID()
	fmt.Println(hostname, ip, mac, sessionID)
	data := fmt.Sprintf("%s|%s|%s|%s", hostname, ip, mac, sessionID)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SystemInfo 系统信息结构
type SystemInfo struct {
	// 基本信息
	OS            string `json:"os"`             // 操作系统
	Arch          string `json:"arch"`           // 系统架构
	KernelName    string `json:"kernel_name"`    // 内核名称
	KernelVersion string `json:"kernel_version"` // 内核版本
	OSVersion     string `json:"os_version"`     // 操作系统版本
	Hostname      string `json:"hostname"`       // 主机名

	// 硬件信息
	CPUCores      int     `json:"cpu_cores"`       // CPU核数
	CPUModel      string  `json:"cpu_model"`       // CPU型号
	MemoryTotal   uint64  `json:"memory_total"`    // 总内存(字节)
	MemoryTotalGB float64 `json:"memory_total_gb"` // 总内存(GB)
	DiskTotal     uint64  `json:"disk_total"`      // 总硬盘(字节)
	DiskTotalGB   float64 `json:"disk_total_gb"`   // 总硬盘(GB)

	// 软件版本
	GoVersion      string `json:"go_version"`      // Go版本
	PythonVersion  string `json:"python_version"`  // Python版本
	Python3Version string `json:"python3_version"` // Python3版本
	ShellVersion   string `json:"shell_version"`   // Shell版本
}

// GetSystemInfo 获取完整系统信息
func GetSystemInfo() *SystemInfo {
	info := &SystemInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
	}

	// 获取主机名
	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	// 获取CPU信息
	info.CPUCores = runtime.NumCPU()
	info.CPUModel = getCPUModel()

	// 获取内核和操作系统版本
	info.KernelName, info.KernelVersion = getKernelInfo()
	info.OSVersion = getOSVersion()

	// 获取内存信息
	info.MemoryTotal = getMemoryTotal()
	info.MemoryTotalGB = float64(info.MemoryTotal) / (1024 * 1024 * 1024)

	// 获取硬盘信息
	info.DiskTotal = getDiskTotal()
	info.DiskTotalGB = float64(info.DiskTotal) / (1024 * 1024 * 1024)

	// 获取软件版本
	info.PythonVersion = getPythonVersion()
	info.Python3Version = getPython3Version()
	info.ShellVersion = getShellVersion()

	return info
}

// getCPUModel 获取CPU型号
func getCPUModel() string {
	switch runtime.GOOS {
	case "linux":
		return getCPUModelLinux()
	case "windows":
		return getCPUModelWindows()
	case "darwin":
		return getCPUModelDarwin()
	default:
		return "unknown"
	}
}

// getCPUModelLinux 获取Linux系统CPU型号
func getCPUModelLinux() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

// getCPUModelWindows 获取Windows系统CPU型号
func getCPUModelWindows() string {
	cmd := exec.Command("wmic", "cpu", "get", "name", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Name=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Name="))
		}
	}
	return "unknown"
}

// getCPUModelDarwin 获取macOS系统CPU型号
func getCPUModelDarwin() string {
	cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// getKernelInfo 获取内核信息
func getKernelInfo() (name, version string) {
	switch runtime.GOOS {
	case "linux":
		return getKernelInfoLinux()
	case "windows":
		return getKernelInfoWindows()
	case "darwin":
		return getKernelInfoDarwin()
	default:
		return "unknown", "unknown"
	}
}

// getKernelInfoLinux 获取Linux内核信息
func getKernelInfoLinux() (string, string) {
	// 获取内核名称
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	var kernelName string
	if err == nil {
		kernelName = strings.TrimSpace(string(output))
	} else {
		kernelName = "Linux"
	}

	// 获取内核版本
	cmd = exec.Command("uname", "-r")
	output, err = cmd.Output()
	var kernelVersion string
	if err == nil {
		kernelVersion = strings.TrimSpace(string(output))
	} else {
		kernelVersion = "unknown"
	}

	return kernelName, kernelVersion
}

// getKernelInfoWindows 获取Windows内核信息
func getKernelInfoWindows() (string, string) {
	cmd := exec.Command("cmd", "/c", "ver")
	output, err := cmd.Output()
	if err != nil {
		return "Windows NT", "unknown"
	}

	version := strings.TrimSpace(string(output))
	return "Windows NT", version
}

// getKernelInfoDarwin 获取macOS内核信息
func getKernelInfoDarwin() (string, string) {
	cmd := exec.Command("uname", "-v")
	output, err := cmd.Output()
	var version string
	if err == nil {
		version = strings.TrimSpace(string(output))
	} else {
		version = "unknown"
	}

	return "Darwin", version
}

// getOSVersion 获取操作系统版本
func getOSVersion() string {
	switch runtime.GOOS {
	case "linux":
		return getOSVersionLinux()
	case "windows":
		return getOSVersionWindows()
	case "darwin":
		return getOSVersionDarwin()
	default:
		return "unknown"
	}
}

// getOSVersionLinux 获取Linux系统版本
func getOSVersionLinux() string {
	// 尝试读取 /etc/os-release
	if content, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(content), "\n")
		var name, version string

		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				name = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			}
			if strings.HasPrefix(line, "NAME=") && name == "" {
				name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
			}
			if strings.HasPrefix(line, "VERSION=") {
				version = strings.Trim(strings.TrimPrefix(line, "VERSION="), "\"")
			}
		}

		if name != "" {
			if version != "" && !strings.Contains(name, version) {
				return fmt.Sprintf("%s %s", name, version)
			}
			return name
		}
	}

	// 尝试读取 /etc/issue
	if content, err := os.ReadFile("/etc/issue"); err == nil {
		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 && lines[0] != "" {
			return strings.TrimSpace(lines[0])
		}
	}

	return "Linux"
}

// getOSVersionWindows 获取Windows系统版本
func getOSVersionWindows() string {
	cmd := exec.Command("powershell", "-Command", "(Get-WmiObject -Class Win32_OperatingSystem).Caption")
	output, err := cmd.Output()
	if err != nil {
		return "Windows"
	}
	return strings.TrimSpace(string(output))
}

// getOSVersionDarwin 获取macOS系统版本
func getOSVersionDarwin() string {
	cmd := exec.Command("sw_vers", "-productName")
	nameOutput, err1 := cmd.Output()

	cmd = exec.Command("sw_vers", "-productVersion")
	versionOutput, err2 := cmd.Output()

	if err1 == nil && err2 == nil {
		name := strings.TrimSpace(string(nameOutput))
		version := strings.TrimSpace(string(versionOutput))
		return fmt.Sprintf("%s %s", name, version)
	}

	return "macOS"
}

// getMemoryTotal 获取总内存大小
func getMemoryTotal() uint64 {
	switch runtime.GOOS {
	case "linux":
		return getMemoryTotalLinux()
	case "windows":
		return getMemoryTotalWindows()
	case "darwin":
		return getMemoryTotalDarwin()
	default:
		return 0
	}
}

// getMemoryTotalLinux 获取Linux系统内存大小
func getMemoryTotalLinux() uint64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if kb, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
					return kb * 1024 // 转换为字节
				}
			}
		}
	}
	return 0
}

// getMemoryTotalWindows 获取Windows系统内存大小
func getMemoryTotalWindows() uint64 {
	cmd := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "TotalPhysicalMemory=") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "TotalPhysicalMemory="))
			if memory, err := strconv.ParseUint(value, 10, 64); err == nil {
				return memory
			}
		}
	}
	return 0
}

// getMemoryTotalDarwin 获取macOS系统内存大小
func getMemoryTotalDarwin() uint64 {
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	value := strings.TrimSpace(string(output))
	if memory, err := strconv.ParseUint(value, 10, 64); err == nil {
		return memory
	}
	return 0
}

// getDiskTotal 获取总硬盘大小
func getDiskTotal() uint64 {
	switch runtime.GOOS {
	case "linux":
		return getDiskTotalLinux()
	case "windows":
		return getDiskTotalWindows()
	case "darwin":
		return getDiskTotalDarwin()
	default:
		return 0
	}
}

// getDiskTotalLinux 获取Linux系统硬盘大小
func getDiskTotalLinux() uint64 {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return 0
	}

	// 计算根分区大小
	totalSize := uint64(stat.Blocks) * uint64(stat.Bsize)

	// 尝试获取所有挂载点的总大小
	if totalAll := getDiskTotalAllLinux(); totalAll > totalSize {
		return totalAll
	}

	return totalSize
}

// getDiskTotalAllLinux 获取Linux所有磁盘总大小
func getDiskTotalAllLinux() uint64 {
	cmd := exec.Command("lsblk", "-b", "-d", "-n", "-o", "SIZE")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	var total uint64
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if size, err := strconv.ParseUint(line, 10, 64); err == nil {
				total += size
			}
		}
	}
	return total
}

// getDiskTotalWindows 获取Windows系统硬盘大小
func getDiskTotalWindows() uint64 {
	cmd := exec.Command("wmic", "diskdrive", "get", "size", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	var total uint64
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Size=") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "Size="))
			if size, err := strconv.ParseUint(value, 10, 64); err == nil {
				total += size
			}
		}
	}
	return total
}

// getDiskTotalDarwin 获取macOS系统硬盘大小
func getDiskTotalDarwin() uint64 {
	cmd := exec.Command("diskutil", "list", "-plist")
	output, err := cmd.Output()
	if err != nil {
		// 备用方法：使用df命令
		cmd = exec.Command("df", "-k", "/")
		output, err = cmd.Output()
		if err != nil {
			return 0
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 2 {
				if kb, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
					return kb * 1024
				}
			}
		}
		fmt.Println("diskutil list -plist failed, using df -k /")

		return 0
	}
	fmt.Println(string(output))
	// 这里可以解析plist输出，暂时使用简单方法
	var stat syscall.Statfs_t
	err = syscall.Statfs("/", &stat)
	if err != nil {
		return 0
	}

	return uint64(stat.Blocks) * uint64(stat.Bsize)
}

// getPythonVersion 获取Python版本
func getPythonVersion() string {
	// 尝试不同的Python命令
	pythonCommands := []string{"python3", "python", "python3.11", "python3.10", "python3.9", "python3.8"}

	for _, cmd := range pythonCommands {
		if version := tryGetPythonVersion(cmd); version != "" {
			return version
		}
	}

	return "not installed"
}

// tryGetPythonVersion 尝试获取指定Python命令的版本
func tryGetPythonVersion(pythonCmd string) string {
	cmd := exec.Command(pythonCmd, "--version")
	output, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		// Python 2输出到stderr，Python 3输出到stdout
		if version != "" {
			return version
		}
	}

	// 尝试stderr
	cmd = exec.Command(pythonCmd, "--version")
	if output, err := cmd.CombinedOutput(); err == nil {
		version := strings.TrimSpace(string(output))
		if strings.HasPrefix(version, "Python") {
			return version
		}
	}

	return ""
}

// getPython3Version 获取Python3版本
func getPython3Version() string {
	cmd := exec.Command("python3", "--version")
	output, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		if version != "" {
			return version
		}
	}
	return ""
}

// getShellVersion 获取Shell版本
func getShellVersion() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}

	// 获取shell名称
	shellName := shell
	if idx := strings.LastIndex(shell, "/"); idx != -1 {
		shellName = shell[idx+1:]
	}

	// 尝试获取版本
	cmd := exec.Command(shellName, "--version")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) > 0 {
			return fmt.Sprintf("%s (%s)", shellName, strings.TrimSpace(lines[0]))
		}
	}

	return shellName
}

// GetSystemInfoMap 获取系统信息并转换为map格式
func GetSystemInfoMap() map[string]interface{} {
	info := GetSystemInfo()

	return map[string]interface{}{
		"os":              info.OS,
		"arch":            info.Arch,
		"kernel_name":     info.KernelName,
		"kernel_version":  info.KernelVersion,
		"os_version":      info.OSVersion,
		"hostname":        info.Hostname,
		"cpu_cores":       info.CPUCores,
		"cpu_model":       info.CPUModel,
		"memory_total":    info.MemoryTotal,
		"memory_total_gb": fmt.Sprintf("%.2f GB", info.MemoryTotalGB),
		"disk_total":      info.DiskTotal,
		"disk_total_gb":   fmt.Sprintf("%.2f GB", info.DiskTotalGB),
		"go_version":      info.GoVersion,
		"python_version":  info.PythonVersion,
		"python3_version": info.Python3Version,
		"shell_version":   info.ShellVersion,
	}
}
