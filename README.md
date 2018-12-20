# nodeinfo

[![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/nodeinfo)](https://goreportcard.com/report/github.com/m-lab/nodeinfo)

This collects data about the hardware, software, and configs on nodes on the
[M-Lab](https://www.measurementlab.net) platform.  Every hour (in expectation,
with some randomization) the output of `lspci`, `lshw`, `ifconfig`, and others
is written to disk. This allows us to track the configuration of fleet nodes
over time.

Available as a container in
[measurementlab/nodeinfo](https://hub.docker.com/r/measurementlab/nodeinfo/) on
Docker Hub.
