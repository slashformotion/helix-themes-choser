all:
	go run image-builder/main.go
	cp -r out/* public/
	go run website-builder/main.go
