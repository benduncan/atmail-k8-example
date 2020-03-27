GO_PROJECT_NAME := atmail-rbl
K8_UPDATEDATE := $(shell date +'%s')
IMAGE_TAG := latest

# GO commands
go_build:
	@echo "\n....Building $(GO_PROJECT_NAME)"
	go build -o ./bin/$(GO_PROJECT_NAME) `ls -1 *.go`

go_dep_install:
	@echo "\n....Installing dependencies for $(GO_PROJECT_NAME)...."
	go get -v .

go_run:
	@echo "\n....Running $(GO_PROJECT_NAME)...."
	./bin/$(GO_PROJECT_NAME)

go_test:
	@echo "\n....Running tests for $(GO_PROJECT_NAME)...."
	go test

# Project rules
build:
	$(MAKE) go_dep_install
	$(MAKE) go_build

test:
	UNIT_TESTING=1 go test -v

run:
	$(MAKE) go_build
	$(MAKE) go_run

clean:
	rm -rf ./pkg/*
	rm -rf ./bin/*

docker:
	@echo "\n....Building docker image ($(IMAGE_TAG)) and uploading to GCR ...."
	#$(MAKE) test
	gcloud auth configure-docker
	docker build -t $(GO_PROJECT_NAME) .
	docker tag $(GO_PROJECT_NAME) gcr.io/$(PROJECT_ID)/$(GO_PROJECT_NAME):$(IMAGE_TAG)
	docker push gcr.io/$(PROJECT_ID)/$(GO_PROJECT_NAME):$(IMAGE_TAG)

kubernetes:
	@echo "\n....Updating Kubernetes ...."
	kubectl patch deployment/$(GO_PROJECT_NAME)-deployment -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"lastupdated\":\"$(K8_UPDATEDATE)\"}}}}}"

.PHONY: docker kubernetes db
