.PHONY: proto tools

# Update all protocol buffers
proto:
	make -C proto/

# Build all the tools
alltools: 
	make -C tools/

# Build all the nodes
allnodes:
	make -C node/

allclients:
	make -C client/

testfiles: 
	make -C tools/ testfiles

clean:
	@make -C proto/ clean
	@make -C node/ clean
	@make -C tools/e2c clean
	@make -C client/ clean