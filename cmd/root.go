package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	remoteModule = "github.com/rajatjindal/daggerverse/wasi@main"
	localModule  = "../../rajatjindal/daggerverse/wasi"
)

var localdev = false

type buildFlagsType struct {
	up           bool
	from         string
	debug        bool
	componentIds []string
}

var buildFlags = buildFlagsType{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "buildx",
	Short: "spin buildx is a spin plugin that improves UX when dealing with multiple toolchain versions for spin apps",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if buildFlags.up {
			err := runSpinUp(ctx)
			if err != nil {
				fmt.Println("Error", err)
				os.Exit(1)
			}

			return
		}

		err := runSpinBuild(ctx)
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

func runSpinUp(ctx context.Context) error {
	dagger, args, err := getDefaultDaggerCmdAndArgs()
	if err != nil {
		return err
	}

	args = append(args, []string{
		"up",
		"--source=.",
		"up",
	}...)

	_, err = runAndGetStdout(exec.CommandContext(ctx, dagger, args...))
	if err != nil {
		return err
	}

	return nil
}

func runSpinBuild(ctx context.Context) error {
	dagger, args, err := getDefaultDaggerCmdAndArgs()
	if err != nil {
		return err
	}

	args = append(args, []string{
		"build",
		"--source=.",
		"directory",
		"--path=/app",
		"export",
		"--path=.",
	}...)

	_, err = runAndGetStdout(exec.CommandContext(ctx, dagger, args...))
	if err != nil {
		return err
	}

	return nil
}

// func buildxx() error {
// 	ctx := context.TODO()

// 	dagger, err := checkDagger()
// 	if err != nil {
// 		return err
// 	}

// 	args := []string{"call"}
// 	if !buildFlags.debug {
// 		args = append(args, "-s")
// 	} else {
// 		// open terminal when build in container fails
// 		args = append(args, "-i")
// 	}

// 	if buildFlags.up {
// 		args = append(args, []string{
// 			"-m=github.com/rajatjindal/daggerverse/wasi@main",
// 			// "-m=../../rajatjindal/daggerverse/wasi",
// 			"up",
// 			"--source=.",
// 			fmt.Sprintf("--args=%q", strings.Join([]string{"--build"}, ",")),
// 		}...)
// 	} else {
// 		args = append(args, []string{
// 			"-m=github.com/rajatjindal/daggerverse/wasi@main",
// 			// "-m=../../rajatjindal/daggerverse/wasi",
// 			"build",
// 			"--source=.",
// 			"directory",
// 			"--path=/app",
// 			"export",
// 			"--path=.",
// 		}...)
// 	}

// 	cmd := exec.CommandContext(ctx, dagger, args...)

// 	// DO NOT SEND TRACES TO DAGGER CLOUD
// 	// Setting XDG_CONFIG_HOME is just a hack, otherwise dagger
// 	// sends the traces to the cloud.
// 	cmd.Env = append(os.Environ(), "DAGGER_NO_NAG=1", "DO_NOT_TRACK=1", "XDG_CONFIG_HOME=/foo")

// 	return run(cmd)
// }

// func run(cmd *exec.Cmd) error {
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stdout
// 	cmd.Stdin = os.Stdin

// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("failed to run cmd %v\n", err)
// 	}

// 	return nil
// }
//
// func getBuildEnv(ctx context.Context) (string, error) {
// 	startTime := time.Now()
// 	logrus.Info("setting up build env...")

// 	dagger, args, err := getDefaultDaggerCmdAndArgs()
// 	if err != nil {
// 		return "", err
// 	}

// 	args = append(args, []string{
// 		"build-env",
// 		"--source=.",
// 	}...)

// 	stdout, err := runAndGetStdout(exec.CommandContext(ctx, dagger, args...))
// 	if err != nil {
// 		return "", err
// 	}

// 	logrus.Infof("setting up build env... (completed) - took %.2fs", time.Since(startTime).Seconds())

// 	buildEnvId, err := runSpinBuild(ctx, stdout)
// 	if err != nil {
// 		fmt.Println("Error", err)
// 		os.Exit(1)
// 	}

// 	return buildEnvId, nil
// }

// func runSpinBuild(ctx context.Context, buildEnvId string) (string, error) {
// 	startTime := time.Now()
// 	logrus.Info("running spin build...")

// 	dagger, args, err := getDefaultDaggerCmdAndArgs()
// 	if err != nil {
// 		return "", err
// 	}

// 	args = append(args, []string{
// 		"build-in-ctr",
// 		fmt.Sprintf("--env-id=%s", buildEnvId),
// 	}...)

// 	ctrId, err := runAndGetStdout(exec.CommandContext(ctx, dagger, args...))
// 	if err != nil {
// 		return "", err
// 	}

// 	logrus.Infof("running build... (completed) - took %.2fs", time.Since(startTime).Seconds())
// 	return ctrId, nil
// }

// func exportBuildArtifacts(ctx context.Context, buildEnvId string) (string, error) {
// 	startTime := time.Now()
// 	logrus.Info("exporting build artifacts...")

// 	dagger, args, err := getDefaultDaggerCmdAndArgs()
// 	if err != nil {
// 		return "", err
// 	}

// 	args = append(args, []string{
// 		"get-app-dir-in-ctr",
// 		fmt.Sprintf("--env-id=%s", buildEnvId),
// 		"export",
// 		"--path=.",
// 	}...)

// 	ctrId, err := runAndGetStdout(exec.CommandContext(ctx, dagger, args...))
// 	if err != nil {
// 		return "", err
// 	}

// 	logrus.Infof("exporting build artifacts... (completed) - took %.2fs", time.Since(startTime).Seconds())
// 	return ctrId, nil
// }

// func runSpinUp(ctx context.Context, buildEnvId string) (string, error) {
// 	startTime := time.Now()
// 	logrus.Info("running spin up...")

// 	dagger, args, err := getDefaultDaggerCmdAndArgs()
// 	if err != nil {
// 		return "", err
// 	}

// 	args = append(args, []string{
// 		"up-in-ctr",
// 		fmt.Sprintf("--env-id=%s", buildEnvId),
// 	}...)

// 	// TODO(rajatjindal): run in background and stream stdout/stderr
// 	ctrId, err := runAndGetStdout(exec.CommandContext(ctx, dagger, args...))
// 	if err != nil {
// 		return "", err
// 	}

// 	logrus.Infof("running spin up... (completed) - ran for %.2fs", time.Since(startTime).Seconds())
// 	return ctrId, nil
// }

func runAndGetStdout(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run cmd %v\n", err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func getDefaultDaggerCmdAndArgs() (string, []string, error) {
	dagger, err := checkDagger()
	if err != nil {
		return "", nil, err
	}

	// disable dagger traces
	// TODO(rajatjindal): allow user to enable them
	os.Setenv("DAGGER_NO_NAG", "1")
	os.Setenv("DO_NOT_TRACK", "1")
	os.Setenv("XDG_CONFIG_HOME", "/foo")

	currentModule := remoteModule
	if localdev {
		currentModule = localModule
	}

	args := []string{"call", fmt.Sprintf("-m=%s", currentModule)}
	if buildFlags.debug {
		args = append(args, "-i") // open terminal when build in container fails
	}

	return dagger, args, nil
}

func checkDagger() (string, error) {
	path, err := exec.LookPath("dagger")
	if err != nil {
		return "", fmt.Errorf("dagger is not found. Please install using https://docs.dagger.io/install")
	}

	return path, nil
}
