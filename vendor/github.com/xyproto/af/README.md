# af [![Build Status](https://travis-ci.org/xyproto/af.svg?branch=master)](https://travis-ci.org/xyproto/af) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/af)](https://goreportcard.com/report/github.com/xyproto/af) [![GoDoc](https://godoc.org/github.com/xyproto/af?status.svg)](https://godoc.org/github.com/xyproto/af)

Activation functions for neural networks.

These activation functions are included:

* Swish (`x / (1 + exp(-x))`)
* Sigmoid (`1 / (1 + exp(-x))`)
* SoftPlus (`log(1 + exp(x))`)
* Gaussian01 (`exp(-(x * x) / 2.0)`)
* Sin (`math.Sin(math.Pi * x)`)
* Cos (`math.Cos(math.Pi * x)`)
* Linear (`x`)
* Inv (`-x`)
* ReLU (`x >= 0 ? x : 0`)
* Squared (`x * x`)

These `math` functions are included just for convenience:

* Abs (`math.Abs`)
* Tanh (`math.Tanh`)

One functions that takes two arguments is also included:

* PReLU (`x >= 0 ? x : x * a`)

## Requirements

* Go 1.11 or later.

## General information

* License: MIT
* Version: 0.3.2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
