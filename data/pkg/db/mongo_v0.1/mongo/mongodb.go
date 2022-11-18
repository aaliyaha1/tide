package mongo

import (
	"github.com/bitly/go-simplejson"
	"github.com/globalsign/mgo"
	"log"
	"time"
)

type MongoDB struct {
	session    *mgo.Session
	db         string
	collection string
}

func (m *MongoDB) connect(db, collection string) (*mgo.Session, *mgo.Collection) {

	if db == "" {
		db = m.db
	}

	if collection == "" {
		collection = m.collection
	}

	ms := m.session.Copy()
	c := ms.DB(db).C(collection)
	ms.SetMode(mgo.Monotonic, true)
	return ms, c
}

func (m *MongoDB) getDb(db string) (*mgo.Session, *mgo.Database) {

	if db == "" {
		db = m.db
	}

	ms := m.session.Copy()
	return ms, ms.DB(db)
}

func (m *MongoDB) IsEmpty(db, collection string) bool {

	ms, c := m.connect(db, collection)
	defer ms.Close()
	count, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

func (m *MongoDB) Count(db, collection string, query interface{}) (int, error) {

	ms, c := m.connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

func (m *MongoDB) Insert(db, collection string, docs ...interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Insert(docs...)
}

func (m *MongoDB) FindOne(db, collection string, query, selector, result interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).One(result)
}

func (m *MongoDB) FindAll(db, collection, filed string, query, selector, result interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Find(query).Sort(filed).Select(selector).All(result)
}

func (m *MongoDB) FindLimit(db, collection string, limit int, query, selector, result interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Limit(limit).All(result)
}

func (m *MongoDB) FindPage(db, collection string, page, limit int, query, selector, result interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Skip(page * limit).Limit(limit).All(result)
}

func (m *MongoDB) FindIter(db, collection string, query interface{}) *mgo.Iter {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Find(query).Iter()
}

func (m *MongoDB) Update(db, collection string, selector, update interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Update(selector, update)
}

func (m *MongoDB) Upsert(db, collection string, selector, update interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

func (m *MongoDB) UpdateAll(db, collection string, selector, update interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	_, err := c.UpdateAll(selector, update)
	return err
}

func (m *MongoDB) Remove(db, collection string, selector interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	return c.Remove(selector)
}

func (m *MongoDB) RemoveAll(db, collection string, selector interface{}) error {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	_, err := c.RemoveAll(selector)
	return err
}

// BulkInsert insert one or multi documents
func (m *MongoDB) BulkInsert(db, collection string, docs ...interface{}) (*mgo.BulkResult, error) {

	ms, c := m.connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Insert(docs...)
	return bulk.Run()
}

func (m *MongoDB) BulkRemove(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {

	ms, c := m.connect(db, collection)
	defer ms.Close()

	bulk := c.Bulk()
	bulk.Remove(selector...)
	return bulk.Run()
}

func (m *MongoDB) BulkRemoveAll(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {

	ms, c := m.connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.RemoveAll(selector...)
	return bulk.Run()
}

func (m *MongoDB) BulkUpdate(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {

	ms, c := m.connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Update(pairs...)
	return bulk.Run()
}

func (m *MongoDB) BulkUpdateAll(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := m.connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.UpdateAll(pairs...)
	return bulk.Run()
}

func (m *MongoDB) BulkUpsert(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := m.connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Upsert(pairs...)
	return bulk.Run()
}

func (m *MongoDB) PipeAll(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := m.connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.All(result)
}

func (m *MongoDB) PipeOne(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := m.connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.One(result)
}

func (m *MongoDB) PipeIter(db, collection string, pipeline interface{}, allowDiskUse bool) *mgo.Iter {
	ms, c := m.connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}

	return pipe.Iter()

}

func (m *MongoDB) Explain(db, collection string, pipeline, result interface{}) error {
	ms, c := m.connect(db, collection)
	defer ms.Close()
	pipe := c.Pipe(pipeline)
	return pipe.Explain(result)
}
func (m *MongoDB) GridFSCreate(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := m.getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Create(name)
}

func (m *MongoDB) GridFSFindOne(db, prefix string, query, result interface{}) error {
	ms, d := m.getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).One(result)
}

func (m *MongoDB) GridFSFindAll(db, prefix string, query, result interface{}) error {
	ms, d := m.getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).All(result)
}

func (m *MongoDB) GridFSOpen(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := m.getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Open(name)
}

func (m *MongoDB) GridFSRemove(db, prefix, name string) error {
	ms, d := m.getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Remove(name)
}

func NewMongoDB(params []byte) (*MongoDB, error) {

	sj, err := simplejson.NewJson(params)
	if err != nil {
		return nil, err
	}

	url := sj.Get("url").MustString()
	authDb := sj.Get("authDb").MustString()
	user := sj.Get("user").MustString()
	password := sj.Get("password").MustString()
	db := sj.Get("db").MustString()
	collection := sj.Get("table").MustString()

	dialInfo := &mgo.DialInfo{
		Addrs:     []string{url},
		Timeout:   60 * time.Second,
		Source:    authDb,
		Username:  user,
		Password:  password,
		PoolLimit: 4096,
	}

	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	mongo := new(MongoDB)
	mongo.session = s

	if db != "" && collection != "" {
		mongo.db = db
		mongo.collection = collection
	}

	return mongo, nil
}

func NewMongoDBWithoutBytes(url, authDb, user, password, db, collection string) (*MongoDB, error) {
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{url},
		Timeout:   60 * time.Second,
		Source:    authDb,
		Username:  user,
		Password:  password,
		PoolLimit: 4096,
	}

	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	mongo := new(MongoDB)
	mongo.session = s
	if db != "" && collection != "" {
		mongo.db = db
		mongo.collection = collection
	}

	return mongo, nil
}
