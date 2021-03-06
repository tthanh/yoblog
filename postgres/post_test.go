package postgres

import (
	"reflect"
	"testing"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	uuid "github.com/satori/go.uuid"
	"github.com/tthanh/yoblog"
)

func Test_Post(t *testing.T) {
	now := time.Now()

	account := &yoblog.Account{
		ID:        uuid.NewV4().String(),
		Email:     randomdata.Email(),
		Name:      randomdata.FirstName(randomdata.Female),
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	accountID, err := accountStore.Create(account)
	if err != nil {
		t.Fatal(err)
	}

	account.ID = accountID

	post := &yoblog.Post{
		OwnerID:   account.ID,
		OwnerName: account.Name,
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	t.Run("Create", func(t *testing.T) {
		postID, err := postStore.Create(post)
		if err != nil {
			t.Fatal(err)
		}

		post.ID = postID
	})

	t.Run("GetByID", func(t *testing.T) {
		_post, err := postStore.GetByID(post.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(*post, _post) {
			t.Fatalf("Mismatch: %v != %v", *post, _post)
		}
	})

	t.Run("GetByOwnerID", func(t *testing.T) {
		posts, err := postStore.GetByOwnerID(account.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(posts) != 1 {
			t.Fatalf("Expected 1. Got %d", len(posts))
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		posts, err := postStore.GetAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(posts) == 0 {
			t.Fatalf("Expected not empty")
		}
	})

	comment := &yoblog.Comment{
		OwnerID: account.ID,
		PostID:  post.ID,
		Content: "Lorem Ipsum",
	}
	t.Run("CreateComment", func(t *testing.T) {
		_, err := postStore.CreateComment(comment)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetPostComments", func(t *testing.T) {
		comments, err := postStore.GetPostComments(post.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(comments) <= 0 {
			t.Fatal("Expected not empty")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := postStore.Delete(post.ID)
		if err != nil {
			t.Fatal(err)
		}
	})
}
