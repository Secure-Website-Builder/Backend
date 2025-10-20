-- name: CreateStoreOwner :one
INSERT INTO store_owner (name, email, password_hash)
VALUES ($1, $2, $3)
RETURNING store_owner_id, name, email, password_hash, created_at;

-- name: GetStoreOwnerByEmail :one
SELECT store_owner_id, name, email, password_hash, created_at
FROM store_owner
WHERE email = $1;

-- name: CreateStore :one
INSERT INTO store (store_owner_id, name, domain)
VALUES ($1, $2, $3)
RETURNING store_id, store_owner_id, name, domain, created_at, updated_at;

-- name: GetStoreWithOwner :one
SELECT s.store_id, s.name AS store_name, s.domain, s.created_at, s.updated_at,
       o.store_owner_id, o.name AS owner_name, o.email AS owner_email
FROM store s
JOIN store_owner o ON s.store_owner_id = o.store_owner_id
WHERE s.store_id = $1;
