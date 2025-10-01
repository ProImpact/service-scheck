
-- name: ServiceGetAll :many
SELECT * FROM "services";

-- name: ServiceGetByName :one
SELECT * FROM "services" WHERE "service_name" = ?;

-- name: ServiceDeleteByName :exec
DELETE FROM "services" WHERE "service_name" = ?;

-- name: ServiceCreate :exec
INSERT INTO "services" ("id","service_name","command_type","cmd_check_command","startup_time","status","execute_command","healtcheck_endpoint","ping_time","pid")
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ServiceFullUpdate :exec
UPDATE services SET 
    "service_name" = ?,
    "status" = ?,
    "execute_command" = ?,
    "healtcheck_endpoint" = ?,
    "ping_time" = ?,
    "pid" = ?,
    "command_type" = ?,
    "cmd_check_command" = ?
WHERE "service_name" = ?;

-- name: ServiceUpdateStatus :exec
UPDATE services SET 
    "status" = ?
WHERE "service_name" = ?;


-- name: ServiceChangeStatus :exec
UPDATE services SET "status" = ? WHERE "service_name" = ?;

