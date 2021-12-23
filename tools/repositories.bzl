"""Workspace rules (tools/repositories)"""

load("@rules_haskell//haskell:cabal.bzl", "stack_snapshot")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

def rules_haskell_worker_dependencies(**stack_kwargs):
    """Provide all repositories that are necessary for `rules_haskell`'s tools to
    function.
    """
    excludes = native.existing_rules().keys()

    if "rules_haskell_worker_dependencies" not in excludes:
        stack_snapshot(
            name = "rules_haskell_worker_dependencies",
            packages = [
                "base",
                "bytestring",
                "filepath",
                "ghc",
                "ghc-paths",
                "microlens",
                "process",
                "profunctors-5.5.2",
                "proto-lens-0.7.0.0",
                "proto-lens-runtime-0.7.0.0",
                "text",
                "vector",
            ],
            snapshot = "lts-18.0",
            **stack_kwargs
        )

def bazel_binaries_for_integraion_testing():
    http_file(
        name = "bazel_bin_linux",
        executable = True,
        sha256 = "0eb2e378d2782e7810753e2162245ad1179c1bb12f848c692b4a595b4edf779b",
        urls = ["https://github.com/bazelbuild/bazel/releases/download/4.1.0/bazel-4.1.0-linux-x86_64"],
    )

    http_file(
        name = "bazel_bin_darwin",
        executable = True,
        sha256 = "74d93848f0c9d592e341e48341c53c87e3cb304a54a2a1ee9cff3df422f0b23c",
        urls = ["https://github.com/bazelbuild/bazel/releases/download/4.2.1/bazel-4.2.1-darwin-x86_64"],
    )

    http_file(
        name = "bazel_bin_windows",
        executable = True,
        sha256 = "7b2077af7055b421fe31822f83c3c3c15e36ff39b69560ba2472dde92dd45b46",
        urls = ["https://github.com/bazelbuild/bazel/releases/download/4.1.0/bazel-4.1.0-windows-x86_64.exe"],
    )
