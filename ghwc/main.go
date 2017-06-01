package main

import (
    "fmt"
    "os"

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
    rootCommand.SilenceUsage = true
}

var rootCommand = &cobra.Command{
    Use: "ghwc",
    Short: "ghwc - Discover hardware information.",
    Long: "ghwc - Discover hardware information.",
    RunE: showAll,
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
    return nil
}

var memoryCommand = &cobra.Command{
    Use: "memory",
    Short: "Show memory information for the host system",
    RunE: showMemory,
}

func showMemory(cmd *cobra.Command, args []string) error {
    mem := info.Memory
    fmt.Printf("%v\n", mem)
    return nil
}

var cpuCommand = &cobra.Command{
    Use: "cpu",
    Short: "Show CPU information for the host system",
    RunE: showCPU,
}

func showCPU(cmd *cobra.Command, args []string) error {
    cpu := info.CPU
    fmt.Printf("%v\n", cpu)
    return nil
}

var blockCommand = &cobra.Command{
    Use: "block",
    Short: "Show block storage information for the host system",
    RunE: showBlock,
}

func showBlock(cmd *cobra.Command, args []string) error {
    block := info.Block
    fmt.Printf("%v\n", block)
    return nil
}
