docker build -t kochbuch .
docker container rm kochbuch-dev
docker run -p 80:80 -e DB_Server=172.17.0.2 -e DB_User=kbdev -e DB_Password=foo -e DB_Name=kochbuch --name kochbuch-dev kochbuch