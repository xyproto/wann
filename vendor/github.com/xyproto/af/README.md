# af [![Build Status](https://travis-ci.org/xyproto/af.svg?branch=master)](https://travis-ci.org/xyproto/af) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/af)](https://goreportcard.com/report/github.com/xyproto/af) [![GoDoc](https://godoc.org/github.com/xyproto/af?status.svg)](https://godoc.org/github.com/xyproto/af)

Activation functions, intended for use in neural networks.

Provides the following activation functions, that take just one argument:

* Sigmoid (optimized, from the [swish](https://github.com/xyproto/swish) package).
* Swish (optimized, from the [swish](https://github.com/xyproto/swish) package).
* SoftPlus (optimized, from the [swish](https://github.com/xyproto/swish) package).
* Abs (`math.Abs`)
* Tanh (`math.Tanh`)
* Sin (`math.Sin`)
* Cos (`math.Cos`)
* Inv (`-x`)
* ReLU (`x >= 0 ? x : 0`)

The `math` functions are included just for convenience.

And also these functions, that take two arguments:

* PReLU (`x >= 0 ? x : x * a`)


## Requirements

* Require Go 1.11 or later.

## General information

* License: MIT
* Version: 0.3.0
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
