package integration_testing

import (
        "bytes"
        "fmt"
        "os"
        "os/exec"
        "runtime"
        "strings"
        "testing"
        "github.com/bazelbuild/rules_go/go/tools/bazel"
        "github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M, workspace string) {
        if err := ParseArgs(); err != nil {
                fmt.Fprint(os.Stderr, err)
                return
        }
	fmt.Fprintf(os.Stderr, "Testing context: %v\n", Context)
	defer exec.Command(Context.BazelBinary, "shutdown")

        bazel_testing.TestMain(m, bazel_testing.Args{
                Main: workspace + GenerateBazelrc(),
        })
}

func AssertOutput(t *testing.T, output []byte, expected string) {
        if string(output) != expected {
                t.Fatalf("output of bazel process is invalid.\nExpected: %v\n, Actual: %v\n", expected, string(output))
        }
}

var Context struct {
        Nixpkgs bool
        BazelBinary string
}

func ParseArgs() error {
        bazelPath := ""
        for _, arg := range os.Args {
                if strings.HasPrefix(arg, "nixpkgs=") {
                        fmt.Sscanf(arg, "nixpkgs=%t", &Context.Nixpkgs)
                } else if strings.HasPrefix(arg, "bazel_bin=") {
                        fmt.Sscanf(arg, "bazel_bin=%s", &bazelPath)
                }
        }
        bazelAbsPath, err := bazel.Runfile(bazelPath)
        Context.BazelBinary = bazelAbsPath
        return err
}

func GenerateBazelrc() string {
        bazelrc := "-- .bazelrc --\n"
        if Context.Nixpkgs {
                bazelrc += `
build --host_platform=@io_tweag_rules_nixpkgs//nixpkgs/platforms:host
build --incompatible_enable_cc_toolchain_resolution
`
        } else if runtime.GOOS == "windows" {
                bazelrc += `
build --crosstool_top=@rules_haskell_ghc_windows_amd64//:cc_toolchain
`
        }
        return bazelrc
}

func BazelCmd(bazelPath string, args ...string) *exec.Cmd {
        insertBazelFlags := func (flags ...string) []string {
                for i, arg := range args {
                        switch arg {
                                case
                                    "build",
                                    "test",
                                    "run":
                                    return append(append(append([]string{}, args[:i + 1]...), flags...), args[i + 1:]...)
                        }
                }
                return args
        }

        cmd := exec.Command(bazelPath)
        if bazel_testing.OutputUserRoot != "" {
                cmd.Args = append(cmd.Args, "--output_user_root="+bazel_testing.OutputUserRoot)
        }
        cmd.Args = append(cmd.Args, insertBazelFlags("--announce_rc", "-s", "--toolchain_resolution_debug=true")...)
	// It's important value of $HOME to be invariant between different integration test runs
        // and to be writable directory for bazel test. Probably TEST_TMPDIR is a valid choice
        // but documentation is not clear about it's default value
        // cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", os.Getenv("TEST_TMPDIR")))
        cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", os.TempDir()))
        if runtime.GOOS == "darwin" {
                cmd.Env = append(cmd.Env, "BAZEL_USE_CPP_ONLY_TOOLCHAIN=1")
        }
        for _, e := range os.Environ() {
                // Filter environment variables set by the bazel test wrapper script.
                // These confuse recursive invocations of Bazel.
                if strings.HasPrefix(e, "TEST_") || strings.HasPrefix(e, "RUNFILES_") {
                        continue
                }
                cmd.Env = append(cmd.Env, e)
        }
        fmt.Fprintf(os.Stderr, "bazel cmd: %v\n", cmd.Args)
        return cmd
}

func RunBazel(bazelPath string, args ...string) error {
        cmd := BazelCmd(bazelPath, args...)

        buf := &bytes.Buffer{}
        cmd.Stderr = buf
        err := cmd.Run()
        if eErr, ok := err.(*exec.ExitError); ok {
                eErr.Stderr = buf.Bytes()
                err = &bazel_testing.StderrExitError{Err: eErr}
        }
        return err
}

func BazelOutput(bazelPath string, args ...string) ([]byte, error) {
        cmd := BazelCmd(bazelPath, args...)
        stdout := &bytes.Buffer{}
        stderr := &bytes.Buffer{}
        cmd.Stdout = stdout
        cmd.Stderr = stderr
        err := cmd.Run()
        if eErr, ok := err.(*exec.ExitError); ok {
                eErr.Stderr = stderr.Bytes()
                err = &bazel_testing.StderrExitError{Err: eErr}
        }
        fmt.Fprintf(os.Stderr, "bazel stderr: %v\n", stderr)
        return stdout.Bytes(), err
}
