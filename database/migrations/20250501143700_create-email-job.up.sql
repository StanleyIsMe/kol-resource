CREATE TYPE "email_job_status" AS ENUM (
    'pending',
    'processing',
    'success',
    'partially_success',
    'failed',
    'canceled'
);

CREATE TYPE "email_log_status" AS ENUM (
    'pending',
    'success',
    'failed',
    'canceled'
);

CREATE TABLE "email_sender" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL,
    "key" varchar(255) NOT NULL,
    "domain" varchar(255) NOT NULL,
    "rate_limit" integer NOT NULL,
    "last_send_at" timestamp NOT NULL,
    "updated_admin_id" uuid NOT NULL,
    "enabled" boolean NOT NULL DEFAULT TRUE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (email)
);

CREATE TABLE "email_job" (
    "id" bigserial PRIMARY KEY,
    "expected_reciver_count" integer NOT NULL DEFAULT 0,
    "success_count" integer NOT NULL DEFAULT 0,
    "sender_id" uuid NOT NULL,
    "sender_name" varchar(255) NOT NULL,
    "sender_email" varchar(255) NOT NULL,
    "admin_id" uuid NOT NULL,
    "admin_name" varchar(255) NOT NULL,
    "product_id" uuid NOT NULL,
    "product_name" varchar(255) NOT NULL,
    "updated_admin_id" uuid NOT NULL,
    "memo" text NOT NULL,
    "payload" jsonb NOT NULL DEFAULT '{}' ::jsonb,
    "status" email_job_status NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_execute_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_status_created_at ON email_job(status,created_at);
CREATE INDEX IF NOT EXISTS idx_sender_id_created_at_status ON email_job(sender_id,created_at ASC) WHERE status IN ('pending', 'processing');

CREATE TABLE "email_log" (
    "id" bigserial PRIMARY KEY,
    "job_id" bigint NOT NULL,
    "product_id" uuid NOT NULL,
    "sender_id" uuid NOT NULL,
    "email" varchar(255) NOT NULL,
    "message_id" varchar(255) NOT NULL,
    "reply" boolean NOT NULL DEFAULT FALSE,
    "kol_id" uuid NOT NULL,
    "kol_name" varchar(255) NOT NULL,
    "status" email_log_status NOT NULL,
    "momo" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "sended_at" timestamp,
    UNIQUE (product_id,email)
);

CREATE INDEX IF NOT EXISTS idx_job_id_status ON email_log(job_id,status);
CREATE INDEX IF NOT EXISTS idx_sender_id_status_sended_at ON email_log(sender_id, status, sended_at) WHERE status = 'success';