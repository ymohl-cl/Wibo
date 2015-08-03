package ballon_test

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
	_, err = db.Query("INSERT INTO \"user\" (login, password, salt, lastlogin, creationDate, mail, id_type_g, groupName) VALUES ('mathieudeclemy', 'testpass', '12345', '1971-07-13', '1971-07-13', 'mathieudeclemy@hotmail.fr', 1, 'particuler');")
	assert.Nil(t, err)
	_, err = db.Query("INSERT INTO container (direction, speed, id_type_c, typename, creationDate, device_id, location_ct) VALUES (23.9, 222, 4, 'text', '1971-07-13', 2, ST_GeographyFromText('SRID=4326; POINT(-110 30)'), );")
	if err != nil {
		fmt.Println(err.Error())
		assert.Nil(t, err)
	}
}
