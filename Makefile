
sync-run:
	cd sync/service && PLANNER_DB_HOST=localhost PLANNER_DB_PORT=5432 PLANNER_DB_NAME=planner PLANNER_DB_USER=test PLANNER_DB_PASSWORD=test PLANNER_PORT=8092 PLANNER_API_KEY=testKey go run .

database:
	docker run -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=planner -p 5432:5432 postgres:16
