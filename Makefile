.PHONY: proto clean

# Update all protocol buffers
proto:
	make -C proto/

clean:
	@make -C node/ clean
	@make -C tools/e2c clean
	@make -C client/ clean