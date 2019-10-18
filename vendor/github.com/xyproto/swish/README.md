# Swish

[![Build Status](https://travis-ci.org/xyproto/swish.svg?branch=master)](https://travis-ci.org/xyproto/swish) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/swish)](https://goreportcard.com/report/github.com/xyproto/swish) [![GoDoc](https://godoc.org/github.com/xyproto/swish?status.svg)](https://godoc.org/github.com/xyproto/swish)

An optimized Swish activation function ([Ramachandran, Zoph and Le, 2017](https://arxiv.org/abs/1710.05941)), for neural networks.

## Screenshots

![](img/swish.png)

![](img/sigmoid.png)

The graphs above were drawn using the program in `cmd/graph`, which uses [goterm](https://github.com/buger/goterm).

## Benchmark Results

### Using a `Swish` function that uses `math.Exp`

First run:

```
goos: linux
goarch: amd64
pkg: github.com/xyproto/swish
BenchmarkSwish07-8   	200000000	         8.93 ns/op
BenchmarkSwish03-8   	200000000	         8.95 ns/op
PASS
ok  	github.com/xyproto/swish	5.391s
```

### Using the optimized `Swish` function that uses `exp256`

```
goos: linux
goarch: amd64
pkg: github.com/xyproto/swish
BenchmarkSwish07-8   	2000000000	         0.26 ns/op
BenchmarkSwish03-8   	2000000000	         0.26 ns/op
PASS
ok  	github.com/xyproto/swish	1.108s
```

The optimized `Swish` function is **34x** faster than the one that uses `math.Exp`, and quite a bit faster than my (apparently bad) attempt at a hand-written assembly version.

The average error (difference in output value) between the optimized and non-optimized version is `+-0.0013` and the maximum error is `+-0.0024`. This is for `x` in the range `[5,3]`. See the program in `cmd/precision` for how this was calculated.

```
0.00015
0.00001
goos: linux
goarch: amd64
pkg: github.com/xyproto/swish
BenchmarkSwishAssembly07-8      500000000                3.63 ns/op
BenchmarkSwishAssembly03-8      500000000                3.65 ns/op
BenchmarkSwish07-8              2000000000               0.30 ns/op
BenchmarkSwish03-8              2000000000               0.26 ns/op
BenchmarkSwishPrecise07-8       200000000                9.07 ns/op
BenchmarkSwishPrecise03-8       200000000                9.25 ns/op
PASS
ok      github.com/xyproto/swish        11.100s
```

I have no idea why the assembly version is so slow, but `0.26 ns/op` isn't bad for a non-hand-optimized version.

## General info

* Version: 1.3.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
