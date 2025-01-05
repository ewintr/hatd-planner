package sqlite

var migrations = []string{
	`CREATE TABLE events ("id" TEXT UNIQUE, "title" TEXT, "start" TIMESTAMP, "duration" TEXT)`,
	`PRAGMA journal_mode=WAL`,
	`PRAGMA synchronous=NORMAL`,
	`PRAGMA cache_size=2000`,
	`CREATE TABLE localids ("id" TEXT UNIQUE, "local_id" INTEGER)`,
	`CREATE TABLE items (
    id TEXT PRIMARY KEY NOT NULL,
    kind TEXT NOT NULL,
    updated TIMESTAMP NOT NULL,
    deleted BOOLEAN NOT NULL,
    body TEXT NOT NULL
)`,
	`ALTER TABLE events ADD COLUMN recur_period TEXT`,
	`ALTER TABLE events ADD COLUMN recur_count INTEGER`,
	`ALTER TABLE events ADD COLUMN recur_start TIMESTAMP`,
	`ALTER TABLE events ADD COLUMN recur_next TIMESTAMP`,
	`ALTER TABLE events DROP COLUMN recur_period`,
	`ALTER TABLE events DROP COLUMN recur_count`,
	`ALTER TABLE events DROP COLUMN recur_start`,
	`ALTER TABLE events DROP COLUMN recur_next`,
	`ALTER TABLE events ADD COLUMN recur TEXT`,
	`ALTER TABLE items ADD COLUMN recurrer TEXT`,
	`ALTER TABLE events DROP COLUMN recur`,
	`ALTER TABLE events ADD COLUMN recurrer TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE events ADD COLUMN recur_next TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE events DROP COLUMN start`,
	`ALTER TABLE events ADD COLUMN date TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE events ADD COLUMN time TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE items ADD COLUMN recur_next TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE events RENAME TO tasks`,

	// add unique constraint to localids.local_id
	`CREATE TABLE localids_backup AS SELECT * FROM localids`,
	`DROP TABLE localids`,
	`CREATE TABLE localids ("id" TEXT UNIQUE, "local_id" INTEGER UNIQUE)`,
	`INSERT INTO localids (id, local_id)
    SELECT id, local_id FROM localids_backup`,
	`DROP TABLE localids_backup`,

	`ALTER TABLE items ADD COLUMN date TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE tasks ADD COLUMN project TEXT NOT NULL DEFAULT ''`,
	`CREATE TABLE syncupdate ("timestamp" TIMESTAMP NOT NULL)`,
	`INSERT INTO syncupdate (timestamp) VALUES ("0001-01-01T00:00:00Z")`,
}
