module github.com/jaypipes/ghw

go 1.18

require (
	github.com/StackExchange/wmi v1.2.1
	github.com/ghodss/yaml v1.0.0
	github.com/jaypipes/pcidb v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/safchain/ethtool v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v0.0.3
	howett.net/plist v1.0.0
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/spf13/pflag v1.0.2 // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/safchain/ethtool => github.com/fromanirh/ethtool-ioctl v0.2.1-0.20220510154755-7ca867c90cb0
