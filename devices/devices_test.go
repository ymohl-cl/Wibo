package devices

import (
	"Wibo/db"
	"Wibo/users"
	"container/list"
	"testing"
)

func BenchmarkAddDeviceOnBdd(b *testing.B){
	var err error
		Lst_users := new(users.All_users)
		myDb := new(db.Env)
		Lst_users.Ulist = list.New()
		Db, err := myDb.OpenCo(err)
		LDevices := new(All_Devices)
		LDevices.Dlist = list.New()
		if err != nil {
			b.Fatalf("benchmarkConnection: %s", err)
		}
		for n := 0; n < b.N; n++ {
			 LDevices.AddDeviceOnBdd("20", Lst_users, Db)
			
        }
}