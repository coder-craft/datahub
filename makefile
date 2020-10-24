#Go parameters

build:
	rm -rf target/
	mkdir target/
	cp server/datahub/config.json datahub/datahub.cfg
	go build -o target/datahub server/datahub/main.go
clean:
	rm -rf target
run:
	nohup target/datahub --conf=target/datahub.cfg > target/datahub.out &

stop:
	pkill -f target/datahub
