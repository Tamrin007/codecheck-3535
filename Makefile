setup:
	which godep || go get github.com/tools/godep
	godep restore

save:
	godep save

build:
	go build server.go

run:
	make build && ./server

clean:
	[ -e ./server ] && rm server

deploy:
	sudo service nginx restart
	killall server || true
	git pull origin master
	make clean
	go build server.go
	nohup ./server &
