package ghw

import (
	"time"

	"golang.org/x/sys/windows/registry"
)

// Packages is...
func (ctx *context) packagesFillInfo(info *PackagesInfo) error {
	// Opening main key to find all installed packages
	mainKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return err
	}
	defer mainKey.Close()
	// Collecting sub keys for current main key
	subKeys, err := mainKey.ReadSubKeyNames(-1)
	if err != nil {
		return err
	}
	// Collecting installed software data from subkeys
	for _, subKeyLabel := range subKeys {
		subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\`+subKeyLabel, registry.QUERY_VALUE)
		if err != nil {
			return err
		}
		defer subKey.Close()
		label, _, _ := subKey.GetStringValue("DisplayName")
		version, _, _ := subKey.GetStringValue("DisplayVersion")
		installDate, _, _ := subKey.GetStringValue("InstallDate")
		// Parsing installDate if set
		if installDate != "" {
			parsedTime, _ := time.Parse("20060102", installDate)
			installDate = parsedTime.Format(time.RFC3339)
		}
		// Appending converted package if the package has at least a name...
		if label != "" {
			info.Installed = append(info.Installed, &PackageInfo{
				Label:       label,
				Version:     version,
				InstallDate: installDate,
			})
		}
	}
	return nil
}
