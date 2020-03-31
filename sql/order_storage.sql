CREATE TABLE IF NOT EXISTS public.order_storage (
     id bigserial NOT NULL CONSTRAINT order_storage_pkey PRIMARY KEY,
     content character varying(2) COLLATE pg_catalog."default" NOT NULL,
     active bool NOT NULL DEFAULT TRUE,
     counter bigint NOT NULL DEFAULT 0,
     UNIQUE(content)
);