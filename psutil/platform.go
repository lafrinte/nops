package psutil

import (
	"github.com/lafrinte/nops/fs"
	l "github.com/lafrinte/nops/log"
	"github.com/shirou/gopsutil/v3/host"
	"reflect"
	"strings"
)

var hostInfo *host.InfoStat
var log = l.DefaultLogger().GetLogger()

func init() {
	i, err := host.Info()
	if err != nil {
		log.Error().Err(err).Msg("failed get host info in psutil.init")
		return
	}

	hostInfo = i
}

func get(prop string) any {
	if hostInfo == nil {
		return ""
	}

	obj := reflect.ValueOf(hostInfo).Elem()

	if field := obj.FieldByName(prop); field.IsValid() {
		return field.Interface()
	}

	return ""
}

func GetOS() string {
	return get("OS").(string)
}

func GetHostname() string {
	return get("Hostname").(string)
}

func GetPlatform() string {
	return get("Platform").(string)
}

func GetPlatformVersion() string {
	return get("PlatformVersion").(string)
}

func GetArch() string {
	return get("KernelArch").(string)
}

func IsWindows() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "windows")
}

func IsLinux() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "linux")
}

func IsDarwin() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "darwin")
}

func IsSunOs() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "sunos")
}

func IsSmartOs() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "joyent_")
}

func IsFreeBsd() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "freebsd")
}

func IsOpenBsd() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "openbsd")
}

func IsAix() bool {
	return strings.Contains(strings.ToLower(GetPlatform()), "aix")
}

func IsArch64() bool {
	return strings.Contains(strings.ToLower(GetArch()), "arch64")
}

func IsX86() bool {
	return strings.Contains(strings.ToLower(GetArch()), "x86_64")
}

func IsContainer() bool {
	dockerEnv := "/.dockerenv"
	cgroup := "/proc/self/cgroup"

	return fs.IsFile(dockerEnv) || fs.IsFile(cgroup)
}
