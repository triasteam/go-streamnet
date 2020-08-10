package dag

type storage interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
}

type RocksStorage struct {
	rdb    *gorocksdb.DB
	wo     *gorocksdb.WriteOptions
	ro     *gorocksdb.ReadOptions
}

func NewRocksStorage(dir string) *Storage {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	rdb, err := gorocksdb.OpenDb(opts, dir)
	if err != nil {
		panic(err)
	}

	return rdb
}

func New(rdb *gorocksdb.DB) *DB {
	db := &DB{rdb: rdb}
	db.wo = gorocksdb.NewDefaultWriteOptions()
	db.ro = gorocksdb.NewDefaultReadOptions()
	db.RawSet([]byte{MAXBYTE}, nil) // for Enumerator seek to last
	return db
}

func (d *DB) List(key []byte) *ListElement {
	return d.objFromCache(key, LIST).(*ListElement)
}

func (d *DB) FLushAll() {
	// delete all
}

func (d *DB) Keys() []string {
	keyList := []string{}
	batch := gorocksdb.NewWriteBatch()
	defer batch.Destroy()

	d.PrefixEnumerate(KEY, IterForward, func(i int, key, value []byte, quit *bool) {
		keyName, _ := SplitKeyName(key)
		keyList = append(keyList, keyName)
	})

	return keyList
}

func (d *DB) Delete(key []byte) error {
	return d.RawDelete(key)
}

func (d *DB) TypeOf(key []byte) ElementType {
	c := ElementType(NONE)
	prefix := bytes.Join([][]byte{KEY, key, SEP}, nil)
	d.PrefixEnumerate(prefix, IterForward, func(i int, key, value []byte, quit *bool) {
		c = ElementType(key[len(prefix):][0])
		*quit = true
	})
	return c
}

func (d *DB) Get(key []byte) ([]byte, error) {
	return d.RawGet(rawKey(key, STRING))
}

func (d *DB) GetList(key []byte) ([]byte, error) {
	return d.RawGet(rawKey(key, LIST))
}

func (d *DB) Set(key, value []byte) error {
	return d.RawSet(rawKey(key, STRING), value)
}

func (d *DB) WriteBatch(batch *gorocksdb.WriteBatch) error {
	return d.rdb.Write(d.wo, batch)
}

func (d *DB) RawGet(key []byte) ([]byte, error) {
	return d.rdb.GetBytes(d.ro, key)
}

func (d *DB) RawSet(key, value []byte) error {
	return d.rdb.Put(d.wo, key, value)
}

func (d *DB) RawDelete(key []byte) error {
	return d.rdb.Delete(d.wo, key)
}

func (d *DB) Close() {
	d.wo.Destroy()
	d.ro.Destroy()
	d.rdb.Close()
}
