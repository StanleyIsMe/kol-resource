CREATE TYPE "sex" AS ENUM (
    'm',
    'f'
);

CREATE TABLE "kol" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar(50) NOT NULL,
    "email" varchar(50) NOT NULL,
    "description" text NOT NULL,
    "sex" sex NOT NULL,
    "enable" bool NOT NULL,
    "updated_admin_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    UNIQUE (email)
);

CREATE TABLE "admin" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar(50) NOT NULL,
    "username" varchar(50) NOT NULL,
    "password" varchar(50) NOT NULL,
    "salt" varchar(50) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    UNIQUE (username)
);

CREATE TABLE "tag" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar(50) NOT NULL,
    "updated_admin_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    UNIQUE (name)
);

CREATE TABLE "product" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "description" text NOT NULL,
    "updated_admin_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    UNIQUE (name)
);

CREATE TABLE "kol_tag" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "kol_id" uuid NOT NULL,
    "tag_id" uuid NOT NULL,
    "updated_admin_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    UNIQUE (kol_id, tag_id)
);

CREATE TABLE "send_email_log" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "kol_id" uuid NOT NULL,
    "kol_name" varchar(50) NOT NULL,
    "email" uuid NOT NULL,
    "admin_id" uuid NOT NULL,
    "admin_name" varchar(50) NOT NULL,
    "product_id" uuid NOT NULL, 
    "product_name" varchar(50) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);  
