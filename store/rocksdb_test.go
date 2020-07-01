package store

import (
	"os"
	"testing"
)

const (
	DB_PATH = "/tmp/gorocksdb_test"
	value   = "value_test"
)

var key = []byte{}

func TestInit(t *testing.T) {
	// check directory, if exists, delete it.
	if _, err := os.Stat(DB_PATH); err == nil {
		t.Logf("Directory %s exists, and now it is deleted.\n", DB_PATH)
		if err = os.RemoveAll(DB_PATH); err != nil {
			t.Fatalf("Delete %s failed: %v\n", DB_PATH, err)
		}
	}

	// open db
	store, err := Init(DB_PATH)
	if err != nil {
		t.Fatalf("Open database failed!")
	}
	defer store.CloseDB()

	// test directory
	if _, err := os.Stat(DB_PATH); err != nil {
		t.Fatalf("Create db failed.\n")
	}
}

func TestSave(t *testing.T) {
	// open db
	store, err := Init(DB_PATH)
	if err != nil {
		t.Fatalf("Open database failed!")
	}
	defer store.CloseDB()

	// save
	err = store.Save([]byte("key"), []byte("value"))
	if err != nil {
		t.Fatalf("Save data to database failed: %v\n", err)
	}

	key, err = store.SaveValue([]byte(value))
	if err != nil {
		t.Fatalf("Save data '%v' to database failed\n", value)
	} else {
		t.Logf("Key of value '%v' is '%v'\n", value, key)
	}
}

func TestGet(t *testing.T) {
	// open db
	store, err := Init(DB_PATH)
	if err != nil {
		t.Fatalf("Open database failed!")
	}
	defer store.CloseDB()

	// get
	value2, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get data '%v' from database failed: %v\n", key, err)
	}

	if value != string(value2) {
		t.Fatalf("Not equal: '%v' -- '%v'\n", value, value2)
	}
}

/*func test() {
	// delete old directory

	db, err := OpenDB(DB_PATH)
	if err != nil {
		log.Println("fail to open db,", nil, db)
	}

	readOptions := gorocksdb.NewDefaultReadOptions()
	readOptions.SetFillCache(true)

	writeOptions := gorocksdb.NewDefaultWriteOptions()
	writeOptions.SetSync(true)

	for i := 0; i < 10000; i++ {
		keyStr := "aa" + strconv.Itoa(i)
		var key []byte = []byte(keyStr)
		db.Put(writeOptions, key, key)
		log.Println(i, keyStr)
		slice, err2 := db.Get(readOptions, key)
		if err2 != nil {
			log.Println("获取数据异常：", key, err2)
			continue
		}
		log.Println("获取数据：", slice.Size(), string(slice.Data()))
	}

	//defer readOptions.Destroy()
	//defer writeOptions.Destroy()
}*/
