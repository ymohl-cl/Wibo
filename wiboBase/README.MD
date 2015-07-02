
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

CREATE EXTENSION postgis;
SELECT postgis_full_version();
http://workshops.boundlessgeo.com/postgis-intro/creating_db.html
