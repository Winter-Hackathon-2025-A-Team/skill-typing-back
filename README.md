# skill-typing-back
## ローカル環境構築
1. dockerコンテナを起動
```
$ make up
```

2. go 開発用サーバーを起動
```
$ make dev 
```

3. dockerコンテナを削除
```
$ make down 
```

db/db.goを変更したとき
```
$ GO_ENV=dev go run migrate/migrate.go
```

dockerコンテナを削除した場合、ターミナルやコマンドプロンプトでのコマンド操作が必要です。
```
$ docker exec -it dev-mysql mysql -u root -p
```           
でdockerのmysqlに接続

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
