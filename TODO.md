# TODO

- [ ] Train and test with the Mnist dataset.
- [ ] Fix any remaining issues with drawing SVG diagrams.
- [ ] Fix any remaining issues with generating expressions.
- [ ] Draw an "O" on the output node in the diagram.
- [ ] Fix an issue with mutating the Network structs (the Neurons needs to be mutated too. And there are pointers everywhere to them).
- [ ] Store all neurons within the network, but keep pointers to input nodes (and the output node might work).
      Then all those neurons can be modified when the Network mutates, but the pointers can be kept the same.
