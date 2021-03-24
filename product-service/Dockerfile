FROM golang:1.16

# Download and install any required third party dependencies into the container.
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/gorilla/mux
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/urfave/negroni
RUN go get github.com/juang77/GoMSArch/shared/util
RUN go get github.com/juang77/GoMSArch/shared/models
RUN go get github.com/juang77/GoMSArch/shared/helper
RUN go get github.com/joho/godotenv
RUN go get github.com/sirupsen/logrus
RUN go get github.com/meatballhat/negroni-logrus

# 
ADD . /go/src/github.com/juang77/LaikaGoMSV1/product-service
WORKDIR /go/src/github.com/juang77/LaikaGoMSV1/product-service
RUN go build main.go

# Expose port 5005 to the host so we can access our application
EXPOSE 5005

# Tell Docker what command to run when the container starts
CMD ["/go/src/github.com/juang77/LaikaGoMSV1/product-service/main"]