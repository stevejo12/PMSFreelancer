FROM golang:1.12.0-alpine3.9
LABEL maintainer="SPIRITS <spirits.project.thesis@gmail.com>"
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk update && apk add git && apk add openssh

ADD id_rsa /root/.ssh/id_rsa
ADD id_rsa.pub /root/.ssh/id_rsa.pub


RUN chmod 0700 /root/.ssh/id_rsa && \
    chmod 0600 /root/.ssh/id_rsa.pub

RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN ssh-keyscan 13.229.188.59 >> /root/.ssh/known_hosts

RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

RUN go get github.com/stevejo12/PMSFreelancer
RUN go get github.com/gin-gonic/gin
RUN go get github.com/swaggo/files
RUN go get github.com/swaggo/gin-swagger
RUN go get github.com/alecthomas/template
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/adlio/trello
RUN go get github.com/dgrijalva/jwt-go
RUN go get golang.org/x/oauth2
RUN go get cloud.google.com/go
RUN go get gopkg.in/gomail.v2
RUN go get golang.org/x/crypto/bcrypt

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

CMD ["/app/main"]