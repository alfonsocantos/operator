.PHONY: deps
deps:
	go mod tidy -v;
	go mod download;
	go mod vendor;


.PHONY: test
test: start _test stop

.PHONY: _test
_test:
	INTEGRATION=true go test ./... -cover

.PHONY: buildserver
buildserver:
	docker build -t my_custom_sshd .

.PHONY: start
start: buildserver
	docker run --name my_custom_sshd --rm -d -p 2222:22 my_custom_sshd

.PHONY: stop
stop:
	docker stop my_custom_sshd