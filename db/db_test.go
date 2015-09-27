package db_test

/*
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

// request data base test
func TestDBConnections(t *testing.T) {
	db, err := sql.Open("postgres", "user=wibo  password='wibo' dbname=wibo_base sslmode=disable host=localhost port=49155")

	defer db.Close()

	assert.Nil(t, err)
	//insert users test
	_, err = db.Query("INSERT INTO \"user\" (login, password, salt, lastlogin, creationDate, mail, id_type_g, groupName) VALUES ('testlogin1', 'testpass', '12345', '1971-07-13', '1971-07-13', 'jasds@test.com', 1, 'particuler');")
	assert.Nil(t, err)
	_, err = db.Query("INSERT INTO \"user\" (login, password, salt, lastlogin, creationDate, mail, id_type_g, groupName) VALUES ('testlogin2', 'testpass', '12345', '1971-07-13', '1971-07-13', 'jasds@test.com', 1, 'particuler');")
	assert.Nil(t, err)
	_, err = db.Query("INSERT INTO \"user\" (login, password, salt, lastlogin, creationDate, mail, id_type_g, groupName) VALUES ('testlogin3', 'testpass', '12345', '1971-07-13', '1971-07-13', 'jasds@test.com', 1, 'particuler');")
	_, err = db.Query("INSERT INTO \"user\" (login, password, salt, lastlogin, creationDate, mail, id_type_g, groupName) VALUES ('testlogin4', 'testpass', '12345', '1971-07-13', '1971-07-13', 'jasds@test.com', 1, 'particuler');")
	assert.Nil(t, err)
	// insert device
	_, err = db.Query("INSERT INTO device (macAddr, user_id_user, lastUseMagnet, id_type_d, typeName) VALUES ('2222', 1, '1971-07-13', 1, 'testdevice1');")
	_, err = db.Query("INSERT INTO device (macAddr, user_id_user, lastUseMagnet, id_type_d, typeName) VALUES ('2222', 2, '1971-07-13', 1, 'testdevice2');")
	_, err = db.Query("INSERT INTO device (macAddr, user_id_user, lastUseMagnet, id_type_d, typeName) VALUES ('2222', 3, '1971-07-13', 1, 'testdevice2');")
	_, err = db.Query("INSERT INTO device (macAddr, user_id_user, lastUseMagnet, id_type_d, typeName) VALUES ('2222', 4, '1971-07-13', 1, 'testdevice2');")
	// insert container

	_, err = db.Query("INSERT INTO container (direction, speed, id_type_c, typename, creationDate, device_id, location_ct) VALUES (23.9, 222, 1, 'testdevice', '1971-07-13', 2, ST_GeographyFromText('SRID=4326; POINT(-110 30)'));")
	_, err = db.Query("INSERT INTO container (direction, speed, id_type_c, typename, creationDate, device_id, location_ct) VALUES (23.9, 222, 2, 'testdevice', '1971-07-13', 2, ST_GeographyFromText('SRID=4326; POINT(-110 30)'));")
	_, err = db.Query("INSERT INTO container (direction, speed, id_type_c, typename, creationDate, device_id, location_ct) VALUES (23.9, 222, 3, 'testdevice', '1971-07-13', 2, ST_GeographyFromText('SRID=4326; POINT(-110 30)'));")
	_, err = db.Query("INSERT INTO container (direction, speed, id_type_c, typename, creationDate, device_id, location_ct) VALUES (23.9, 222, 4, 'testdevice', '1971-07-13', 2, ST_GeographyFromText('SRID=4326; POINT(-110 30)'));")
	if err != nil {
		fmt.Println(err.Error())
		assert.Nil(t, err)
	}
}*/
