plan-deploy:
	cd plan && go build -o plan . && mv plan ~/bin

sync-run:
	cd sync/service && go run . -dbname localhost -dbport 5432 -dbname planner -dbuser test -dbpassword test -port 8092 -key testKey

database:
	docker run -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=planner -p 5432:5432 postgres:16
