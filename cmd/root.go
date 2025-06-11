package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	debug bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "buildx",
	Short: "spin buildx is a spin plugin that improves UX when dealing with multiple toolchain versions for spin apps",
	Run: func(cmd *cobra.Command, args []string) {
		err := buildx()
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&debug, "debug", false, "if enabled, print dagger logs")
}

func buildx() error {
	ctx := context.TODO()

	dagger, err := checkDagger()
	if err != nil {
		return err
	}

	debugFlag := "-s"
	if debug {
		debugFlag = ""
	}

	cmd := exec.CommandContext(ctx, dagger, []string{
		"call",
		debugFlag,
		"-m=github.com/rajatjindal/daggerverse/wasi@main",
		"build",
		"--source=.",
		"export",
		"--path=.",
	}...)

	// DO NOT SEND TRACES TO DAGGER CLOUD
	cmd.Env = append(cmd.Environ(), "DAGGER_NO_NAG=1", "SHUTUP=1")

	return run(cmd)
}

func run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run cmd %v\n", err)
	}

	return nil
}

func checkDagger() (string, error) {
	path, err := exec.LookPath("dagger")
	if err != nil {
		return "", fmt.Errorf("dagger is not found. Please install using https://docs.dagger.io/install")
	}

	return path, nil
}
