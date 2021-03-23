-- +migrate Up

-- +migrate StatementBegin
SET search_path TO 'nx_worker';


CREATE TYPE job_status AS ENUM ('ONPROGRESS', 'ONPROGRESS-ERROR', 'OK', 'ERROR');

CREATE SEQUENCE IF NOT EXISTS audit_system_pkey_seq;
CREATE TABLE "audit_system"
(
    id             BIGINT NOT NULL             DEFAULT nextval('audit_system_pkey_seq'::regclass),
    table_name     VARCHAR(256),
    uuid_key       uuid,
    primary_key    BIGINT,
    data           TEXT,
    action         INTEGER,
	description	   VARCHAR(256),
    created_by     BIGINT,
    created_client VARCHAR(256),
    created_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted        BOOLEAN DEFAULT FALSE,
    CONSTRAINT pk_auditsystem_id PRIMARY KEY (id)
);

CREATE INDEX idx_auditsystem_tablename_primarykey ON audit_system (table_name, primary_key);
CREATE INDEX idx_auditsystem_createdat_tablename ON audit_system (created_at, table_name);
CREATE INDEX idx_auditsystem_createdby_createdat ON audit_system (created_by, table_name);
CREATE INDEX idx_auditsystem_createdclient_createdat ON audit_system (created_at, table_name);

CREATE SEQUENCE IF NOT EXISTS job_process_pkey_seq;
CREATE TABLE "job_process"
(
    id                 BIGINT NOT NULL             DEFAULT nextval('job_process_pkey_seq'::regclass),
    uuid_key           uuid                        DEFAULT public.uuid_generate_v4(),
    parent_job_id      VARCHAR(256),
    level              INTEGER,
    job_id             VARCHAR(256),
    "group"            VARCHAR(256),
    type               VARCHAR(256),
    name               VARCHAR(256),
    counter            INTEGER,
    total              INTEGER,
    status             job_status                  DEFAULT 'ONPROGRESS',
    message_alert      VARCHAR(256),
    alert_id           VARCHAR(256),
    parameter          TEXT,
    alert_content_data TEXT,
    url                VARCHAR(256),
    filename           VARCHAR(256),
    created_by         BIGINT,
    created_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_client     VARCHAR(256),
    updated_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted            BOOLEAN                     DEFAULT FALSE,
    CONSTRAINT pk_jobprocess_id PRIMARY KEY (id),
    CONSTRAINT uq_jobprocess_jobid UNIQUE (job_id)
);



-- +migrate StatementEnd
