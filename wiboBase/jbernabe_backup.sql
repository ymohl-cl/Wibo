--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: Group; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE "Group" (
    id integer NOT NULL,
    groupname character varying(20) NOT NULL
);


ALTER TABLE "Group" OWNER TO jbernabe;

--
-- Name: User; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE "User" (
    id integer NOT NULL,
    login character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    salt character varying(255) NOT NULL,
    lastlogin date NOT NULL,
    creationdate date NOT NULL,
    mail character varying(510) NOT NULL,
    groupid integer NOT NULL,
    usemagnet date NOT NULL,
    device_id integer NOT NULL
);


ALTER TABLE "User" OWNER TO jbernabe;

--
-- Name: checkpoints; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE checkpoints (
    id integer NOT NULL,
    longitude numeric(5,2) NOT NULL,
    latitude numeric(5,2) NOT NULL,
    date date NOT NULL,
    containerid integer NOT NULL,
    attractbymagnet boolean NOT NULL
);


ALTER TABLE checkpoints OWNER TO jbernabe;

--
-- Name: container; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE container (
    id integer NOT NULL,
    longitude numeric(5,2) NOT NULL,
    latitude numeric(5,2) NOT NULL,
    direction numeric(5,2) NOT NULL,
    speed integer NOT NULL,
    creationdate date NOT NULL,
    userid integer NOT NULL,
    typecontainerid integer NOT NULL
);


ALTER TABLE container OWNER TO jbernabe;

--
-- Name: container_type_information; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE container_type_information (
    id_type integer NOT NULL,
    name_type character varying(255) NOT NULL
);


ALTER TABLE container_type_information OWNER TO jbernabe;

--
-- Name: device; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE device (
    id integer NOT NULL,
    macaddr character varying(18) NOT NULL,
    typedeviceid integer NOT NULL
);


ALTER TABLE device OWNER TO jbernabe;

--
-- Name: followed; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE followed (
    id integer NOT NULL,
    user_id integer NOT NULL,
    container_id integer NOT NULL
);


ALTER TABLE followed OWNER TO jbernabe;

--
-- Name: message; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE message (
    id integer NOT NULL,
    content text NOT NULL,
    containerid integer NOT NULL,
    userid integer NOT NULL,
    typemessageid integer NOT NULL
);


ALTER TABLE message OWNER TO jbernabe;

--
-- Name: reception; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE reception (
    id integer NOT NULL,
    receptiontime date NOT NULL,
    longitude numeric(5,2) NOT NULL,
    latitude numeric(5,2) NOT NULL,
    userid integer NOT NULL,
    idcontainer integer NOT NULL
);


ALTER TABLE reception OWNER TO jbernabe;

--
-- Name: session; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE session (
    id integer NOT NULL,
    logged boolean NOT NULL,
    logintime date NOT NULL,
    logouttime boolean NOT NULL,
    userid integer
);


ALTER TABLE session OWNER TO jbernabe;

--
-- Name: shared; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE shared (
    id integer NOT NULL,
    type_shared integer NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE shared OWNER TO jbernabe;

--
-- Name: type_container; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE type_container (
    id integer NOT NULL,
    typename character varying(255) NOT NULL
);


ALTER TABLE type_container OWNER TO jbernabe;

--
-- Name: type_device; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE type_device (
    id integer NOT NULL,
    typename character varying(255) NOT NULL
);


ALTER TABLE type_device OWNER TO jbernabe;

--
-- Name: type_message; Type: TABLE; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE TABLE type_message (
    id integer NOT NULL,
    typename character varying(255) NOT NULL
);


ALTER TABLE type_message OWNER TO jbernabe;

--
-- Data for Name: Group; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY "Group" (id, groupname) FROM stdin;
\.


--
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY "User" (id, login, password, salt, lastlogin, creationdate, mail, groupid, usemagnet, device_id) FROM stdin;
\.


--
-- Data for Name: checkpoints; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY checkpoints (id, longitude, latitude, date, containerid, attractbymagnet) FROM stdin;
\.


--
-- Data for Name: container; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY container (id, longitude, latitude, direction, speed, creationdate, userid, typecontainerid) FROM stdin;
\.


--
-- Data for Name: container_type_information; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY container_type_information (id_type, name_type) FROM stdin;
\.


--
-- Data for Name: device; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY device (id, macaddr, typedeviceid) FROM stdin;
\.


--
-- Data for Name: followed; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY followed (id, user_id, container_id) FROM stdin;
\.


--
-- Data for Name: message; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY message (id, content, containerid, userid, typemessageid) FROM stdin;
\.


--
-- Data for Name: reception; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY reception (id, receptiontime, longitude, latitude, userid, idcontainer) FROM stdin;
\.


--
-- Data for Name: session; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY session (id, logged, logintime, logouttime, userid) FROM stdin;
\.


--
-- Data for Name: shared; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY shared (id, type_shared, user_id) FROM stdin;
\.


--
-- Data for Name: type_container; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY type_container (id, typename) FROM stdin;
\.


--
-- Data for Name: type_device; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY type_device (id, typename) FROM stdin;
\.


--
-- Data for Name: type_message; Type: TABLE DATA; Schema: public; Owner: jbernabe
--

COPY type_message (id, typename) FROM stdin;
\.


--
-- Name: checkpoints_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY checkpoints
    ADD CONSTRAINT checkpoints_pk PRIMARY KEY (id);


--
-- Name: container_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY container
    ADD CONSTRAINT container_pk PRIMARY KEY (id);


--
-- Name: device_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY device
    ADD CONSTRAINT device_pk PRIMARY KEY (id);


--
-- Name: followed_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY followed
    ADD CONSTRAINT followed_pk PRIMARY KEY (id);


--
-- Name: group_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY "Group"
    ADD CONSTRAINT group_pk PRIMARY KEY (id);


--
-- Name: message_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY message
    ADD CONSTRAINT message_pk PRIMARY KEY (id);


--
-- Name: reception_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY reception
    ADD CONSTRAINT reception_pk PRIMARY KEY (id);


--
-- Name: session_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_pk PRIMARY KEY (id);


--
-- Name: shared_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY shared
    ADD CONSTRAINT shared_pk PRIMARY KEY (id);


--
-- Name: type_container_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY type_container
    ADD CONSTRAINT type_container_pk PRIMARY KEY (id);


--
-- Name: type_container_type_information_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY container_type_information
    ADD CONSTRAINT type_container_type_information_pk PRIMARY KEY (id_type);


--
-- Name: type_device_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY type_device
    ADD CONSTRAINT type_device_pk PRIMARY KEY (id);


--
-- Name: type_message_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY type_message
    ADD CONSTRAINT type_message_pk PRIMARY KEY (id);


--
-- Name: user_pk; Type: CONSTRAINT; Schema: public; Owner: jbernabe; Tablespace: 
--

ALTER TABLE ONLY "User"
    ADD CONSTRAINT user_pk PRIMARY KEY (id);


--
-- Name: checkpoints_container_id_idx; Type: INDEX; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE INDEX checkpoints_container_id_idx ON checkpoints USING btree (containerid, date);


--
-- Name: message_idx_container; Type: INDEX; Schema: public; Owner: jbernabe; Tablespace: 
--

CREATE INDEX message_idx_container ON message USING btree (containerid);


--
-- Name: Session_User; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY session
    ADD CONSTRAINT "Session_User" FOREIGN KEY (userid) REFERENCES "User"(id);


--
-- Name: checkpoints_container; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY checkpoints
    ADD CONSTRAINT checkpoints_container FOREIGN KEY (containerid) REFERENCES container(id);


--
-- Name: container_type_container; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY container
    ADD CONSTRAINT container_type_container FOREIGN KEY (typecontainerid) REFERENCES type_container(id);


--
-- Name: container_user; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY container
    ADD CONSTRAINT container_user FOREIGN KEY (userid) REFERENCES "User"(id);


--
-- Name: device_type_device; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY device
    ADD CONSTRAINT device_type_device FOREIGN KEY (typedeviceid) REFERENCES type_device(id);


--
-- Name: followed_container; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY followed
    ADD CONSTRAINT followed_container FOREIGN KEY (container_id) REFERENCES container(id);


--
-- Name: followed_user; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY followed
    ADD CONSTRAINT followed_user FOREIGN KEY (user_id) REFERENCES "User"(id);


--
-- Name: message_container; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY message
    ADD CONSTRAINT message_container FOREIGN KEY (containerid) REFERENCES container(id);


--
-- Name: message_type_message; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY message
    ADD CONSTRAINT message_type_message FOREIGN KEY (typemessageid) REFERENCES type_message(id);


--
-- Name: message_user; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY message
    ADD CONSTRAINT message_user FOREIGN KEY (userid) REFERENCES "User"(id);


--
-- Name: reception_user; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY reception
    ADD CONSTRAINT reception_user FOREIGN KEY (userid) REFERENCES "User"(id);


--
-- Name: shared_user; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY shared
    ADD CONSTRAINT shared_user FOREIGN KEY (user_id) REFERENCES "User"(id);


--
-- Name: user_device; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY "User"
    ADD CONSTRAINT user_device FOREIGN KEY (device_id) REFERENCES device(id);


--
-- Name: user_group; Type: FK CONSTRAINT; Schema: public; Owner: jbernabe
--

ALTER TABLE ONLY "User"
    ADD CONSTRAINT user_group FOREIGN KEY (groupid) REFERENCES "Group"(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: jbernabe
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM jbernabe;
GRANT ALL ON SCHEMA public TO jbernabe;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

