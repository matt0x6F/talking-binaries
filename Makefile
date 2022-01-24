build:
	go build -o dist/program main.go
	go build -o dist/plugin-echo ./plugins/echo
	go build -o dist/plugin-shell ./plugins/shell