# nodeinfo

[![GoDoc](https://godoc.org/github.com/m-lab/nodeinfo?status.svg)](https://godoc.org/github.com/m-lab/nodeinfo) [![Build Status](https://travis-ci.org/m-lab/nodeinfo.svg?branch=master)](https://travis-ci.org/m-lab/nodeinfo) [![Coverage Status](https://coveralls.io/repos/github/m-lab/nodeinfo/badge.svg?branch=master)](https://coveralls.io/github/m-lab/nodeinfo?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/nodeinfo)](https://goreportcard.com/report/github.com/m-lab/nodeinfo)

This collects data about the hardware, software, and configs on nodes on the
[M-Lab](https://www.measurementlab.net) platform.  Every hour (in expectation,
with some randomization) the output of `lspci`, `lshw`, `ifconfig`, and others
is written to disk. This allows us to track the configuration of fleet nodes
over time.

Available as a container in
[measurementlab/nodeinfo](https://hub.docker.com/r/measurementlab/nodeinfo/) on
Docker Hub.

# design

As simple as possible. This system is called `nodeinfo`. Every command produces its own type of data, and so is it own datatype.  These two facts, together with [M-Lab's unified naming scheme for data](http://example.com), and the best practices for [Pusher](http://github.com/m-lab/pusher) mean that the directory structure for output is fully determined. 

This program calls a series of other programs, and directs the output of each call to the appropriate output file. The set of programs to call is currently hard-coded in the binary. If any of the commands run unsuccessfully, this crashes.  Every command is rerun every hour on average, with some randomness. The inter-run times are drawn from the exponential distribution to try and make sure the resulting series of measurements has the [PASTA property](https://en.wikipedia.org/wiki/Arrival_theorem#Theorem_for_arrivals_governed_by_a_Poisson_process).
