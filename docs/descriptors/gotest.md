# Go tests - GOTEST

    GOTEST(foo
    )

Uses `go test -c` to produce a binary called foo.test that's by default
installed to the /gotest destination directory. This binary can be executed to
run the tests, but do note that some go unit tests might assume that they're
run from a specific directory or have other expectations.

In addition, you can also run the go tests directly via ninja. The advantage of
using ninja here to run go test is that it adds include paths etc. for cgo
tests that can otherwise be difficult to figure out.

To allow to run go test via ninja, for each GOTEST descriptor a target
`build/flavor/gotest/name` is created which you can give to ninja to run the
test. In addition, you can also run `go test -cover` to generate an html file
with this target:

    build/<flavor>/gocover/<name>.html

And also `go test -bench` with

    build/<flavor>/gobench/<name>

You can override the default of running all benchmarks by adding a `benchflags`
argument to the `GOTEST` directive and putting a regexp there, possibly also
adding additional flags.

The gobuild tool used by GOTEST checks a number of ninja variables when
executed. These are described on the [GOPROG page](goprog.md). The `gopkg`
argument described there also works for GOTEST.

As an advanced feature, GOTEST targets are also automatically collected in
variables, see the
[collect_target_var argument](../arguments/collect-target-var.md) for more more
information about this.

Unfortunately, the go coverage html generating does not currently work in other
go modules than the main one. There is no easy way to solve this so it's kept
as a known issue.

To easily generate all go coverage report, the target

    build/<flavor>/gocover

maps to the list of reports and can be used to generate all of them.

By default, dependency tracking is disabled for Go tests, since it can be
quite slow. See the [go_track_deps](gonfig.md#go_track_deps) CONFIG argument
for more details.
