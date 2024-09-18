
sync-run:
	cd sync/service && PLANNER_DB_PATH=test.db PLANNER_PORT=8092 PLANNER_API_KEY=testKey go run .

