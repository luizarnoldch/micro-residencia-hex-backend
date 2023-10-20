TEMPLATE_FILE := templates/main.yml
STACK_NAME := residentes


init:
	go mod init main
update:
	go mod tidy
build:
	./scripts/build.sh
test:
	go test ./tests/...
f_test:
	./scripts/func_test.sh
mock:
	mockery --all --output ./tests/mocks/
sam:
	sam build --template-file $(TEMPLATE_FILE)
deploy:
	sam deploy --template-file $(TEMPLATE_FILE) --stack-name $(STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
b-deploy:
	make build
	make deploy
	make f_test
e2e:
	make destroy
	sleep 7
	make test
	make build
	make deploy
	sleep 3
	make f_test


