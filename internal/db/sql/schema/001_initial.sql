-- +goose Up
-- +goose StatementBegin
CREATE TABLE "services" (
    "id" TEXT PRIMARY KEY,
    "service_name" TEXT NOT NULL UNIQUE,
    "startup_time" TIMESTAMP NOT NULL DEFAULT NOW,
    "status" TEXT NOT NULL,
    "command_type" TEXT NOT NULL,
    "cmd_check_command" TEXT,
    "execute_command" TEXT NOT NULL,
    "healtcheck_endpoint" TEXT NOT NULL,
    "ping_time" TEXT NOT NULL,
    "pid" INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  "services" ;
-- +goose StatementEnd