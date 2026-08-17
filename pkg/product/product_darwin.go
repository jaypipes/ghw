// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"howett.net/plist"

	"github.com/jaypipes/ghw/internal/config"
	"github.com/jaypipes/ghw/internal/log"
	"github.com/jaypipes/ghw/pkg/util"
)

type hardwareOverview struct {
	MachineName  string `json:"machine_name"`  // MacBook Pro
	MachineModel string `json:"machine_model"` // Mac14,7
	ModelNumber  string `json:"model_number"`  // Z1AB0001CD/A
	SerialNumber string `json:"serial_number"` // C02ABC123DEF
	PlatformUUID string `json:"platform_UUID"` // some-uuid-foo-bar
}

func run(ctx context.Context, args ...string) ([]byte, error) {
	log.Debug(ctx, "running %q", strings.Join(args, " "))
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", args[0], err)
	}
	return out, nil
}

// platformExpertDevice returns the properties of the I/O Registry's
// IOPlatformExpertDevice node, the macOS equivalent of the DMI/SMBIOS system
// information Linux exposes under /sys/class/dmi/id.
func platformExpertDevice(ctx context.Context) (map[string]any, error) {
	out, err := run(ctx, "ioreg",
		"-a",      // use XML output
		"-d", "1", // limit device tree output depth to the matched node
		"-r",                           // root device tree at matched node
		"-c", "IOPlatformExpertDevice", // match by class
	)
	if err != nil {
		return nil, err
	}
	return parsePlatformExpertDevice(out)
}

func parsePlatformExpertDevice(out []byte) (map[string]any, error) {
	var data []map[string]any
	if len(out) > 0 {
		if _, err := plist.Unmarshal(out, &data); err != nil {
			return nil, fmt.Errorf("ioreg unmarshal for IOPlatformExpertDevice failed: %w", err)
		}
	}
	if len(data) == 0 {
		return nil, errors.New("ioreg returned no IOPlatformExpertDevice node")
	}
	return data[0], nil
}

func hardwareInfo(ctx context.Context) (*hardwareOverview, error) {
	out, err := run(ctx, "system_profiler", "-json", "SPHardwareDataType")
	if err != nil {
		return nil, err
	}
	return parseHardwareInfo(out)
}

func parseHardwareInfo(out []byte) (*hardwareOverview, error) {
	var data struct {
		Items []hardwareOverview `json:"SPHardwareDataType"`
	}
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("system_profiler SPHardwareDataType unmarshal failed: %w", err)
	}
	if len(data.Items) == 0 {
		return nil, errors.New("system_profiler returned no SPHardwareDataType item")
	}
	return &data.Items[0], nil
}

// property returns the named I/O Registry property as a string. Device tree
// properties are NUL-terminated (and NUL-padded) byte arrays, everything else
// is already a plain string.
func property(node map[string]any, key string) string {
	switch val := node[key].(type) {
	case string:
		return strings.TrimSpace(val)
	case []byte:
		return strings.TrimSpace(string(bytes.TrimRight(val, "\x00")))
	}
	return ""
}

func (i *Info) load(ctx context.Context) error {
	if !config.ToolsEnabled(ctx) {
		return errors.New("DisableTools=true on darwin disables product support entirely.")
	}

	dev, err := platformExpertDevice(ctx)
	if err != nil {
		return err
	}

	// Query system_profiler for the machine name
	hw, err := hardwareInfo(ctx)
	if err != nil {
		log.Warn(ctx, "%s\n", err)
	}

	i.fill(dev, hw)

	return nil
}

// fill maps the I/O Registry properties and the system_profiler(8) overview,
// which may be nil, onto the product fields.
func (i *Info) fill(dev map[string]any, hw *hardwareOverview) {
	i.Name = property(dev, "model")
	i.Vendor = property(dev, "manufacturer")
	i.SerialNumber = property(dev, "IOPlatformSerialNumber")
	if i.SerialNumber == "" {
		i.SerialNumber = property(dev, "serial-number")
	}
	i.UUID = property(dev, "IOPlatformUUID")
	// Apple's part number is the number itself ("Z1AB0001C") followed by the
	// region it was sold in ("D/A").
	if sku := property(dev, "model-number"); sku != "" {
		i.SKU = sku + property(dev, "region-info")
	}
	i.Version = property(dev, "version")

	if hw != nil {
		i.Family = hw.MachineName
		if i.Name == "" {
			i.Name = hw.MachineModel
		}
		if i.SerialNumber == "" {
			i.SerialNumber = hw.SerialNumber
		}
		if i.UUID == "" {
			i.UUID = hw.PlatformUUID
		}
		if i.SKU == "" {
			i.SKU = hw.ModelNumber
		}
	}

	for _, field := range []*string{
		&i.Family, &i.Name, &i.Vendor, &i.SerialNumber, &i.UUID, &i.SKU, &i.Version,
	} {
		if *field == "" {
			*field = util.UNKNOWN
		}
	}
}
