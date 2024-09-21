# デフォルトの値を設定
SUITE ?=

# 条件分岐
ifeq ($(SUITE),)
  GINKGO_CMD := 	docker compose exec app bash -c 'cd bbs && ginkgo -v -r'
else
  GINKGO_CMD := 	docker compose exec app bash -c 'cd bbs && ginkgo -v -r -focus=$(SUITE)'
endif

build:
	docker compose up -d --build
down:
	docker compose down --rmi all -v
status:
	docker compose ps
app:
	docker compose exec app bash
db:
	docker compose exec db bash -c 'mysql -uroot -ppassword bbs-dev'
delete:
	docker compose down --rmi all --volumes --remove-orphans
test:
	@echo "Running Ginkgo with command: $(GINKGO_CMD)"
	@$(GINKGO_CMD)