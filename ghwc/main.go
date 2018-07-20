package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jaypipes/ghw"
)

var (
	info *ghw.HostInfo
)

func main() {
	i, err := ghw.Host()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	info = i
	err = rootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCommand.AddCommand(memoryCommand)
	rootCommand.AddCommand(cpuCommand)
	rootCommand.AddCommand(blockCommand)
	rootCommand.AddCommand(topologyCommand)
	rootCommand.AddCommand(netCommand)
	rootCommand.AddCommand(gpuCommand)
	rootCommand.SilenceUsage = true
}

var rootCommand = &cobra.Command{
	Use:   "ghwc",
	Short: "ghwc - Discover hardware information.",
	Long:  "ghwc - Discover hardware information.",
	RunE:  showAll,
}

func showAll(cmd *cobra.Command, args []string) error {
	err := showMemory(cmd, args)
	if err != nil {
		return err
	}
	err = showCPU(cmd, args)
	if err != nil {
		return err
	}
	err = showBlock(cmd, args)
	if err != nil {
		return err
	}
	err = showTopology(cmd, args)
	if err != nil {
		return err
	}
	err = showNetwork(cmd, args)
	if err != nil {
		return err
	}
	err = showGPU(cmd, args)
	if err != nil {
		return err
	}
	return nil
}

var memoryCommand = &cobra.Command{
	Use:   "memory",
	Short: "Show memory information for the host system",
	RunE:  showMemory,
}

func showMemory(cmd *cobra.Command, args []string) error {
	mem := info.Memory
	fmt.Printf("%v\n", mem)
	return nil
}

var cpuCommand = &cobra.Command{
	Use:   "cpu",
	Short: "Show CPU information for the host system",
	RunE:  showCPU,
}

func showCPU(cmd *cobra.Command, args []string) error {
	cpu := info.CPU
	fmt.Printf("%v\n", cpu)

	for _, proc := range cpu.Processors {
		fmt.Printf(" %v\n", proc)
		for _, core := range proc.Cores {
			fmt.Printf("  %v\n", core)
		}
		if len(proc.Capabilities) > 0 {
			// pretty-print the (large) block of capability strings into rows
			// of 6 capability strings
			rows := int(math.Ceil(float64(len(proc.Capabilities)) / float64(6)))
			for row := 1; row < rows; row = row + 1 {
				rowStart := (row * 6) - 1
				rowEnd := int(math.Min(float64(rowStart+6), float64(len(proc.Capabilities))))
				rowElems := proc.Capabilities[rowStart:rowEnd]
				capStr := strings.Join(rowElems, " ")
				if row == 1 {
					fmt.Printf("  capabilities: [%s\n", capStr)
				} else if rowEnd < len(proc.Capabilities) {
					fmt.Printf("                 %s\n", capStr)
				} else {
					fmt.Printf("                 %s]\n", capStr)
				}
			}
		}
	}
	return nil
}

var blockCommand = &cobra.Command{
	Use:   "block",
	Short: "Show block storage information for the host system",
	RunE:  showBlock,
}

func showBlock(cmd *cobra.Command, args []string) error {
	block := info.Block
	fmt.Printf("%v\n", block)

	for _, disk := range block.Disks {
		fmt.Printf(" %v\n", disk)
		for _, part := range disk.Partitions {
			fmt.Printf("  %v\n", part)
		}
	}
	return nil
}

var topologyCommand = &cobra.Command{
	Use:   "topology",
	Short: "Show topology information for the host system",
	RunE:  showTopology,
}

func showTopology(cmd *cobra.Command, args []string) error {
	topology := info.Topology
	fmt.Printf("%v\n", topology)

	for _, node := range topology.Nodes {
		fmt.Printf(" %v\n", node)
		for _, cache := range node.Caches {
			fmt.Printf("  %v\n", cache)
		}
	}
	return nil
}

var netCommand = &cobra.Command{
	Use:   "net",
	Short: "Show network information for the host system",
	RunE:  showNetwork,
}

func showNetwork(cmd *cobra.Command, args []string) error {
	net := info.Network
	fmt.Printf("%v\n", net)

	for _, nic := range net.NICs {
		fmt.Printf(" %v\n", nic)
	}
	return nil
}

var gpuCommand = &cobra.Command{
	Use:   "gpu",
	Short: "Show graphics/GPU information for the host system",
	RunE:  showGPU,
}

func showGPU(cmd *cobra.Command, args []string) error {
	gpu := info.GPU
	fmt.Printf("%v\n", gpu)

	for _, card := range gpu.GraphicsCards {
		fmt.Printf(" %v\n", card)
	}
	return nil
}
