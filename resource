### mongodb ###
https://www.mongodb.com/docs/drivers/go/current/quick-start/

go get go.mongodb.org/mongo-driver/mongo

### gofiber ###
go get github.com/gofiber/fiber/v2

### docker ###
docker run --name mongodb -d mongo:lastest -p 27017:27017


### command ###
go mod init github.com/jayzedx/hotel-reservation
make build
make run

./bin/api --listenAddr :7000
