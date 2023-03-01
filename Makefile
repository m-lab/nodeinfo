#
# The main purpose of this Makefile is to help local development and testing.
#
SOURCE_FILES=api/node_info.go config/config.go data/gather.go main.go metrics/metrics.go
CONFIG=./testdata/config.json
DATADIR=./testdata
DATATYPE=nodeinfo1

run: nodeinfo
	rm -rf $(DATADIR)/$(DATATYPE)
	./nodeinfo -config $(CONFIG) -datadir $(DATADIR) -once -smoketest -wait 1s; echo; tree $(DATADIR); echo

nodeinfo: $(SOURCE_FILES)
	go build -race .

check:
	@find $(DATADIR)/$(DATATYPE) -name '*.json' -exec echo -e '\n>>>' {} \; -a -exec jq . {} \; | \
		sed -e 's/./.../80' -e 's/\.\.\..*/.../'
