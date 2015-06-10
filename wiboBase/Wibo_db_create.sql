-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2015-06-10 07:23:13.02




-- tables
-- Table: Checkpoints
CREATE TABLE Checkpoints (
    id int  NOT NULL,
    longitude decimal(5,2)  NOT NULL,
    latitude decimal(5,2)  NOT NULL,
    date date  NOT NULL,
    containerId int  NOT NULL,
    attractByMagnet boolean  NOT NULL,
    CONSTRAINT Checkpoints_pk PRIMARY KEY (id)
);

CREATE INDEX Checkpoints_container_id_idx on Checkpoints (containerId ASC,date ASC);




-- Table: Container
CREATE TABLE Container (
    id int  NOT NULL,
    longitude decimal(5,2)  NOT NULL,
    latitude decimal(5,2)  NOT NULL,
    direction decimal(5,2)  NOT NULL,
    speed int  NOT NULL,
    creationDate date  NOT NULL,
    userId int  NOT NULL,
    typeContainerId int  NOT NULL,
    titre varchar(255) NOT NULL, 
    CONSTRAINT Container_pk PRIMARY KEY (id)
);



-- Table: Device
CREATE TABLE Device (
    id int  NOT NULL,
    macAddr varchar(18)  NOT NULL,
    typeDeviceId int  NOT NULL,
    idUser int  NOT NULL,
    user_id int  NOT NULL,
    CONSTRAINT Device_pk PRIMARY KEY (id)
);



-- Table: Followed
CREATE TABLE Followed (
    id int  NOT NULL,
    User_id int  NOT NULL,
    Container_id int  NOT NULL,
    CONSTRAINT Followed_pk PRIMARY KEY (id)
);



-- Table: "Group"
CREATE TABLE "Group" (
    id int  NOT NULL,
    groupName varchar(20)  NOT NULL,
    CONSTRAINT Group_pk PRIMARY KEY (id)
);



-- Table: Message
CREATE TABLE Message (
    id int  NOT NULL,
    content text  NOT NULL,
    containerId int  NOT NULL,
    userId int  NOT NULL,
    typeMessageId int  NOT NULL,
    CONSTRAINT Message_pk PRIMARY KEY (id)
);

CREATE INDEX Message_idx_container on Message (containerId ASC);




-- Table: Reception
CREATE TABLE Reception (
    id int  NOT NULL,
    receptionTime date  NOT NULL,
    longitude decimal(5,2)  NOT NULL,
    latitude decimal(5,2)  NOT NULL,
    userId int  NOT NULL,
    IdContainer int  NOT NULL,
    CONSTRAINT Reception_pk PRIMARY KEY (id)
);



-- Table: Session
CREATE TABLE Session (
    id int  NOT NULL,
    logged boolean  NOT NULL,
    loginTime date  NOT NULL,
    logoutTime boolean  NOT NULL,
    userId int  NULL,
    CONSTRAINT Session_pk PRIMARY KEY (id)
);



-- Table: Type_Container
CREATE TABLE Type_Container (
    id int  NOT NULL,
    typeName varchar(255)  NOT NULL,
    CONSTRAINT Type_Container_pk PRIMARY KEY (id)
);



-- Table: Type_Device
CREATE TABLE Type_Device (
    id int  NOT NULL,
    typeName varchar(255)  NOT NULL,
    CONSTRAINT Type_Device_pk PRIMARY KEY (id)
);



-- Table: Type_Message
CREATE TABLE Type_Message (
    id int  NOT NULL,
    typeName varchar(255)  NOT NULL,
    CONSTRAINT Type_Message_pk PRIMARY KEY (id)
);



-- Table: shared
CREATE TABLE shared (
    id int  NOT NULL,
    type_shared int  NOT NULL,
    User_id int  NOT NULL,
    CONSTRAINT shared_pk PRIMARY KEY (id)
);



-- Table: "user"
CREATE TABLE "user" (
    id int  NOT NULL,
    login varchar(255)  NOT NULL,
    password varchar(255)  NOT NULL,
    salt varchar(255)  NOT NULL,
    lastLogin date  NOT NULL,
    creationDate date  NOT NULL,
    mail varchar(510)  NOT NULL,
    groupId int  NOT NULL,
    UseMagnet date  NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (id)
);

-- Table: "container_type_information"
CREATE TABLE container_type_information (
    id_type integer NOT NULL,
    name_type character varying(255) NOT NULL,
    CONSTRAINT container_type_information_pk PRIMARY KEY (id_type)
);



-- foreign keys
-- Reference:  Checkpoints_Container (table: Checkpoints)


ALTER TABLE Checkpoints ADD CONSTRAINT Checkpoints_Container 
    FOREIGN KEY (containerId)
    REFERENCES Container (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Container_Type_Container (table: Container)


ALTER TABLE Container ADD CONSTRAINT Container_Type_Container 
    FOREIGN KEY (typeContainerId)
    REFERENCES Type_Container (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Container_User (table: Container)


ALTER TABLE Container ADD CONSTRAINT Container_User 
    FOREIGN KEY (userId)
    REFERENCES "user" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Device_Type_Device (table: Device)


ALTER TABLE Device ADD CONSTRAINT Device_Type_Device 
    FOREIGN KEY (typeDeviceId)
    REFERENCES Type_Device (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Device_user (table: Device)


ALTER TABLE Device ADD CONSTRAINT Device_user 
    FOREIGN KEY (user_id)
    REFERENCES "user" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Followed_Container (table: Followed)


ALTER TABLE Followed ADD CONSTRAINT Followed_Container 
    FOREIGN KEY (Container_id)
    REFERENCES Container (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Followed_User (table: Followed)


ALTER TABLE Followed ADD CONSTRAINT Followed_User 
    FOREIGN KEY (User_id)
    REFERENCES "user" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Message_Container (table: Message)


ALTER TABLE Message ADD CONSTRAINT Message_Container 
    FOREIGN KEY (containerId)
    REFERENCES Container (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Message_Type_message (table: Message)


ALTER TABLE Message ADD CONSTRAINT Message_Type_message 
    FOREIGN KEY (typeMessageId)
    REFERENCES Type_Message (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Message_User (table: Message)


ALTER TABLE Message ADD CONSTRAINT Message_User 
    FOREIGN KEY (userId)
    REFERENCES "user" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Reception_User (table: Reception)


ALTER TABLE Reception ADD CONSTRAINT Reception_User 
    FOREIGN KEY (userId)
    REFERENCES "user" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  "Session_User" (table: Session)


ALTER TABLE Session ADD CONSTRAINT "Session_User" 
    FOREIGN KEY (userId)
    REFERENCES "user" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference:  Shared_User (table: shared)


ALTER TABLE shared ADD CONSTRAINT Shared_User 
    FOREIGN KEY (User_id)
    REFERENCES "user" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;

-- Reference: Shared_container_type_information 

ALTER TABLE container_type_information ADD CONSTRAINT Shared_container_type_information 
    FOREIGN KEY (type_shared)
    REFERENCES container_type_information (id_type)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference:  User_Group (table: "user")


ALTER TABLE "user" ADD CONSTRAINT User_Group 
    FOREIGN KEY (groupId)
    REFERENCES "Group" (id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE 
;






-- End of file.

