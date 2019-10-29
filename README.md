## wann

[![Build Status](https://travis-ci.org/xyproto/wann.svg?branch=master)](https://travis-ci.org/xyproto/wann) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/wann)](https://goreportcard.com/report/github.com/xyproto/wann) [![GoDoc](https://godoc.org/github.com/xyproto/wann?status.svg)](https://godoc.org/github.com/xyproto/wann)

Weight Agnostic Neural Networks, implemented in Go.

Inspired by: https://ai.googleblog.com/2019/08/exploring-weight-agnostic-neural.html

## Features and limitations

* Neural networks can be trained, and they work, but I have only tried this on very simple training data and there is surely a lot of room for improvement, both in term of benchmarking/profiling and improving how the rate of mutation is controlled.
* Currently, a random weight is chosen when training, instead of looping over the range of the weight, as described in the paper. It's unclear to me what the ideal step size for the weight is, when looping.

## Quick start

This might require Go 1.12 or later.

Clone the repository:

    git clone https://github.com/xyproto/wann

Enter the `cmd/evolve` directory:

    cd wann/cmd/evolve

Build and run the example:

    go build && ./evolve

Take a look at the best image for the last generation:

    xdg-open best.svg

(If needed, use your favorite SVG viewer intead of the `xdg-open` command).

## Ideas

* Adding convolution nodes might give interesting results.

## Diagrams

There is included functionality for drawing networks and saving them as SVG files. This functionality needs more testing, but it can output a variety of diagrams.

The activation functions are plotted directly onto the nodes.

Here are a few examples:

<img alt=diagram src=img/diagram.svg width=128 />

<img alt=diagram src=img/test.svg width=128 />

<img alt=diagram src=img/before.svg width=128 />

<img alt=diagram src=img/after.svg width=128 />

This one happened during debugging, and is included just for fun:

<img alt=diagram src=img/wip.svg />

## General info

* Version: 0.0.1
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;