# skill-typing-back
## ローカル環境構築
1. dockerコンテナを起動
```
make up
```
2. mysqlのuser作成(初回起動時のみ)
```
make db_create_user
```
3. mysqlのDB初期化 (DBを再度作成したい場合も使用可)
```
make db_init
```
4. migrationコマンド実行
```
make db_migrate
```
5. go 開発用サーバーを起動
```
make dev 
```
6. dockerコンテナを削除
```
make down 
```

DBにログインして、mysqlクライアントで操作をしたい場合は下記のコマンドが使えます。
```
make db_login
```

DBを一旦全消去して、再度migrationをやりなおしたいときは、下記の２つを実行してください。
```
make db_init
make db_migrate
```
<!--
dockerコンテナを削除した場合、ターミナルやコマンドプロンプトでのコマンド操作が必要です。
```
$ dockedb_initr exec -it dev-mysql mysql -u root -p
```           
でdockerのmysqlに接続
db_init
```
CREATE USER 'myuser'@'%' IDENTIFIED WITH mysql_native_password BY 'root';
GRANT ALL PRIVILEGES ON *.* TO 'myuser'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
CREATE DATABASE mydatabase;
GRANT ALL PRIVILEGES ON mydatabase.* TO 'myuser'@'%';
FLUSH PRIVILEGES;
```  
'myuser'と'root'とmydatabaseは環境変数で指定したもの。

```  
EXIT
```

db/db.goを変更したとき
```
$ GO_ENV=dev go run migrate/migrate.go
```
-->
