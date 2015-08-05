-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2015-06-11 08:29:39.327




-- tables

-- Table: type_message
CREATE TABLE type_message (
    id_type_m int  NOT NULL,
    typeName varchar(255)  NOT NULL,
    CONSTRAINT type_message_pk PRIMARY KEY (id_type_m)
);

-- Table: type_container
CREATE TABLE type_container (
    id_type_c int  NOT NULL,
    typeName varchar(255)  NOT NULL,
    CONSTRAINT type_container_pk PRIMARY KEY (id_type_c)
);



-- Table: type_device
CREATE TABLE type_device (
    id_type_d int  NOT NULL,
    typeName varchar(255)  NOT NULL,
    CONSTRAINT type_device_pk PRIMARY KEY (id_type_d)
);



-- Table: type_group
CREATE TABLE type_group (
    id_type_g int  NOT NULL,
    groupName varchar(20)  NOT NULL,
    CONSTRAINT type_group_pk PRIMARY KEY (id_type_g)
);


-- Table: type_information
CREATE TABLE type_information (
    id_type_info int  NOT NULL,
    name_info varchar(255)  NOT NULL,
    CONSTRAINT type_information_pk PRIMARY KEY (id_type_info)
);
-- Table: checkpoints
CREATE TABLE checkpoints (
    id serial  NOT NULL,
    date date  NOT NULL,
    containerId int  NOT NULL,
    attractByMagnet boolean  NOT NULL,
	location_ckp GEOGRAPHY(POINT,4326),
    CONSTRAINT checkpoints_pk PRIMARY KEY (id)
);

CREATE INDEX Checkpoints_container_id_idx on checkpoints (containerId ASC,date ASC);




-- Table: container
CREATE TABLE container (
    id serial  NOT NULL,
    direction decimal(5,2)  NOT NULL,
    speed int  NOT NULL,
    creationDate date  NOT NULL,
    Device_id int  NOT NULL,
	location_ct GEOGRAPHY(POINT,4326),
    CONSTRAINT container_pk PRIMARY KEY (id)
) INHERITS (type_container);



-- Table: device
CREATE TABLE device (
    id serial  NOT NULL,
    macAddr varchar(18)  NOT NULL,
    user_id_user int  NOT NULL,
    lastUseMagnet date  NOT NULL,
    CONSTRAINT device_pk PRIMARY KEY (id)
) INHERITS (type_device);



-- Table: followed
CREATE TABLE followed (
    id serial  NOT NULL,
    Container_id int  NOT NULL,
    Device_id int  NOT NULL,
    CONSTRAINT followed_pk PRIMARY KEY (id)
);



-- Table: message
CREATE TABLE message (
    id serial  NOT NULL,
    content text  NOT NULL,
    containerId int  NOT NULL,
    Device_id int  NOT NULL,
    CONSTRAINT message_pk PRIMARY KEY (id)
) INHERITS (type_message);

CREATE INDEX Message_idx_container on message (containerId ASC);




-- Table: reception
CREATE TABLE reception (
    id serial  NOT NULL,
    receptionTime date  NOT NULL,
	location_rc GEOGRAPHY(POINT,4326),
    IdContainer int  NOT NULL,
    Device_id int  NOT NULL,
    CONSTRAINT reception_pk PRIMARY KEY (id)
);



-- Table: shared
CREATE TABLE shared (
    id serial  NOT NULL,
    type_shared int  NOT NULL,
    Device_id int  NOT NULL,
    CONSTRAINT shared_pk PRIMARY KEY (id)
) INHERITS (type_information);








-- Table: "user"
CREATE TABLE "user" (
    id_user serial  NOT NULL,
    login varchar(255)  NOT NULL,
    password varchar(255)  NOT NULL,
    salt varchar(255)  NOT NULL,
    lastLogin date  NOT NULL,
    creationDate date  NOT NULL,
    mail varchar(510)  NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (id_user)
) INHERITS (type_group);







-- foreign keys
-- Reference:  Checkpoints_Container (table: checkpoints)


ALTER TABLE checkpoints ADD CONSTRAINT Checkpoints_Container
    FOREIGN KEY (containerId)
    REFERENCES container (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  Container_Device (table: container)


ALTER TABLE container ADD CONSTRAINT Container_Device
    FOREIGN KEY (Device_id)
    REFERENCES device (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  Followed_Container (table: followed)


ALTER TABLE followed ADD CONSTRAINT Followed_Container
    FOREIGN KEY (Container_id)
    REFERENCES container (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  Followed_Device (table: followed)


ALTER TABLE followed ADD CONSTRAINT Followed_Device
    FOREIGN KEY (Device_id)
    REFERENCES device (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  Message_Container (table: message)


ALTER TABLE message ADD CONSTRAINT Message_Container
    FOREIGN KEY (containerId)
    REFERENCES container (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  Message_Device (table: message)


ALTER TABLE message ADD CONSTRAINT Message_Device
    FOREIGN KEY (Device_id)
    REFERENCES device (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  Reception_Device (table: reception)


ALTER TABLE reception ADD CONSTRAINT Reception_Device
    FOREIGN KEY (Device_id)
    REFERENCES device (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  device_user (table: device)


ALTER TABLE device ADD CONSTRAINT device_user
    FOREIGN KEY (user_id_user)
    REFERENCES "user" (id_user)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  shared_Device (table: shared)


ALTER TABLE shared ADD CONSTRAINT shared_Device
    FOREIGN KEY (Device_id)
    REFERENCES device (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;






-- End of file.
