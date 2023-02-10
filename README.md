# go serverles api

env GOOS=linux go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o build/main cmd/main.go

zip -jrm build/main.zip build/main
