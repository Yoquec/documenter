MAIN = src/main.go
SHARE_DIR = /usr/share/go-documenter

build:
	go build -o bin/go-documenter ${MAIN}

run:
	go run ${MAIN}

install: bin/go-documenter
	# move assets
	rm -rf ${SHARE_DIR}
	mkdir -p ${SHARE_DIR}
	cp -r static/* ${SHARE_DIR}/
	chmod 0655 ${SHARE_DIR}/*
	# move binary
	rm -f /usr/bin/go-documenter
	install -C -Dm 755 -t /usr/bin bin/go-documenter
