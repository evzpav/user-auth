package mysql_test

// func TestNew(t *testing.T) {
// 	container := mock.NewPGContainer(t)

// 	defer container.Shutdown()

// 	dbLogTest, err := postgres.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 0, true)
// 	if err != nil {
// 		t.Fatalf("Error establishing connection %v", err)
// 	}
// 	dbLogTest.Close()

// 	db, err := postgres.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 1, true)
// 	if err != nil {
// 		t.Fatalf("Error establishing connection %v", err)
// 	}

// 	var user model.User
// 	db.Select(&user)

// 	assert.NotNil(t, db)

// 	db.Close()

// }
