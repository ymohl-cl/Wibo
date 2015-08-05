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


--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry, geography, and raster spatial types and functions';


SET search_path = public, pg_catalog;

--
-- Name: getalldevices(); Type: FUNCTION; Schema: public; Owner: wibo
--

CREATE FUNCTION getalldevices() RETURNS TABLE(deviceid integer, userid integer, loginuser character varying, usertype integer, usermail character varying)
    LANGUAGE plpgsql
    AS $$ BEGIN RETURN QUERY EXECUTE 'SELECT device.id AS DeviceId, device.user_id_user AS UserId, "user".login AS LoginUser, "user".id_type_g AS UserType, "user".mail AS UserMail FROM device LEFT OUTER JOIN "user" ON (device.user_id_user = "user".id_user)'; END  $$;


ALTER FUNCTION public.getalldevices() OWNER TO wibo;

--
-- Name: getcontainers(); Type: FUNCTION; Schema: public; Owner: wibo
--

CREATE FUNCTION getcontainers() RETURNS TABLE(idballon integer, idtype integer, direction numeric, speedcont integer, creationdate date, deviceid integer, locationcont text)
    LANGUAGE plpgsql
    AS $$ 
 BEGIN  RETURN  QUERY EXECUTE 'SELECT  container.id AS contIndex, container.id_type_c AS TypeCode, container.direction AS contDirection, container.speed AS contSpeed, container.creationdate as dateCreationCon, container.device_id as idDevice,  ST_AsText(container.location_ct) FROM container';  END$$;


ALTER FUNCTION public.getcontainers() OWNER TO wibo;

--
-- Name: getcontainersbyuserid(integer); Type: FUNCTION; Schema: public; Owner: wibo
--

CREATE FUNCTION getcontainersbyuserid(iduser integer) RETURNS TABLE(idballon integer, titlename character varying, idtype integer, direction numeric, speedcont integer, creationdate timestamp with time zone, deviceid integer, locationcont text)
    LANGUAGE plpgsql
    AS $$  BEGIN RETURN QUERY SELECT container.id AS contIndex, container.titlename AS     TitleName, container.id_type_c AS TypeCode, container.direction AS contDirection, container.speed AS contSpeed, container.creationdate as dateCreationCon,     container.device_id as idDevice,  ST_AsText(container.location_ct)  FROM container  WHERE idcreator = iduser;  END $$;


ALTER FUNCTION public.getcontainersbyuserid(iduser integer) OWNER TO wibo;

--
-- Name: getcontainersbyuseridjson(integer); Type: FUNCTION; Schema: public; Owner: wibo
--

CREATE FUNCTION getcontainersbyuseridjson(iduser integer) RETURNS TABLE(idballon integer, titlename character varying, idtype integer, direction numeric, speedcont integer, creationdate date, deviceid integer, locationcont text)
    LANGUAGE plpgsql
    AS $$  BEGIN RETURN QUERY SELECT container.id AS contIndex, container.titlename AS TitleName, container.id_type_c AS TypeCode, container.direction AS contDirection, container.speed AS contSpeed, container.creationdate as dateCreationCon, container.device_id as idDevice,  ST_AsGeoJson(container.location_ct)  FROM container  WHERE idcreator = iduser;  END $$;


ALTER FUNCTION public.getcontainersbyuseridjson(iduser integer) OWNER TO wibo;

--
-- Name: getcontainersjson(); Type: FUNCTION; Schema: public; Owner: wibo
--

CREATE FUNCTION getcontainersjson() RETURNS TABLE(idballon integer, idtype integer, direction numeric, speedcont integer, creationdate date, deviceid integer, locationcont text)
    LANGUAGE plpgsql
    AS $$ 
 BEGIN  RETURN  QUERY EXECUTE 'SELECT  container.id AS contIndex, container.id_type_c AS TypeCode, container.direction AS contDirection, container.speed AS contSpeed, container.creationdate as dateCreationCon, container.device_id as idDevice,  ST_AsGeoJson(container.location_ct) FROM container';  END$$;


ALTER FUNCTION public.getcontainersjson() OWNER TO wibo;

--
-- Name: getdevicesbyuserid(integer); Type: FUNCTION; Schema: public; Owner: wibo
--

CREATE FUNCTION getdevicesbyuserid(iduser integer) RETURNS TABLE(macaddr character varying)
    LANGUAGE plpgsql
    AS $$  BEGIN RETURN QUERY SELECT  device.macaddr AS MacAddressDevice FROM device WHERE user_id_user = iduser; END  $$;


ALTER FUNCTION public.getdevicesbyuserid(iduser integer) OWNER TO wibo;

--
-- Name: insertcontainer(integer, integer, integer, integer, double precision, double precision, text, integer); Type: FUNCTION; Schema: public; Owner: wibo
--

CREATE FUNCTION insertcontainer(idcreatorc integer, latitudec integer, longitudec integer, device integer, directionc double precision, speedc double precision, title text, idx integer) RETURNS SETOF integer
    LANGUAGE plpgsql
    AS $$  BEGIN RETURN QUERY INSERT INTO container (direction, speed, device_id, location_ct, idcreator, titlename, ianix) VALUES(directionc, speedc , device, ST_SetSRID(ST_MakePoint(latitudec, longitudec), 4326), idcreatorc, title, idx) RETURNING id;  END; $$;


ALTER FUNCTION public.insertcontainer(idcreatorc integer, latitudec integer, longitudec integer, device integer, directionc double precision, speedc double precision, title text, idx integer) OWNER TO wibo;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: checkpoints; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE checkpoints (
    id integer NOT NULL,
    date date NOT NULL,
    containerid integer NOT NULL,
    attractbymagnet boolean NOT NULL,
    location_ckp geography(Point,4326)
);


ALTER TABLE checkpoints OWNER TO wibo;

--
-- Name: checkpoints_id_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE checkpoints_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE checkpoints_id_seq OWNER TO wibo;

--
-- Name: checkpoints_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE checkpoints_id_seq OWNED BY checkpoints.id;


--
-- Name: type_container; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE type_container (
    id_type_c integer NOT NULL,
    typename character varying(255) NOT NULL
);


ALTER TABLE type_container OWNER TO wibo;

--
-- Name: container; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE container (
    id integer NOT NULL,
    direction numeric(5,2) NOT NULL,
    speed integer NOT NULL,
    device_id integer NOT NULL,
    location_ct geography(Point,4326),
    idcreator integer,
    titlename character varying(255),
    ianix integer,
    creationdate timestamp with time zone DEFAULT now()
)
INHERITS (type_container);


ALTER TABLE container OWNER TO wibo;

--
-- Name: container_id_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE container_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE container_id_seq OWNER TO wibo;

--
-- Name: container_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE container_id_seq OWNED BY container.id;


--
-- Name: type_device; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE type_device (
    id_type_d integer NOT NULL,
    typename character varying(255) NOT NULL
);


ALTER TABLE type_device OWNER TO wibo;

--
-- Name: device; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE device (
    id integer NOT NULL,
    macaddr character varying(18) NOT NULL,
    user_id_user integer NOT NULL,
    lastusemagnet date
)
INHERITS (type_device);


ALTER TABLE device OWNER TO wibo;

--
-- Name: device_id_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE device_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE device_id_seq OWNER TO wibo;

--
-- Name: device_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE device_id_seq OWNED BY device.id;


--
-- Name: followed; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE followed (
    id integer NOT NULL,
    container_id integer NOT NULL,
    device_id integer,
    iduser integer
);


ALTER TABLE followed OWNER TO wibo;

--
-- Name: followed_id_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE followed_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE followed_id_seq OWNER TO wibo;

--
-- Name: followed_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE followed_id_seq OWNED BY followed.id;


--
-- Name: type_message; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE type_message (
    id_type_m integer NOT NULL,
    typename character varying(255) NOT NULL
);


ALTER TABLE type_message OWNER TO wibo;

--
-- Name: message; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE message (
    id integer NOT NULL,
    content text NOT NULL,
    containerid integer NOT NULL,
    device_id integer NOT NULL,
    creationdate timestamp with time zone DEFAULT now()
)
INHERITS (type_message);


ALTER TABLE message OWNER TO wibo;

--
-- Name: message_id_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE message_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE message_id_seq OWNER TO wibo;

--
-- Name: message_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE message_id_seq OWNED BY message.id;


--
-- Name: reception; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE reception (
    id integer NOT NULL,
    receptiontime date NOT NULL,
    location_rc geography(Point,4326),
    idcontainer integer NOT NULL,
    device_id integer NOT NULL
);


ALTER TABLE reception OWNER TO wibo;

--
-- Name: reception_id_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE reception_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE reception_id_seq OWNER TO wibo;

--
-- Name: reception_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE reception_id_seq OWNED BY reception.id;


--
-- Name: type_information; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE type_information (
    id_type_info integer NOT NULL,
    name_info character varying(255) NOT NULL
);


ALTER TABLE type_information OWNER TO wibo;

--
-- Name: shared; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE shared (
    id integer NOT NULL,
    type_shared integer NOT NULL,
    device_id integer NOT NULL
)
INHERITS (type_information);


ALTER TABLE shared OWNER TO wibo;

--
-- Name: shared_id_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE shared_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE shared_id_seq OWNER TO wibo;

--
-- Name: shared_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE shared_id_seq OWNED BY shared.id;


--
-- Name: type_group; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE type_group (
    id_type_g integer NOT NULL,
    groupname character varying(20) NOT NULL
);


ALTER TABLE type_group OWNER TO wibo;

--
-- Name: user; Type: TABLE; Schema: public; Owner: wibo; Tablespace: 
--

CREATE TABLE "user" (
    id_user integer NOT NULL,
    login character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    salt character varying(255) NOT NULL,
    lastlogin date NOT NULL,
    creationdate date NOT NULL,
    mail character varying(510) NOT NULL
)
INHERITS (type_group);


ALTER TABLE "user" OWNER TO wibo;

--
-- Name: user_id_user_seq; Type: SEQUENCE; Schema: public; Owner: wibo
--

CREATE SEQUENCE user_id_user_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_id_user_seq OWNER TO wibo;

--
-- Name: user_id_user_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wibo
--

ALTER SEQUENCE user_id_user_seq OWNED BY "user".id_user;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY checkpoints ALTER COLUMN id SET DEFAULT nextval('checkpoints_id_seq'::regclass);


--
-- Name: id_type_c; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY container ALTER COLUMN id_type_c SET DEFAULT 1;


--
-- Name: typename; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY container ALTER COLUMN typename SET DEFAULT 'red'::character varying;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY container ALTER COLUMN id SET DEFAULT nextval('container_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY device ALTER COLUMN id SET DEFAULT nextval('device_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY followed ALTER COLUMN id SET DEFAULT nextval('followed_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY message ALTER COLUMN id SET DEFAULT nextval('message_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY reception ALTER COLUMN id SET DEFAULT nextval('reception_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY shared ALTER COLUMN id SET DEFAULT nextval('shared_id_seq'::regclass);


--
-- Name: id_user; Type: DEFAULT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY "user" ALTER COLUMN id_user SET DEFAULT nextval('user_id_user_seq'::regclass);


--
-- Data for Name: checkpoints; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY checkpoints (id, date, containerid, attractbymagnet, location_ckp) FROM stdin;
\.


--
-- Name: checkpoints_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('checkpoints_id_seq', 1, false);


--
-- Data for Name: container; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY container (id_type_c, typename, id, direction, speed, device_id, location_ct, idcreator, titlename, ianix, creationdate) FROM stdin;
1	text	10	123.00	222	2	0101000020E61000000000000000805BC00000000000003E40	2	myball	1	2015-08-04 13:05:51.943019+00
1	text	11	123.00	222	2	0101000020E61000000000000000805BC00000000000003E40	3	myball2	1	2015-08-04 13:13:10.401201+00
1	red	12	123.00	222	2	0101000020E61000000000000000805BC00000000000003E40	3	myball42	1	2015-08-04 13:14:04.278887+00
1	red	13	22.30	222	3	0101000020E61000000000000000805BC00000000000003E40	2	tatat	3	2015-08-04 13:46:43.184382+00
1	red	14	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	DAVE	5	2015-08-04 16:10:36.265758+00
1	red	15	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	DAVE	5	2015-08-05 10:30:04.323494+00
1	red	16	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	DAVE	5	2015-08-05 14:04:36.062689+00
1	red	17	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	DAVE	5	2015-08-05 14:05:03.376293+00
1	red	18	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	DAVE	5	2015-08-05 14:53:39.149932+00
1	red	19	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	DAVE	5	2015-08-05 14:53:48.349937+00
1	red	20	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	MonPremieurBallon	5	2015-08-05 14:54:23.779848+00
1	red	21	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	MonPremieurBallon	5	2015-08-05 14:55:30.587851+00
1	red	22	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	MonPremieurBallon	5	2015-08-05 14:57:53.649374+00
1	red	27	22.00	222	3	0101000020E610000000000000000059C00000000000003E40	2	test	3	2015-08-05 15:28:20.161492+00
1	red	28	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	MonPremieurBallon	5	2015-08-05 15:28:31.252488+00
1	red	29	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	MonPremieurBallon	5	2015-08-05 15:29:29.507925+00
1	red	30	23.90	222	3	0101000020E61000000000000000805BC00000000000003E40	2	MonPremieurBallon	5	2015-08-05 15:30:13.144889+00
\.


--
-- Name: container_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('container_id_seq', 30, true);


--
-- Data for Name: device; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY device (id_type_d, typename, id, macaddr, user_id_user, lastusemagnet) FROM stdin;
1	testdevice1	1	2222	1	1971-07-13
1	testdevice2	2	2222	2	1971-07-13
1	testdevice2	3	2222	3	1971-07-13
1	testdevice2	4	2222	4	1971-07-13
1	testdevice1	5	2222	1	1971-07-13
1	testdevice2	6	2222	2	1971-07-13
1	testdevice2	7	2222	3	1971-07-13
1	testdevice2	8	2222	4	1971-07-13
1	testdevice1	9	2222	1	1971-07-13
1	testdevice2	10	2222	2	1971-07-13
1	testdevice2	11	2222	3	1971-07-13
1	testdevice2	12	2222	4	1971-07-13
\.


--
-- Name: device_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('device_id_seq', 12, true);


--
-- Data for Name: followed; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY followed (id, container_id, device_id, iduser) FROM stdin;
\.


--
-- Name: followed_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('followed_id_seq', 1, false);


--
-- Data for Name: message; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY message (id_type_m, typename, id, content, containerid, device_id, creationdate) FROM stdin;
1	text	3	osar, querer, poder, callar	10	2	2015-07-29 08:54:51.262523+00
1	text	4	osar, querer, poder, callar	10	2	2015-07-29 08:56:34.407764+00
\.


--
-- Name: message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('message_id_seq', 4, true);


--
-- Data for Name: reception; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY reception (id, receptiontime, location_rc, idcontainer, device_id) FROM stdin;
\.


--
-- Name: reception_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('reception_id_seq', 1, false);


--
-- Data for Name: shared; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY shared (id_type_info, name_info, id, type_shared, device_id) FROM stdin;
\.


--
-- Name: shared_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('shared_id_seq', 1, false);


--
-- Data for Name: spatial_ref_sys; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY spatial_ref_sys (srid, auth_name, auth_srid, srtext, proj4text) FROM stdin;
\.


--
-- Data for Name: type_container; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY type_container (id_type_c, typename) FROM stdin;
\.


--
-- Data for Name: type_device; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY type_device (id_type_d, typename) FROM stdin;
\.


--
-- Data for Name: type_group; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY type_group (id_type_g, groupname) FROM stdin;
\.


--
-- Data for Name: type_information; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY type_information (id_type_info, name_info) FROM stdin;
\.


--
-- Data for Name: type_message; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY type_message (id_type_m, typename) FROM stdin;
\.


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: wibo
--

COPY "user" (id_type_g, groupname, id_user, login, password, salt, lastlogin, creationdate, mail) FROM stdin;
1	particuler	1	testlogin1	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	2	testlogin2	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	3	testlogin3	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	4	testlogin4	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	5	testlogin1	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	6	testlogin2	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	7	testlogin3	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	8	testlogin4	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	9	testlogin1	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	10	testlogin2	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	11	testlogin3	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	12	testlogin4	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	13	testlogin1	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	14	testlogin2	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
1	particuler	15	testlogin3	testpass	12345	1971-07-13	1971-07-13	jasds@test.com
\.


--
-- Name: user_id_user_seq; Type: SEQUENCE SET; Schema: public; Owner: wibo
--

SELECT pg_catalog.setval('user_id_user_seq', 17, true);


--
-- Name: checkpoints_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY checkpoints
    ADD CONSTRAINT checkpoints_pk PRIMARY KEY (id);


--
-- Name: container_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY container
    ADD CONSTRAINT container_pk PRIMARY KEY (id);


--
-- Name: device_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY device
    ADD CONSTRAINT device_pk PRIMARY KEY (id);


--
-- Name: followed_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY followed
    ADD CONSTRAINT followed_pk PRIMARY KEY (id);


--
-- Name: message_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY message
    ADD CONSTRAINT message_pk PRIMARY KEY (id);


--
-- Name: reception_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY reception
    ADD CONSTRAINT reception_pk PRIMARY KEY (id);


--
-- Name: shared_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY shared
    ADD CONSTRAINT shared_pk PRIMARY KEY (id);


--
-- Name: type_container_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY type_container
    ADD CONSTRAINT type_container_pk PRIMARY KEY (id_type_c);


--
-- Name: type_device_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY type_device
    ADD CONSTRAINT type_device_pk PRIMARY KEY (id_type_d);


--
-- Name: type_group_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY type_group
    ADD CONSTRAINT type_group_pk PRIMARY KEY (id_type_g);


--
-- Name: type_information_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY type_information
    ADD CONSTRAINT type_information_pk PRIMARY KEY (id_type_info);


--
-- Name: type_message_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY type_message
    ADD CONSTRAINT type_message_pk PRIMARY KEY (id_type_m);


--
-- Name: user_pk; Type: CONSTRAINT; Schema: public; Owner: wibo; Tablespace: 
--

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_pk PRIMARY KEY (id_user);


--
-- Name: checkpoints_container_id_idx; Type: INDEX; Schema: public; Owner: wibo; Tablespace: 
--

CREATE INDEX checkpoints_container_id_idx ON checkpoints USING btree (containerid, date);


--
-- Name: message_idx_container; Type: INDEX; Schema: public; Owner: wibo; Tablespace: 
--

CREATE INDEX message_idx_container ON message USING btree (containerid);


--
-- Name: checkpoints_container; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY checkpoints
    ADD CONSTRAINT checkpoints_container FOREIGN KEY (containerid) REFERENCES container(id);


--
-- Name: container_device; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY container
    ADD CONSTRAINT container_device FOREIGN KEY (device_id) REFERENCES device(id);


--
-- Name: container_user; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY container
    ADD CONSTRAINT container_user FOREIGN KEY (idcreator) REFERENCES "user"(id_user);


--
-- Name: device_user; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY device
    ADD CONSTRAINT device_user FOREIGN KEY (user_id_user) REFERENCES "user"(id_user);


--
-- Name: followed_container; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY followed
    ADD CONSTRAINT followed_container FOREIGN KEY (container_id) REFERENCES container(id);


--
-- Name: followed_device; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY followed
    ADD CONSTRAINT followed_device FOREIGN KEY (device_id) REFERENCES device(id);


--
-- Name: message_container; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY message
    ADD CONSTRAINT message_container FOREIGN KEY (containerid) REFERENCES container(id);


--
-- Name: message_device; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY message
    ADD CONSTRAINT message_device FOREIGN KEY (device_id) REFERENCES device(id);


--
-- Name: reception_device; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY reception
    ADD CONSTRAINT reception_device FOREIGN KEY (device_id) REFERENCES device(id);


--
-- Name: shared_device; Type: FK CONSTRAINT; Schema: public; Owner: wibo
--

ALTER TABLE ONLY shared
    ADD CONSTRAINT shared_device FOREIGN KEY (device_id) REFERENCES device(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

