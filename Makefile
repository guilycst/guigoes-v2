HUGO ?= hugo

install:
	@echo "Install Hugo Extended first: https://gohugo.io/installation/"

run:
	$(HUGO) server -D --disableFastRender

build:
	$(HUGO) --minify

rundocker:
	docker build -t guigoes-hugo:latest .
	docker run --platform linux/amd64 -p 8080:80 --rm -it guigoes-hugo:latest