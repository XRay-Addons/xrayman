-- +goose Up
-- +goose StatementBegin
ALTER TABLE total_nodes_traffic
    RENAME TO nodes_stats;
ALTER TABLE nodes_stats
	ADD COLUMN open_connections INTEGER NOT NULL DEFAULT 0,
	ADD COLUMN cpu_load REAL NOT NULL DEFAULT 0.0,
	ADD COLUMN ram_load REAL NOT NULL DEFAULT 0.0,
    ADD COLUMN mem_load REAL NOT NULL DEFAULT 0.0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nodes_stats
	DROP COLUMN mem_load,
	DROP COLUMN ram_load,
	DROP COLUMN cpu_load,
  	DROP COLUMN open_connections;
ALTER TABLE nodes_stats
    RENAME TO total_nodes_traffic;
-- +goose StatementEnd
