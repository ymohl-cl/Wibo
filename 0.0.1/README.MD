installation postgis:

You could install from brew:

brew install postgis(not recomended)

or get the source from (recomended)

[post gist](http://postgis.net/install/)

config:

./configure --with-projdir=/nfs/zfs-student-3/users/2013/<your login>/.brew/Cellar/proj/4.9.1/
./configure --with-projdir=/nfs/zfs-student-3/users/2013/<your login>/.brew/Cellar/proj/4.9.1/ --prefix=/nfs/zfs-student-3/users/2013/<your login>/Applications/


CREATE EXTENSION postgis;
SELECT postgis_full_version();


[TEST]

example:

INSERT INTO "user" (login, password, salt, lastlogin, creationDate, mail, id_type_g, groupName) VALUES ('testlogin', 'testpass', '12345',  '1971-07-13', '1971-07-13', 'jasds@test.com', 1, 'particuler');

INSERT INTO device  (macAddr, user_id_user, lastUseMagnet, id_type_d, typeName) VALUES ('2222', 4,  '1971-07-13', 1, 'testdevice');

INSERT INTO container (direction, speed, id_type_c,  typename, creationDate, device_id,  location_ct) VALUES (23.9, 222,  1, 'testdevice', '1971-07-13', 2, ST_GeographyFromText('SRID=4326; POINT(-110 30)'));


2 WAYS  TO REPRESENT SPATIAL OBJETS:
	1. Well-known Text (WKT).
	2. Well-known Binary (WKB).
------ BOTH got information about type of the objet and the coordinates

example: 
		WKT

		POINT(0 0)
		MULTIPOINT(0 0, 1 2)

OPENGIS format requires the spacial referencing system indentifier (SRID) to INSERT.

			INPUT /OUPUT

	bytea WKB = asBinary(geometry); 
	text WKT = asText(geometry); ( I RECOMEND)
	geometry = GeomFromWKB(bytea WKB, SRID); 
	geometry = GeometryFromText(text WKT, SRID);

example sql:

	INSERT INTO container (type_container, positionGeom)
		VALUES ('RED', GeomFromText('POINT(-125.2 45.32)', 4236));

there are two OpenGIS meta-data tables: 
	
	1. SPATIAL_REF_SYS
		
		<<The SPATIAL_REF_SYS table holds the numeric IDs and textual descriptions of coordinate systems used in the spatial database.>>
		In other words the weather stations


	2. GEOMETRY_COLUMNS. 

	[very good explanation](http://postgis.net/docs/AddGeometryColumn.html)



[source](http://postgis.refractions.net/documentation/manual-1.3/ch04.html#id2571655)
[create db postgis](http://workshops.boundlessgeo.com/postgis-intro/creating_db.html)


Commmand to pass shape files:

~/.brew/Cellar/postgresql/9.4.3/bin/shp2pgsql -s 4326 -I FRA_roads.shp  schema.FranceRoutes [redirection to ]  routes_france.sql
