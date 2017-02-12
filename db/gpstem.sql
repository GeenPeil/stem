-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.0-alpha1
-- PostgreSQL version: 9.6
-- Project Site: pgmodeler.com.br
-- Model Author: Geert-Johan Riemer

SET check_function_bodies = false;
-- ddl-end --

-- object: rutte | type: ROLE --
-- DROP ROLE IF EXISTS rutte;
CREATE ROLE rutte WITH 
	SUPERUSER
	LOGIN
	UNENCRYPTED PASSWORD 'rutte';
-- ddl-end --


-- Database creation must be done outside an multicommand file.
-- These commands were put in this file only for convenience.
-- -- object: gpstem | type: DATABASE --
-- -- DROP DATABASE IF EXISTS gpstem;
-- CREATE DATABASE gpstem
-- ;
-- -- ddl-end --
-- 

-- object: members | type: SCHEMA --
-- DROP SCHEMA IF EXISTS members CASCADE;
CREATE SCHEMA members;
-- ddl-end --
ALTER SCHEMA members OWNER TO postgres;
-- ddl-end --

-- object: i8n | type: SCHEMA --
-- DROP SCHEMA IF EXISTS i8n CASCADE;
CREATE SCHEMA i8n;
-- ddl-end --
ALTER SCHEMA i8n OWNER TO postgres;
-- ddl-end --

SET search_path TO pg_catalog,public,members,i8n;
-- ddl-end --

-- object: members.accounts | type: TABLE --
-- DROP TABLE IF EXISTS members.accounts CASCADE;
CREATE TABLE members.accounts(
	id serial NOT NULL,
	email varchar(200) NOT NULL,
	loginname varchar(200),
	nickname varchar(30),
	given_name varchar(100) NOT NULL,
	first_names varchar(200) NOT NULL DEFAULT '',
	initials varchar(20) NOT NULL DEFAULT '',
	last_name_prefix varchar(30),
	last_name varchar(200) NOT NULL,
	last_name_full varchar(230),
	birthdate date NOT NULL,
	phonenumber varchar(15) NOT NULL DEFAULT '',
	postalcode varchar(10) NOT NULL DEFAULT '',
	housenumber varchar(10) NOT NULL DEFAULT '',
	housenumber_suffix varchar(10) NOT NULL DEFAULT '',
	streetname varchar(200) NOT NULL DEFAULT '',
	city varchar(200) NOT NULL DEFAULT '',
	province varchar(100) NOT NULL DEFAULT '',
	country_id integer,
	fee_last_payment_date date,
	verified_email boolean NOT NULL DEFAULT false,
	verified_identity boolean NOT NULL DEFAULT false,
	verified_voting_entitlement boolean NOT NULL DEFAULT false,
	textsearch_vector tsvector,
	registration_token varchar(200) NOT NULL,
	registration_date date NOT NULL DEFAULT now(),
	CONSTRAINT accounts_pk PRIMARY KEY (id),
	CONSTRAINT accounts_uq_loginname UNIQUE (loginname),
	CONSTRAINT accounts_uq_nickname UNIQUE (nickname),
	CONSTRAINT accounts_check_nickname_length CHECK (char_length(nickname) > 3),
	CONSTRAINT accounts_check_age_over_14 CHECK (date_part('year', age(birthdate)) >= 14)

);
-- ddl-end --
COMMENT ON COLUMN members.accounts.nickname IS 'When NULL, user has never chosen a nickname.';
-- ddl-end --
COMMENT ON COLUMN members.accounts.first_names IS 'Separated by space';
-- ddl-end --
COMMENT ON COLUMN members.accounts.initials IS 'Initials, no separation signs';
-- ddl-end --
COMMENT ON COLUMN members.accounts.last_name IS 'Includes prefix';
-- ddl-end --
COMMENT ON COLUMN members.accounts.fee_last_payment_date IS 'When set to NULL, member has never paid.';
-- ddl-end --
COMMENT ON CONSTRAINT accounts_check_age_over_14 ON members.accounts  IS 'One can only become a member at age 14 or older.';
-- ddl-end --
ALTER TABLE members.accounts OWNER TO postgres;
-- ddl-end --

-- Appended SQL commands --
CREATE OR REPLACE FUNCTION is_adult(members.accounts) RETURNS boolean AS $$
	SELECT date_part('year', age($1.birthdate)) >= 18
$$ VOLATILE LANGUAGE SQL;

CREATE OR REPLACE FUNCTION fee_paid(members.accounts) RETURNS boolean AS $$
	SELECT coalesce(date_part('year', age($1.fee_last_payment_date)) < 1, false)
$$ VOLATILE LANGUAGE SQL;
-- ddl-end --

-- object: members.accounts_fn_lock_identity | type: FUNCTION --
-- DROP FUNCTION IF EXISTS members.accounts_fn_lock_identity() CASCADE;
CREATE FUNCTION members.accounts_fn_lock_identity ()
	RETURNS trigger
	LANGUAGE plpgsql
	VOLATILE 
	CALLED ON NULL INPUT
	SECURITY INVOKER
	COST 1
	AS $$
DECLARE
	item_dimension_id integer;
	image_dimension_id integer;
BEGIN
	IF OLD.verified_identity AND NOT NEW.verified_identity THEN 
		RAISE EXCEPTION 'Cannot undo verified identity';
	END IF;
	IF OLD.verified_identity THEN
		IF OLD.initials != NEW.initials THEN
			RAISE EXCEPTION 'Cannot change initials when identity is verified';
		END IF;
		IF OLD.last_name != NEW.last_name THEN
			RAISE EXCEPTION 'Cannot change last_name when identity is verified';
		END IF;
		IF OLD.birthdate != NEW.birthdate THEN
			RAISE EXCEPTION 'Cannot change birthdate when identity is verified';
		END IF;
	END IF;
	RETURN NEW;
END;
$$;
-- ddl-end --
ALTER FUNCTION members.accounts_fn_lock_identity() OWNER TO postgres;
-- ddl-end --

-- object: accounts_tr_lock_identity | type: TRIGGER --
-- accounts_tr_lock_identity ON members.accounts CASCADE;
CREATE CONSTRAINT TRIGGER accounts_tr_lock_identity
	AFTER UPDATE
	ON members.accounts
	NOT DEFERRABLE 
	FOR EACH ROW
	EXECUTE PROCEDURE members.accounts_fn_lock_identity();
-- ddl-end --

-- object: members.sessions | type: TABLE --
-- DROP TABLE IF EXISTS members.sessions CASCADE;
CREATE TABLE members.sessions(
	id serial NOT NULL,
	account_id integer NOT NULL,
	token varchar(200) NOT NULL,
	timeout timestamp NOT NULL DEFAULT NOW(),
	CONSTRAINT sessions_pk PRIMARY KEY (id)

);
-- ddl-end --
ALTER TABLE members.sessions OWNER TO postgres;
-- ddl-end --

-- object: members.entitlement_proofs | type: TABLE --
-- DROP TABLE IF EXISTS members.entitlement_proofs CASCADE;
CREATE TABLE members.entitlement_proofs(
	id serial NOT NULL,
	account_id integer NOT NULL,
	filename varchar(100) NOT NULL,
	original_filename varchar(200) NOT NULL,
	CONSTRAINT entitlement_proofs_pk PRIMARY KEY (id),
	CONSTRAINT entitlement_proofs_uq_filename UNIQUE (filename)

);
-- ddl-end --
ALTER TABLE members.entitlement_proofs OWNER TO postgres;
-- ddl-end --

-- object: members.accounts_fn_update_textsearch | type: FUNCTION --
-- DROP FUNCTION IF EXISTS members.accounts_fn_update_textsearch() CASCADE;
CREATE FUNCTION members.accounts_fn_update_textsearch ()
	RETURNS trigger
	LANGUAGE plpgsql
	VOLATILE 
	CALLED ON NULL INPUT
	SECURITY INVOKER
	COST 1
	AS $$
BEGIN
	NEW.textsearch_vector = to_tsvector('english', 
		coalesce(NEW.given_name,'') || ' ' ||
		coalesce(NEW.first_names,'') || ' ' ||
		coalesce(NEW.initials,'') || ' ' ||
		coalesce(NEW.last_name,'') || ' ' ||
		coalesce(NEW.first_names,'') || ' ' ||
		coalesce(NEW.streetname, '') || ' ' ||
		coalesce(NEW.city, ''));
RETURN NEW;
END;

$$;
-- ddl-end --
ALTER FUNCTION members.accounts_fn_update_textsearch() OWNER TO postgres;
-- ddl-end --

-- object: accounts_tr_update_textsearch | type: TRIGGER --
-- DROP TRIGGER IF EXISTS accounts_tr_update_textsearch ON members.accounts CASCADE;
CREATE TRIGGER accounts_tr_update_textsearch
	BEFORE INSERT OR UPDATE
	ON members.accounts
	FOR EACH ROW
	EXECUTE PROCEDURE members.accounts_fn_update_textsearch();
-- ddl-end --

-- object: accounts_index_textsearch_vector | type: INDEX --
-- DROP INDEX IF EXISTS members.accounts_index_textsearch_vector CASCADE;
CREATE INDEX accounts_index_textsearch_vector ON members.accounts
	USING gin
	(
	  textsearch_vector
	);
-- ddl-end --

-- object: i8n.countries | type: TABLE --
-- DROP TABLE IF EXISTS i8n.countries CASCADE;
CREATE TABLE i8n.countries(
	id serial NOT NULL,
	code char(2) NOT NULL,
	CONSTRAINT countries_pk PRIMARY KEY (id),
	CONSTRAINT countries_uq_code UNIQUE (code)

);
-- ddl-end --
ALTER TABLE i8n.countries OWNER TO postgres;
-- ddl-end --

-- object: i8n.languages | type: TABLE --
-- DROP TABLE IF EXISTS i8n.languages CASCADE;
CREATE TABLE i8n.languages(
	id serial NOT NULL,
	code char(5) NOT NULL,
	name varchar(64) NOT NULL,
	CONSTRAINT languages_pk PRIMARY KEY (id),
	CONSTRAINT languages_uq_code UNIQUE (code)

);
-- ddl-end --
ALTER TABLE i8n.languages OWNER TO postgres;
-- ddl-end --

-- object: i8n.country_names | type: TABLE --
-- DROP TABLE IF EXISTS i8n.country_names CASCADE;
CREATE TABLE i8n.country_names(
	id serial NOT NULL,
	country_id integer NOT NULL,
	language_id integer NOT NULL,
	name varchar(64) NOT NULL,
	CONSTRAINT country_names_pk PRIMARY KEY (id),
	CONSTRAINT country_names_uq_country_language UNIQUE (country_id,language_id)

);
-- ddl-end --
ALTER TABLE i8n.country_names OWNER TO postgres;
-- ddl-end --

-- object: members.enum_mollie_status | type: TYPE --
-- DROP TYPE IF EXISTS members.enum_mollie_status CASCADE;
CREATE TYPE members.enum_mollie_status AS
 ENUM ('open','cancelled','expired','failed','pending','paid','paidout','refunded','charged_back');
-- ddl-end --
ALTER TYPE members.enum_mollie_status OWNER TO postgres;
-- ddl-end --

-- object: members.email_verifications | type: TABLE --
-- DROP TABLE IF EXISTS members.email_verifications CASCADE;
CREATE TABLE members.email_verifications(
	id serial NOT NULL,
	account_id integer NOT NULL,
	email varchar(200) NOT NULL,
	proof varchar(200) NOT NULL,
	created date NOT NULL DEFAULT now(),
	verified date,
	CONSTRAINT email_verifications_pk PRIMARY KEY (id)

);
-- ddl-end --
ALTER TABLE members.email_verifications OWNER TO postgres;
-- ddl-end --

-- object: email_verifications_idx_email_proof | type: INDEX --
-- DROP INDEX IF EXISTS members.email_verifications_idx_email_proof CASCADE;
CREATE UNIQUE INDEX email_verifications_idx_email_proof ON members.email_verifications
	USING btree
	(
	  email,
	  proof ASC NULLS LAST
	);
-- ddl-end --

-- object: members.payments | type: TABLE --
-- DROP TABLE IF EXISTS members.payments CASCADE;
CREATE TABLE members.payments(
	id serial NOT NULL,
	account_id integer NOT NULL,
	token varchar(80) NOT NULL,
	fee_amount money NOT NULL DEFAULT 12.00,
	donation_amount money NOT NULL DEFAULT 0.00,
	purchase_amount money NOT NULL DEFAULT 0.00,
	purchase_options varchar(200)[] NOT NULL DEFAULT '{}',
	mollie_id varchar(20) NOT NULL,
	mollie_status members.enum_mollie_status NOT NULL,
	mollie_created timestamp NOT NULL,
	mollie_paid timestamp,
	created timestamp NOT NULL DEFAULT now(),
	CONSTRAINT payments_pk PRIMARY KEY (id),
	CONSTRAINT payments_uq_mollie_id UNIQUE (mollie_id),
	CONSTRAINT payment_check_fee_amount CHECK (fee_amount = 12.00::money),
	CONSTRAINT payments_uq_token UNIQUE (token)

);
-- ddl-end --
ALTER TABLE members.payments OWNER TO postgres;
-- ddl-end --

-- Appended SQL commands --
CREATE OR REPLACE FUNCTION total_amount(members.payments) RETURNS money AS $$
	SELECT $1.fee_amount + $1.donation_amount + $1.purchase_amount
$$ VOLATILE LANGUAGE SQL;
-- ddl-end --

-- object: members.payments_fn_update_account_payment_date | type: FUNCTION --
-- DROP FUNCTION IF EXISTS members.payments_fn_update_account_payment_date() CASCADE;
CREATE FUNCTION members.payments_fn_update_account_payment_date ()
	RETURNS trigger
	LANGUAGE plpgsql
	VOLATILE 
	CALLED ON NULL INPUT
	SECURITY INVOKER
	COST 1
	AS $$
BEGIN
	IF NEW.mollie_status = 'paid'::members.enum_mollie_status AND NEW.mollie_status != OLD.mollie_status THEN
		UPDATE members.accounts SET fee_last_payment_date = NEW.mollie_created WHERE id = NEW.account_id;
	END IF;
	RETURN NEW;
END;
$$;
-- ddl-end --
ALTER FUNCTION members.payments_fn_update_account_payment_date() OWNER TO postgres;
-- ddl-end --

-- object: payments_tr_update_account_payment_date | type: TRIGGER --
-- DROP TRIGGER IF EXISTS payments_tr_update_account_payment_date ON members.payments CASCADE;
CREATE TRIGGER payments_tr_update_account_payment_date
	AFTER UPDATE
	ON members.payments
	FOR EACH ROW
	EXECUTE PROCEDURE members.payments_fn_update_account_payment_date();
-- ddl-end --

-- object: members.accounts_fn_set_last_name_full | type: FUNCTION --
-- DROP FUNCTION IF EXISTS members.accounts_fn_set_last_name_full() CASCADE;
CREATE FUNCTION members.accounts_fn_set_last_name_full ()
	RETURNS trigger
	LANGUAGE plpgsql
	VOLATILE 
	CALLED ON NULL INPUT
	SECURITY INVOKER
	COST 1
	AS $$
BEGIN
	NEW.last_name_full = COALESCE(NEW.last_name_prefix, '') + NEW.last_name;
	RETURN NEW;
END;
$$;
-- ddl-end --
ALTER FUNCTION members.accounts_fn_set_last_name_full() OWNER TO postgres;
-- ddl-end --

-- object: accounts_fk_country | type: CONSTRAINT --
-- ALTER TABLE members.accounts DROP CONSTRAINT IF EXISTS accounts_fk_country CASCADE;
ALTER TABLE members.accounts ADD CONSTRAINT accounts_fk_country FOREIGN KEY (country_id)
REFERENCES i8n.countries (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: sessions_fk_account | type: CONSTRAINT --
-- ALTER TABLE members.sessions DROP CONSTRAINT IF EXISTS sessions_fk_account CASCADE;
ALTER TABLE members.sessions ADD CONSTRAINT sessions_fk_account FOREIGN KEY (account_id)
REFERENCES members.accounts (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: entitlement_proofs_fk_account | type: CONSTRAINT --
-- ALTER TABLE members.entitlement_proofs DROP CONSTRAINT IF EXISTS entitlement_proofs_fk_account CASCADE;
ALTER TABLE members.entitlement_proofs ADD CONSTRAINT entitlement_proofs_fk_account FOREIGN KEY (account_id)
REFERENCES members.accounts (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: country_names_fk_country | type: CONSTRAINT --
-- ALTER TABLE i8n.country_names DROP CONSTRAINT IF EXISTS country_names_fk_country CASCADE;
ALTER TABLE i8n.country_names ADD CONSTRAINT country_names_fk_country FOREIGN KEY (country_id)
REFERENCES i8n.countries (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE RESTRICT;
-- ddl-end --

-- object: country_names_fk_language | type: CONSTRAINT --
-- ALTER TABLE i8n.country_names DROP CONSTRAINT IF EXISTS country_names_fk_language CASCADE;
ALTER TABLE i8n.country_names ADD CONSTRAINT country_names_fk_language FOREIGN KEY (language_id)
REFERENCES i8n.languages (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE RESTRICT;
-- ddl-end --

-- object: email_verifications_fk_account | type: CONSTRAINT --
-- ALTER TABLE members.email_verifications DROP CONSTRAINT IF EXISTS email_verifications_fk_account CASCADE;
ALTER TABLE members.email_verifications ADD CONSTRAINT email_verifications_fk_account FOREIGN KEY (account_id)
REFERENCES members.accounts (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: payments_fk_account | type: CONSTRAINT --
-- ALTER TABLE members.payments DROP CONSTRAINT IF EXISTS payments_fk_account CASCADE;
ALTER TABLE members.payments ADD CONSTRAINT payments_fk_account FOREIGN KEY (account_id)
REFERENCES members.accounts (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --


