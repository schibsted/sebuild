# Go Programs - GOPROG

    GOPROG(foo
    )

When used, the directory containing this Builddesc will be compiled into a go
binary. You can add the argument nocgo[] to disable cgo for this program
only. Other common arguments also work, for example `libs[]` can be useful
when compiling a Go program linking with C libraries.

When running the gobuild tool to compile go binaries and tests, a number of
ninja variables are checked. These can be set via a
[configvars](config.md#configvars) file or the [extravars
argument](../arguments/extravars.md).

Additionally, setting the `nocgo` [condition](../conditions.md) disables cgo
for all programs.

By default, dependency tracking is disabled for Go programs, since it can be
quite slow. See the [go_track_deps](gonfig.md#go_track_deps) CONFIG argument
for more details.

## Arguments

### nocgo

Use with an empty value, i.e. `nocgo[]`. Disables cgo for this binary.

### goarch

Sets the goarch to compile for, useful when compiling binaries that will be
deployed. Commonly set together with goos for a specific flavor.
Example:

	goos:release[linux]
	goarch:release[amd64]

### goos

Sets the goos to compile for. See [goarch](#goarch) for more information.

### gopkg

You can use GOPROG without having the go sources present in the same directory.
If you use the `gopkg` argument that Go package will be compiled instead of the
directory containing the Builddesc. If required, it will even be downloaded
first, thus you can use any package available for download and install it
to your local build.

	gopkg[github.com/schibsted/sebuild/v2/cmd/seb]

The gopkg path used must name a `main` package.

## Ninja Variables

### gobuild_flags
Flags added when executing `go build`. This can be used to for example add
build tags or other build options.

Defaults to the `GOBUILD_FLAGS` environment if set, otherwise empty string.

### gobuild_test_flags
Flags added when executing `go test`.

Defaults to the `GOBUILD_TEST_FLAGS` environment if set, otherwise empty string.

### cgo_enabled

Another way to control cgo. Defaults to the `CGO_ENABLED` environment variable,
which the seb binary might set itself if the `nocgo` condition is used.
