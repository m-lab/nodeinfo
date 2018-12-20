# metadata-collector

[![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/metadata-collector)](https://goreportcard.com/report/github.com/m-lab/metadata-collector)

This collects metadata for experiments on the
[M-Lab](https://www.measurementlab.net) platform.  Every hour (in expectation,
with some randomization) the output of `lspci`, `lshw`, `ifconfig`, and others
is written to disk. This allows us to track the configuration of fleet nodes
over time.

Available as a container in [measurementlab/metadata-collector](https://hub.docker.com/r/measurementlabmetadata-collector/) on Docker Hub.
