.PHONY: install build run check rundocker

install:
	sh scripts/install-margo.sh

build: install
	python3 scripts/build.py

check: build
	python3 -m unittest discover -s tests -v
	python3 scripts/check-site.py

run: build
	python3 -m http.server 1313 --bind 127.0.0.1 --directory public

rundocker:
	docker build -t guigoes-margo:latest .
	docker run -p 8080:80 --rm -it guigoes-margo:latest
