.PHONY: all testfiles

all: bin
	@echo "Finished building tools"

testfiles: bin test3 test10 b10 b20 b5 d2-f1 d2-f2 d5-f1 d5-f2
	@echo "Built test files"

bin: genConfig.go 
	go build -o genConfig $^

test3: bin
	./genConfig -o ../../testData/apollo/3-node-test/ -n3

test10: bin
	./genConfig -o ../../testData/apollo/10-node-test/ -n10

b5: bin
	./genConfig -o ../../testData/apollo/b5/ -n3 -b5 -d1000

b10: bin
	./genConfig -o ../../testData/apollo/b10/ -n3 -b10 -d1000

b20: bin
	./genConfig -o ../../testData/apollo/b20/ -n3 -b20 -d1000

d2-f1: bin
	./genConfig -o ../../testData/apollo/d2-f1/ -n3 -f1 -b20 -d2000

d2-f2: bin
	./genConfig -o ../../testData/apollo/d2-f2/ -n5 -f2 -b20 -d2000

d5-f1: bin
	./genConfig -o ../../testData/apollo/d5-f1/ -n3 -b20 -d5000

d5-f2: bin
	./genConfig -o ../../testData/apollo/d5-f2/ -n5 -f2 -b20 -d5000

clean: 
	rm -rf genConfig