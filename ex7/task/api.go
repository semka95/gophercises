package task

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Storage represents storage for tasks
type Storage interface {
	GetAll() ([]Task, error)
	Store(task string) (int, error)
	Delete(id int) error
	Complete(id int) error
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
			task := Task{}

			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}

			tasks = append(tasks, task)
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
	t := Task{
		Value: task,
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		id64, err := b.NextSequence()
		if err != nil {
			return err
		}
		t.ID = int(id64)

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put(itob(t.ID), buf)
	})
	if err != nil {
		return -1, err
	}

	return t.ID, nil
}

// Delete deletes task from database by given id
func (s BoltStore) Delete(id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		return b.Delete(itob(id))
	})
}

// Complete marks task as completed by given id
func (s BoltStore) Complete(id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)

		task := new(Task)
		data := b.Get(itob(id))
		if len(data) == 0 {
			return fmt.Errorf("task %d does not exist", id)
		}

		err := json.Unmarshal(data, &task)
		if err != nil {
			return err
		}

		task.CompletedAt = time.Now()
		enc, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put(itob(id), enc)
	})
}

// itob returns an 8-byte big endian representation of v
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
