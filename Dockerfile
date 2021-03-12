FROM golang:1.14

# WORKDIR /root/projects/src/github.com/239103/webupload
# COPY . .

RUN go get -d -v github.com/239103/webupload
RUN go install -v github.com/239103/webupload

CMD ["webupload"]