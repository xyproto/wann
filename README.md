## wann

[![Build Status](https://travis-ci.org/xyproto/wann.svg?branch=master)](https://travis-ci.org/xyproto/wann) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/wann)](https://goreportcard.com/report/github.com/xyproto/wann) [![GoDoc](https://godoc.org/github.com/xyproto/wann?status.svg)](https://godoc.org/github.com/xyproto/wann)

Weight Agnostic Neural Networks in Go.

Inspired by: https://ai.googleblog.com/2019/08/exploring-weight-agnostic-neural.html

## Features and limitations

* Neural networks can be trained, and they work, but I have only tried this on very simple training data and there is surely a lot of room for improvement, both in term of benchmarking/profiling and controlling the rate of mutation.
* Currently, a random weight is chosen when training, instead of looping over the range of the weight. The paper describes the latter, but it's unclear to me what the ideal step size for the weight is, when looping.
* Complex networks are given a worse score when evolving. A quick benchmark at the start of the program determines which activation function us more complex. This optimizes not only for simple networks, but also for network performance. 

## Quick start

This might require Go 1.12 or later.

Clone the repository:

    git clone https://github.com/xyproto/wann

Enter the `cmd/evolve` directory:

    cd wann/cmd/evolve

Build and run the example:

    go build && ./evolve

Take a look at the best network for judging if a set of numbers that are either 0 or 1 are of one category:

    xdg-open best.svg

(If needed, use your favorite SVG viewer intead of the `xdg-open` command).

## Ideas

* Adding convolution nodes might give interesting results.

## Diagrams

There is included functionality for drawing networks and saving them as SVG files. This functionality needs more testing, but it can output a variety of diagrams.

The activation functions are plotted directly onto the nodes.

Here are a few examples:

<img alt=diagram src=img/diagram.svg width=128 />

<img alt=diagram src=img/best.svg width=128 />

<img alt=diagram src=img/test.svg width=128 />

<img alt=diagram src=img/before.svg width=128 />

This one is from a debugging session and is included just for fun:

<img alt=diagram src=img/wip.svg />

## General info

* Version: 0.0.1
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
