from alpine
run apk update && apk add ca-certificates
add ./rewrite-proxy-go /rewrite-proxy-go
cmd ["./rewrite-proxy-go"]
