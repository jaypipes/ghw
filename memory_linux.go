// +build linux

package ghw

import (
    "bufio"
    "compress/gzip"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "strconv"
)

func memFillInfo(info *MemoryInfo) error {
    tpb := memTotalPhysicalBytes()
    if tpb < 1 {
        return fmt.Errorf("Could not determine total physical bytes of memory")
    }
    info.TotalPhysicalBytes = tpb
    return nil
}

// System log lines will look similar to the following:
// ... kernel: [0.000000] Memory: 24633272K/25155024K ...
var (
    syslogMemLineRe = regexp.MustCompile("Memory:\\s+\\d+K\\/(\\d+)K")
)

func memTotalPhysicalBytes() int64 {
    // In Linux, the total physical memory can be determined by looking at the
    // output of dmidecode, however dmidecode requires root privileges to run,
    // so instead we examine the system logs for startup information containing
    // total physical memory and cache the results of this.
    findPhysicalKb := func (line string) int64 {
        matches :=  syslogMemLineRe.FindStringSubmatch(line)
        if len(matches) == 2 {
            i, err := strconv.Atoi(matches[1])
            if err != nil {
                return -1
            }
            return int64(i * 1024)
        }
        return -1
    }

    // /var/log will contain a file called syslog and 0 or more files called
    // syslog.$NUMBER or syslog.$NUMBER.gz containing system log records. We
    // search each, stopping when we match a system log record line that
    // contains physical memory information.
    logDir := "/var/log"
    logFiles, err := ioutil.ReadDir(logDir)
    if err != nil {
        return -1
    }
    for _, file := range logFiles {
        if strings.HasPrefix(file.Name(), "syslog") {
            fullPath := filepath.Join(logDir, file.Name())
            unzip := strings.HasSuffix(file.Name(), ".gz")
            var r io.ReadCloser
            r, err = os.Open(fullPath)
            if err != nil {
                return -1
            }
            defer r.Close()
            if unzip {
                r, err = gzip.NewReader(r)
                if err != nil {
                    return -1
                }
            }

            scanner := bufio.NewScanner(r)
            for scanner.Scan() {
                line := scanner.Text()
                size := findPhysicalKb(line)
                if size > 0 {
                    return size
                }
            }
        }
    }
    return -1
}
