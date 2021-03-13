package task

import (
	"encoding/binary"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Storage represents storage for tasks
type Storage interface {
	GetAll() ([]Task, error)
	Store(task string) (int, error)
	Delete(id int) error
}

// BoltStore represents BoltDB Storage
type BoltStore struct {
	db         *bolt.DB
	bucketName []byte
}

// NewBoltStore creates new instance of BoltStore
func NewBoltStore(path, bucketName string) (BoltStore, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
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
func (s BoltStore) GetAll() ([]Task, error) {
	var tasks []Task
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				ID:    btoi(k),
				Value: string(v),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Store puts single task to database
func (s BoltStore) Store(task string) (int, error) {
	var id int
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		id64, err := b.NextSequence()
		if err != nil {
			return err
		}
		id = int(id64)

		return b.Put(itob(id), []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

// Delete deletes task from database by given id
func (s BoltStore) Delete(id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		return b.Delete(itob(id))
	})
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
