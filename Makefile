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