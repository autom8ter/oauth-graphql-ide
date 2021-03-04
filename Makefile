version := "0.0.3"

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Makefile Commands:"
	@echo "----------------------------------------------------------------"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo "----------------------------------------------------------------"

run:
	@go run main.go


patch: ## bump sem version by 1 patch
	bumpversion patch --allow-dirty

minor: ## bump sem version by 1 minor
	bumpversion minor --allow-dirty

tag: ## tag the repo (remember to commit changes beforehand)
	git tag v$(version)

push:
	git push origin v$(version)

docker-build:
	@docker build -t colemanword/oauth-graphql-playground:v$(version) .

docker-push:
	@docker push colemanword/oauth-graphql-playground:v$(version)


build: ## build the server to ./bin
	@mkdir -p bin
	@gox -osarch="linux/amd64" -output="bin/linux/{{.Dir}}_linux_amd64"
	@gox -osarch="darwin/amd64" -output="bin/darwin/{{.Dir}}_darwin_amd64"
	@gox -osarch="windows/amd64" -output="bin/windows/{{.Dir}}_windows_amd64"
