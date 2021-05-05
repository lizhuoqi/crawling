# note: call scripts from /scripts
# A phony target is one that is not really the name of a file; rather it is just a name for a recipe to be executed when you make an explicit request. There are two reasons to use a phony target: to avoid a conflict with a file of the same name, and to improve performance.
# If you write a rule whose recipe will not create the target file, the recipe will be executed every time the target comes up for remaking. 

.PHONY: build
build:
# go build  -o . main.go
	go mod tidy
	go build -v

run:
	./crawl

dev: build run
