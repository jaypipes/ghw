package cpu

import (
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
	"strings"
)

var sysctlOutput = make(map[string]string)

func (i *Info) load() error {
	err := populateSysctlOutput()
	if err != nil {
		return errors.Wrap(err, "unable to populate sysctl map")
	}

	i.TotalCores = getTotalCores()
	i.TotalThreads = getTotalThreads()
	i.Processors = getProcessors()

	return nil
}

// getProcessors some more info https://developer.apple.com/documentation/kernel/1387446-sysctlbyname/determining_system_capabilities
func getProcessors() []*Processor {
	p := make([]*Processor, getProcTopCount())
	for i, _ := range p {
		p[i] = new(Processor)
		p[i].Vendor = sysctlOutput[fmt.Sprintf("hw.perflevel%s.name", strconv.Itoa(i))]
		p[i].Model = getVendor()
		p[i].NumCores = getNumberCoresFromPerfLevel(i)
		p[i].Capabilities = getCapabilities()
	}
	return p
}

// valid for ARM see https://developer.apple.com/documentation/kernel/1387446-sysctlbyname/determining_instruction_set_characteristics
func getCapabilities() []string {
	caps := []string{}
	for k, v := range sysctlOutput {
		if strings.HasPrefix(k, "hw.optional.arm.") {
			isEnabled, _ := strconv.Atoi(v)
			if isEnabled == 1 {
				caps = append(caps, k)
			}
		}
	}

	return caps
}

// populateSysctlOutput to populate a map to quickly retrieve values later
func populateSysctlOutput() error {
	// get sysctl output
	o, err := exec.Command("sysctl", "-a").CombinedOutput()
	if err != nil {
		return err
	}

	// clean up and store sysctl output
	oS := strings.Split(string(o), "\n")
	for _, l := range oS {
		if l != "" {
			s := strings.SplitN(l, ":", 2)
			k, v := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
			sysctlOutput[k] = v
		}
	}

	return nil
}

func getNumberCoresFromPerfLevel(i int) uint32 {
	key := fmt.Sprintf("hw.perflevel%s.physicalcpu", strconv.Itoa(i))
	nCores := sysctlOutput[key]
	return stringToUint32(nCores)
}

func getVendor() string {
	v := sysctlOutput["machdep.cpu.brand_string"]
	return v
}

func getProcTopCount() int {
	pC, ok := sysctlOutput["hw.nperflevels"]
	if !ok {
		// most likely intel so no performance/efficiency core seperation
		return 1
	}
	i, _ := strconv.Atoi(pC)
	return i
}

// num of physical cores
func getTotalCores() uint32 {
	nCores := sysctlOutput["hw.physicalcpu"]
	return stringToUint32(nCores)
}

func getTotalThreads() uint32 {
	nThreads := sysctlOutput["machdep.cpu.thread_count"]
	return stringToUint32(nThreads)
}

func stringToUint32(s string) uint32 {
	o, _ := strconv.ParseUint(s, 10, 0)
	return uint32(o)
}
