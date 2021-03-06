// Copyright 2018-2019 Schibsted

package assets

const RulesNinja = `
rule cxx
    command = $cxx $picflag $flavor_cflags $cflags $cxxflags $cwarnflags $gcov_copts $copts $srcopts -I. -I$incdir $includes $defines -MMD -MF $out.d -MT $out -c $in -o $out
    depfile = $out.d
    description = C++ compile $out

rule cc
    command = $cc $picflag $flavor_cflags $cflags $conlyflags $cwarnflags $gcov_copts $copts $srcopts -I. -I$incdir $includes $defines -MMD -MF $out.d -MT $out -c $in -o $out
    depfile = $out.d
    description = C compile $out

rule cxx_analyse
    command = rm -rf $out ; mkdir -p $out ; clang++ $analyser_flags $picflag $flavor_cflags $cflags $cxxflags $cwarnflags $gcov_copts $copts $srcopts -I. -I$incdir $includes $defines $in -o $out
    description = C++ analyse $out

rule cc_analyse
    command = rm -rf $out ; mkdir -p $out ; clang $analyser_flags $picflag $flavor_cflags $cflags $conlyflags $cwarnflags $gcov_copts $copts $srcopts -I. -I$incdir $includes $defines $in -o $out
    description = C analyse $out

rule copy_analyse
    command = seb -tool copy-analyse $out $in
    description = Copy analyse to $out

rule final_analyse
    command = seb -tool copy-analyse -finalize $out $in
    description = Final analyse in $out

link_wrapper = seb -tool link
rule linkxx
    command = $link_wrapper $cxx $ldflags $ldopts $gcov_ldopts -L $libdir -o $out @$out.rsp $ldlibs
    description = link $out
    rspfile = $out.rsp
    rspfile_content = $in

rule link
    command = $link_wrapper $cc $ldflags $ldopts $gcov_ldopts -L $libdir -o $out @$out.rsp $ldlibs
    description = link $out
    rspfile = $out.rsp
    rspfile_content = $in

rule partiallink
    command = $link_wrapper ld -r -o $out @$out.rsp
    description = partially linking $out
    rspfile = $out.rsp
    rspfile_content = $in

rule ar
    command = ar cr $out $in
    description = ar library $out

rule flexx
    command = flex -+ -o$out $in
    description = lex $out

rule yaccxx
    command = bison -y -d --debug -o $out $in
    description = yacc $out

rule install_conf
    command = rm -f $out ; install -p -m0644 $in $out
    description = copy $in to $out

rule install_script
    command = rm -f $out ; install -p -m0755 $in $out
    description = copy $in to $out

rule install_header
    command = seb -tool header-install $in $out
    description = copy header to $out

rule install_py
    command = seb -tool python-install $in $out
    description = python copy $in to $out

rule install_php
    command = php -l $in > /dev/null && install -m0644 $in $out
    description = php syntax check and copy $in to $out

rule copy_go
    command = seb -tool go-install $in $out
    description = go fmt check and copy $in to $out

rule concat
    command = cat $in > $out
    description = concatenate $in to $out

rule touch
    command = touch $out

rule gperf_enum
    command = seb -tool gperf-enum enum $in $out
    description = gperf_enum $out

rule gperf_switch
    command = seb -tool gperf-enum source $in $out
    description = gperf_enum $out

rule gperf
    command = gperf -L ANSI-C --output-file=$out $in
    description = gperf $out

rule build_version
    command = echo \#define BUILD_VERSION \"$buildversion\" > $out

rule generate_ninjas
    command = BUILD_BUILD_FROM_NINJA=1 CC="$cc" BUILDTOOLDIR="$buildtooldir" $build_build
    generator=

# Since go build runs with custom pkgdirs install dependency packages
# they can't be run in parallel, only allow one at a time for these modes.
pool gobuilds_lib
    depth = 1
pool gobuilds_piclib
    depth = 1

# The gobuild tool uses some environment variables. We allow these to be set
# in configvars ninja files. They will default to the normal environment
# variables if not set there however.  This is for dependencies to work more
# properly as configvars script changes retrigger builds but environment
# variables do not.
# Note that for gobuild the depfile is only used if enabled. By default
# the commands are always run and instead use the Go build cache.
gobuild_tool=GOBUILD_FLAGS=$gobuild_flags GOBUILD_TEST_FLAGS=$gobuild_test_flags CGO_ENABLED=$cgo_enabled seb -tool gobuild

rule gobuild
    command = GOOS="$goos" GOARCH="$goarch" $gobuild_tool -pkg="$gopkg" -cflags="-I $incdir $includes $platform_includes" -ldflags="-L $libdir $ldlibs" -mode=$gomode -pkgdir="$builddir" "$in" "$out" "$objdir/depfile-$gomode"
    depfile = $objdir/depfile-$gomode
    description = building go $gomode $out from $in

rule gobuildlib
    command = GOOS="$goos" GOARCH="$goarch" $gobuild_tool -pkg="$gopkg" -cflags="$picflag -I $incdir $includes $platform_includes" -ldflags="-L $libdir $ldlibs" -mode=$gomode -pkgdir="$builddir" $gonoinit "$in" "$out" "$depfile"
    depfile = $depfile
    description = building go library $out from $in
    pool = gobuilds_$gomode

rule gotest
    command = $gobuild_tool -pkg="$gopkg" -cflags="-I $incdir $includes $platform_includes" -ldflags="-L $libdir $ldlibs" -mode=test -pkgdir="$builddir" "$in"
    description = testing go package in $in

rule gobench
    command = $gobuild_tool -pkg="$gopkg" -cflags="-I $incdir $platform_includes" -ldflags="-L $libdir $ldlibs" -mode=bench -pkgdir="$builddir" "$in" "$benchflags"
    description = benching go package in $in

rule gocover
    command = $gobuild_tool -pkg="$gopkg" -cflags="-I $incdir $includes $platform_includes" -ldflags="-L $libdir $ldlibs" -mode=cover -pkgdir="$builddir" "$in" "$out" "$objdir/depfile-cover"
    depfile = $objdir/depfile-cover
    description = testing coverage of go package in $in

rule gocover_html
    command = $gobuild_tool -pkg="$gopkg" -cflags="-I $incdir $includes $platform_includes" -ldflags="-L $libdir $ldlibs" -mode=cover_html -pkgdir="$builddir" "$in" "$out"
    description = coverage to html of go package in $in

rule goaddmain
    command = seb -tool asset -out "$out" main.go
    description = setup go main package $out

rule symlink
    # Need to rm $out first, otherwise might try to create $out/$(basename $out)
    # GNU has -T flag but this is more compatible.
    command = rm -f $out && ln -sf $target $out
    description = symlinking $out

rule protocc
    command = protoc --proto_path=$$(dirname $in) --cpp_out=$$(dirname $out) $in

rule protoc-c
    command = protoc-c --proto_path=$$(dirname $in) --c_out=$$(dirname $out) $in

rule download
    command = ( curl -s -L -o $out $$(head -1 $in) && test "$$(sha256sum $out)" = "$$(tail -1 $in)  $out" ) || ( rm -f $out ; exit 1 )
    description = Downloading url from $in

rule unzip
    command = unzip -p $in $file > $out || ( rm -f $out ; exit 1 )
    description= Unzip $file from $in

rule bunzip2
    command = rm -rf $out ; bunzip2 --keep --quiet --stdout $in > $out
depend_rule_bunzip2=/dev/null

rule untar
    command = tar -x -O -f $in $file > $out || ( rm -f $out ; exit 1 )
    description= Untar $file from $in

rule ronn
    command = seb -tool ronn $in $out
    description = Generate man page using ronn

# The following rules depend on variables that are defined in
# toolrules.ninja. The reason for this is that they in turn
# depend on paths defined by per-flavor ninja files.

rule in
    command = seb -tool in "$inconf" $in $out
    restat=1

rule inconfig
    command = seb -tool invars -buildvars="$buildvars" -invars="$builtin_invars" -out "$out" $in
    depfile = $out.d
`
