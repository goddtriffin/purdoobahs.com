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

.PHONY: watch_sass
watch_sass: ## hot reloads Sass stylesheets
	sass --watch --update --style=compressed --no-source-map --color --unicode ui/static/scss:bin/static/stylesheet

.PHONY: gen_static
gen_static:  ## generates static resources
	# generate/clean bin
	rm -rf bin
	mkdir -p bin/static

	# copy over ui: html, files, images
	cp -R ui/html bin/html
	cp -R ui/static/file bin/static/file
	cp -R ui/static/image bin/static/image
	cp -R ui/static/video bin/static/video
	cp -R ui/static/script bin/static/script

	# generate css from sass
	sass --style=compressed --no-source-map --color --unicode ui/static/scss:bin/static/stylesheet

	# copy over assets
	cp -R assets bin/assets

.PHONY: build
build: ## builds the binary locally
	go build -o bin/website ./cmd/website

.PHONY: dev
dev: gen_static build ## runs the binary locally
	cd bin && ./website -env=development

.PHONY: build_docker
build_docker: ## builds the binary and Docker container
	docker build --tag goddtriffin/purdoobahs-website:latest --file deployment/Dockerfile .

.PHONY: run_docker
run_docker: build_docker ## creates and runs a new Docker container
	docker run \
	--name "purdoobahs-website" \
	-d --restart unless-stopped \
	-p 8080:8080 \
	goddtriffin/purdoobahs-website:latest

.PHONY: start_docker
start_docker: ## resumes a stopped Docker container
	docker start purdoobahs-website

.PHONY: stop_docker
stop_docker: ## stops the Docker container
	docker stop purdoobahs-website

.PHONY: remove_docker
remove_docker: ## removes the Docker container
	docker rm purdoobahs-website

.PHONY: push_docker
push_docker: ## pushes new Docker image to Docker Hub
	docker push goddtriffin/purdoobahs-website:latest

.PHONY: restart_deployment
restart_deployment: ## restarts all pods in the purdoobahs-website k8s deployment
	kubectl rollout restart deployment purdoobahs-website

.PHONY: deploy
deploy: build_docker push_docker restart_deployment # builds/pushes new docker image at :latest and restarts k8s deployment

.PHONY: mem_usage
mem_usage: ## displays the memory usage of the currently running Docker container
	docker stats purdoobahs-website --no-stream --format "{{.Container}}: {{.MemUsage}}"

.PHONY: logs
logs: ## displays logs from the currently running Docker container
	docker logs purdoobahs-website
