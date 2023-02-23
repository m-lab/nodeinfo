SOURCE_FILES=config/config.go data/gather.go main.go metrics/metrics.go config-old.json config-new.json
OLD=old/var/spool/nodeinfo
NEW=new/var/spool/nodeinfo
OLD_DATATYPES=biosversion chassisserial ipaddress iproute4 iproute6 lshw lspci lsusb osrelease uname
NEW_DATATYPE=nodeinfo1

run: nodeinfo
	@rm -rf old new
	./nodeinfo -config ./config-old.json -datadir $(OLD) -once -smoketest -wait 1s; echo; tree $(OLD); echo
	./nodeinfo -config ./config-new.json -datadir $(NEW) -once -smoketest -wait 1s; echo; tree $(NEW); echo

nodeinfo: $(SOURCE_FILES)
	go build .

check:
	@for i in $(OLD_DATATYPES); do \
		find $(OLD)/$$i -name '*.txt' -exec echo -e '\n>>>' {} \; -a -exec head -10 {} \;; \
	done
	@find $(NEW)/$(NEW_DATATYPE) -name '*.json' -exec echo -e '\n>>>' {} \; -a -exec jq . {} \; | \
		sed -e 's/./.../80' -e 's/\.\.\..*/.../'
