.PHONY: all testfiles

all: bin
	@echo "Finished building tools"

bin:
	make -C e2c/ bin
	make -C synchs/ bin
	make -C apollo/ bin

testfiles:
	make -C e2c/ testfiles
	make -C synchs/ testfiles
	make -C apollo/ testfiles

clean: 
	make -C e2c/ clean
	make -C synchs/ clean
	make -C apollo/ clean