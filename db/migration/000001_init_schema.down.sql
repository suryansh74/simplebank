BEGIN;

DROP TABLE entries;
DROP TABLE transfers;
DROP TABLE accounts;

DROP TYPE IF EXISTS "Currency";

COMMIT;
