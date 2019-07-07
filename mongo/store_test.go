package mongo

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/xujl1062/store"
)

func TestDb(t *testing.T) {
	var s store.Store
	var err error
	s, err = NewStore("mongodb://test:test@localhost:27017/test", log.New(os.Stdout, "test", 0))
	if err != nil {
		t.Fatalf("Failed create store ,cause %s", err.Error())
	}
	// get
	t.Run("Save", func(t *testing.T) {
		err := s.Save(context.Background(), "test/col", map[string]string{
			"_id":     "1",
			"title":   "hello",
			"content": "hello world ",
		})
		if err != nil {
			t.Fatalf("Failed to save entity . Cause %s", err.Error())
		}
		t.Run("Then get", func(t *testing.T) {
			out := &map[string]string{}
			err := s.Get(context.Background(), "test/col", map[string]string{"_id": "1"}, out)
			if err != nil {
				t.Fatalf("Failed to get the obj , cause %s", err.Error())
			}

		})
	})
}
