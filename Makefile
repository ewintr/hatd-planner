plan-build:
	go build -o dist/plan ./plan

plan-deploy:
	cp dist/plan ~/bin

sync-run:
	cd sync/service && go run . -dbname localhost -dbport 5432 -dbname planner -dbuser test -dbpassword test -port 8092 -key testKey

sync-debug:
	cd sync/service && dlv debug . -- -dbname localhost -dbport 5432 -dbname planner -dbuser test -dbpassword test -port 8092 -key testKey

sync-build:
	go build -o dist/plannersync ./sync/service/

sync-deploy:
	ssh server sudo /usr/bin/systemctl stop plannersync.service
	scp dist/plannersync server:/usr/local/bin/plannersync
	ssh server sudo /usr/bin/systemctl start plannersync.service

database:
	docker run -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=planner -p 5432:5432 postgres:16


