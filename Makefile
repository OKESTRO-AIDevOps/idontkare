
# SERVER_HOST := localhost
CLIENT_NAME := idk-client

all:

	@echo "idontkare"


build:

	rm -rf cert-server cert-client

	mkdir -p cert-server 

	mkdir -p cert-client

	go run ./hack/certgen $(SERVER_HOST) $(CLIENT_NAME)

	make -C src/server clean 
	
	make -C src/server server

	mv cert-server src/server

	/bin/cp -Rf apix src/server

	make -C src/client clean 

	make -C src/client client

	mv cert-client src/client

	/bin/cp -Rf apix src/client

	mkdir -p ../kiwi-web/cert-client

	cp -Rf src/client/cert-client/* ../kiwi-web/cert-client/

build-agent:

	make -C src/agent clean

	make -C src/agent agent

	/bin/cp -Rf apix src/agent

server:

	make -C src/server server 

client:
	
	make -C src/client client

agent:

	make -C src/agent agent


.PHONY: db
db:

	cd db && docker compose up --build -d

db-down:

	cd db && docker compose down

db-clean:

	cd db && rm -r data

.PHONY: test
test:

	go build -o test.out ./test/

	./test.out

clean:

	rm -rf *.out



	make -C src/server clean 
	make -C src/client clean 
	make -C src/agent clean 