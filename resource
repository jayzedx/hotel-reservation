### mongodb ###
https://www.mongodb.com/docs/drivers/go/current/quick-start/

go get go.mongodb.org/mongo-driver/mongo

### gofiber ###
go get github.com/gofiber/fiber/v2

### go get ###
go get golang.org/x/crypto/bcrypt
go get github.com/xdg-go/stringprep@v1.0.4

### docker ###
docker run --name mongodb -d mongo:lastest -p 27017:27017


### command ###
go mod init github.com/jayzedx/hotel-reservation
make build
make run

./bin/api --listenAddr :7000

### extension ###
thunder client

### jwt ###
go get -u github.com/golang-jwt/jwt/v5


### export environment in ubuntu ###
## nano ~/.bashrc
export JWT_SECRET=adsdsadasdasdasdasd22132das
## source ~/.bashrc // Run the following command to apply the changes
echo $JWT_SECRET //You can verify that the environment 


### configure ###
go get github.com/spf13/viper

### log ###
go get -u go.uber.org/zap

