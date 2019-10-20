#!/bin/sh

echo 'downloading training set images'
curl -O 'http://yann.lecun.com/exdb/mnist/train-images-idx3-ubyte.gz'

echo 'downloading training set labels'
curl -O 'http://yann.lecun.com/exdb/mnist/train-labels-idx1-ubyte.gz'

echo 'downloading test set images'
curl -O 'http://yann.lecun.com/exdb/mnist/t10k-images-idx3-ubyte.gz'

echo 'downloading test set labels'
curl -O 'http://yann.lecun.com/exdb/mnist/t10k-labels-idx1-ubyte.gz'

echo 'extracting'
gunzip -v *.gz
