#!/bin/sh
rm ./src/data/forum.db

go run main.go

sqlite3 ./src/data/forum.db <<EOF
insert into users (uuid,username,password,email,avatar,admin) values ('62,166,203,33,145,11,65,225,156,20,110,101,63,63,139,31','Texio','$2a$15$eX/qrGiqOg6VVshuvAAfiOe.XhOazdre/x54/tvAmYf9t2qRTJtgG',texio974@gmail.com,./uploads/avatars/default.png,1);
EOF