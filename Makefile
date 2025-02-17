up:
	docker-compose up -d
dev: 
	docker-compose exec app air
down: 
	docker-compose down 
db_create_user:
	docker-compose exec dev-mysql mysql -uroot -p -e'CREATE USER `myuser`@`%` IDENTIFIED WITH mysql_native_password BY "root"; GRANT ALL PRIVILEGES ON *.* TO `myuser`@`%` WITH GRANT OPTION; FLUSH PRIVILEGES;'
db_init:
	docker-compose exec dev-mysql mysql -uroot -p -e'DROP DATABASE IF EXISTS `mydatabase`; CREATE DATABASE mydatabase; GRANT ALL PRIVILEGES ON mydatabase.* TO `myuser`@`%`; FLUSH PRIVILEGES;'
db_migrate:
	docker-compose exec -e GO_ENV=dev app go run migrate/migrate.go
db_login:
	docker-compose exec dev-mysql mysql -uroot -p