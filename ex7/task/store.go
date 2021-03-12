package task

import (
	"encoding/binary"
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// Storage represents storage for tasks
type Storage interface {
	GetAll() (string, error)
	Store(task string) error
	Delete(id int) (string, error)
}

// BoltStore represents BoltDB Storage
type BoltStore struct {
	db         *bolt.DB
	bucketName []byte
}

// NewBoltStore creates new instance of BoltStore
func NewBoltStore(path, bucketName string) (BoltStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return BoltStore{}, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return BoltStore{}, err
	}

	st := BoltStore{
		db:         db,
		bucketName: []byte(bucketName),
	}

	return st, nil
}

// GetAll creates numbered list of all tasks from database
func (s BoltStore) GetAll() (string, error) {
	var res strings.Builder
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Fprintf(&res, "%d. %s\n", btoi(k), string(v))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return res.String(), nil
}

// Store puts single task to database
func (s BoltStore) Store(task string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		return b.Put(itob(int(id)), []byte(task))
	})
}

// Delete deletes task from database by given id
func (s BoltStore) Delete(id int) (string, error) {
	var deleted string
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)

		t := b.Get(itob(id))
		if t == nil {
			return fmt.Errorf("task %d does not exists", id)
		}
		deleted = string(t)

		return b.Delete(itob(id))
	})
	if err != nil {
		return "", err
	}
	return deleted, nil
}

// itob returns an 8-byte big endian representation of v
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi converts 8-byte big endian to int
func btoi(b []byte) int {
	i := binary.BigEndian.Uint64(b)
	return int(i)
}
