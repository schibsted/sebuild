# Special Sources - specialsrcs

    specialsrcs[charmap:swedish.txt:charmap-swedish
                charmap:russian.txt:charmap-russian
                charmap:swedish-utf8.txt:charmap-swedish-utf8
    ]           
    
This argument specifies a method to compile a source file that the automatic
file type detection can't figure out how to compile. Each element in is of the
form: `<command>:<src>[,<src>]*:<target>[:<extravars>[,<extravars>]*]`

`command` is the ninja rule to build the target, `src` is a comma separated
list of sources, `target` is the resulting output. `extravars` is a comma
separated list of extra variables added to the build directive for ninja.

Notice that specialsrcs is mostly useless on its own, the script will not pick
up and do anything with the output from the script. You will need to specify
additional arguments that use the output somehow, usually with
[conf](../descriptors/install.md#conf) or [srcs](srcs.md).

The target can also be a comma separated list. In that case a single command
run should generate all the target files. The `$out` variable will have
all the output files in it, in the same order.

## Special Commands

Plugins are allowed to hijack specialsrcs by redirecting any specific commands
to themselves. Such specialsrcs should be documented in the plugin
documentation. Sebase does not currently have any builtin special rules of this
kind.

## Common Commands

The rules.ninja file bundled with sebuild contains some rules that are meant to
be used as commands in specialsrcs. You can further add your own either locally
via the [extravars argument](extravars.md) or in the [CONFIG
descriptor](../descriptors/config.md) via the `rules` argument there.

### concat

Simply concatenates the sources into the target, can also be used to copy files.

### touch

Runs the touch command on the output file, ignores inputs. Useful for
[godeps_rule](../descriptors/config.md#godeps-and-godeps_rule).

### gperf_switch

The source should be a C or C++ file. Inside these you can have
`GPERF_ENUM(x)` or `GPERF_ENUM_NOCASE(x)` markers to start a gperf switch.
Below it you add `switch (lookup_x(str, strlen(str)))` and inside you can use
`case GPERF_CASE("foo"):` to match str vs. "foo" via gperf.
The target header file must have the name `x.h` to match the argument given
to `GPREF_ENUM`. This header files contains definitions for the required macros
and also the inline lookup function.

### protocc, protoc-c

Used to run protobuf generating tools on a `.proto` file. The target file
name should match what's created by protoc.

### download

The source file should have one line with a URL and one line with the
sha256 sum of the file the URL specifies. It will be downloaded and checked.

### unzip, untar

Can be used to extract files from an archive, e.g. one downloaded via the
download command. Set the `file` extravar to the inside archive path to
be extracted.

### bunzip2

Runs bunzip2 to decompress a file.

### ronn

Runs the ronn command to generate man pages as used to generate the
seb man page.

### phony

A null ninja rule that can be used to indigate that the target is generated
by the same command that generates the source. Useful when another specialsrcs
generates two or more files, e.g. a c source and header.
