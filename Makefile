dc:
	@docker compose up --build --force-recreate --watch

k8s:
	@skaffold dev --no-prune=false

doc:
	@swag init -g cmd/main.go -o docs -ot json

test:
	@go test tests/* --json | jq -r .Output | grep -v null