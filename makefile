all:
    go run images-builder/main.go
    cp -r out/* public/
    go run website-builder/main.go
