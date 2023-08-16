package main

import (
	"fmt"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println(os.Args)
	dbPath := os.Args[1]
	var src, des string
	if len(os.Args) >= 4 {
		src, des = os.Args[2], os.Args[3]
	} else {
		src, des = "http:", "https:"
	}
	log.Println(dbPath)
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 3, ReadOnly: false})
	if db != nil {
		defer db.Close()
	}
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	bucketName := "members"
	members := make(map[string]string)
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.ForEach(func(k, v []byte) error {
			members[string(k)] = string(v)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	for k, v := range members {
		fmt.Printf("<%s,%s>\n", k, v)
		members[k] = strings.ReplaceAll(v, src, des)
	}
	err = db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		for k, v := range members {
			err := b.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.ForEach(func(k, v []byte) error {
			log.Printf("new <%s,%s>", string(k), string(v))
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
}
