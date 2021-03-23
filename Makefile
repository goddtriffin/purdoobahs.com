$(VERBOSE).SILENT:
.DEFAULT_GOAL := help

.PHONY: help
help: # displays Makefile target info
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%-30s %s\n" "target" "help" ; \
	printf "%-30s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
			IFS=$$':' ; \
			help_split=($$help_line) ; \
			help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			printf '\033[36m'; \
			printf "%-30s %s" $$help_command ; \
			printf '\033[0m'; \
			printf "%s\n" $$help_info; \
	done

.PHONY: buildstatic
buildstatic:  ## builds static resources
	# generate/clean bin
	mkdir -p bin
	rm -rf bin/*
	mkdir -p bin/static
	rm -rf bin/static/*

	# copy over html, files, images
	cp -Rn ui/html bin/html
	cp -Rn ui/static/file bin/static/file
	cp -Rn ui/static/image bin/static/image

.PHONY: buildlocal
buildlocal: buildstatic  ## builds the binary locally
	# TLS files only necessary for local development
	cp -Rn tls bin/tls
	go build -o bin/website ./...

.PHONY: runlocal
runlocal: buildlocal ## runs the binary locally
	./bin/website -env=development

.PHONY: builddocker
builddocker: ## builds the binary and Docker container
	docker build --tag purdoobahs-com --file build/Dockerfile .

.PHONY: rundocker
rundocker: builddocker ## creates and runs a new Docker container
	docker run \
	--name "purdoobahs_com" \
	-d --restart unless-stopped \
	-p 8080:80 \
	purdoobahs-com

.PHONY: startdocker
startdocker: ## resumes a stopped Docker container
	docker start purdoobahs_com

.PHONY: stopdocker
stopdocker: ## stops the Docker container
	docker stop purdoobahs_com

.PHONY: removedocker
removedocker: ## removes the Docker container
	docker rm purdoobahs_com

.PHONY: memusage
memusage: ## displays the memory usage of the currently running Docker container
	docker stats purdoobahs_com --no-stream --format "{{.Container}}: {{.MemUsage}}"

.PHONY: logs
logs: ## displays logs from the currently running Docker container
	docker logs purdoobahs_com

.PHONY: certs
certs: ## generates a cert.pem and key.pem for https://localhost development
	mkdir -p tls
	rm -rf tls/*
	go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
	mv -nv cert.pem tls
	mv -nv key.pem tls
