.PHONY: all testfiles

all: bin
	@echo "Finished building tools"

bin:
	make -C e2c/ bin
	make -C synchs/ bin

testfiles:
	make -C e2c/ testfiles
	make -C synchs testfiles

clean: 
	make -C e2c/ clean
	make -C synchs clean