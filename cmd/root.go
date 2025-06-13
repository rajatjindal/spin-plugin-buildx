package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type buildFlagsType struct {
	up           bool
	componentIds []string
	from         string
	help         bool
	debug        bool
}

var buildFlags = buildFlagsType{}

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
	rootCmd.Flags().BoolVar(&buildFlags.debug, "debug", false, "if enabled, print dagger logs")

	rootCmd.Flags().BoolVarP(&buildFlags.up, "up", "u", false, "Run the application after building")
	rootCmd.Flags().StringSliceVarP(&buildFlags.componentIds, "component-id", "c", nil, "Component ID to build. This can be specified multiple times. The default is all components")
	rootCmd.Flags().StringVarP(&buildFlags.from, "from", "f", "", "The application to build. This may be a manifest (spin.toml) file, or a directory containing a spin.toml file. If omitted, it defaults to \"spin.toml\"")
}

func buildx() error {
	ctx := context.TODO()

	dagger, err := checkDagger()
	if err != nil {
		return err
	}

	args := []string{"call"}
	if !buildFlags.debug {
		args = append(args, "-s")
	} else {
		// open terminal when build in container fails
		args = append(args, "-i")
	}

	if buildFlags.up {
		args = append(args, []string{
			"-m=github.com/rajatjindal/daggerverse/wasi@main",
			// "-m=../../rajatjindal/daggerverse/wasi",
			"up",
			"--source=.",
			fmt.Sprintf("--args=%q", strings.Join([]string{"--build"}, ",")),
			"up",
		}...)
	} else {
		args = append(args, []string{
			"-m=github.com/rajatjindal/daggerverse/wasi@main",
			// "-m=../../rajatjindal/daggerverse/wasi",
			"build",
			"--source=.",
			"directory",
			"--path=/app",
			"export",
			"--path=.",
		}...)
	}

	cmd := exec.CommandContext(ctx, dagger, args...)

	// DO NOT SEND TRACES TO DAGGER CLOUD
	// Setting XDG_CONFIG_HOME is just a hack, otherwise dagger
	// sends the traces to the cloud.
	cmd.Env = append(os.Environ(), "DAGGER_NO_NAG=1", "DO_NOT_TRACK=1", "XDG_CONFIG_HOME=/foo")

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
