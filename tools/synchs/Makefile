.PHONY: all testfiles

all: bin
	@echo "Finished building tools"

testfiles: bin test3 test10
	@echo "Built test files"

bin: genConfig.go 
	go build -o genConfig $^

test3: bin
	./genConfig -o ../../testData/synchs/3-node-test/ -n3

test10: bin
	./genConfig -o ../../testData/synchs/10-node-test/ -n10

clean: 
	rm -rf genConfig