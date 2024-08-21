

all:

	@echo "idontkare"



.PHONY: test
test:

	go build -o test.out ./test/

	./test.out

clean:

	rm -rf *.out