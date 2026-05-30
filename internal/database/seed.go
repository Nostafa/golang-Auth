package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Nostafa/golang-jwt/internal/store"
)

const (
	seedUsersCount      = 100
	seedPostsPerUser    = 10
	seedCommentsPerPost = 10
	seedWorkerCount     = 10
)

func Seed(storage *store.Storage) error {
	ctx := context.Background()

	log.Printf("🌱 Starting database seed (%d users, %d posts/user, %d comments/post, %d workers)",
		seedUsersCount, seedPostsPerUser, seedCommentsPerPost, seedWorkerCount)

	users := generateUsers(seedUsersCount)
	log.Printf("👥 Generated %d user records", len(users))

	var wg sync.WaitGroup
	sem := make(chan struct{}, seedWorkerCount)
	errCh := make(chan error, seedUsersCount)

	for i, user := range users {
		wg.Add(1)

		go func(index int, user *store.User) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			if err := seedUser(ctx, storage, index+1, seedUsersCount, user); err != nil {
				errCh <- err
			}
		}(i, user)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	totalPosts := seedUsersCount * seedPostsPerUser
	totalComments := totalPosts * seedCommentsPerPost
	log.Printf("🎉 Seed completed: %d users, %d posts, %d comments", seedUsersCount, totalPosts, totalComments)

	return nil
}

func seedUser(ctx context.Context, storage *store.Storage, index, total int, user *store.User) error {
	user, err := storage.User.Create(ctx, user)
	if err != nil {
		log.Printf("❌ Failed to create user %d: %v", index, err)
		return err
	}
	log.Printf("✅ User %d/%d created (id=%d, username=%s)", index, total, user.Id, user.Username)

	posts := generatePosts(seedPostsPerUser, user.Id)
	for j, post := range posts {
		post, err := storage.Post.Create(ctx, post)
		if err != nil {
			log.Printf("❌ Failed to create post for user %d: %v", user.Id, err)
			return err
		}

		comments := generateComments(seedCommentsPerPost, post.Id, user.Id)
		for _, comment := range comments {
			if err := storage.Comment.Create(ctx, comment); err != nil {
				log.Printf("❌ Failed to create comment on post %d: %v", post.Id, err)
				return err
			}
		}

		if j == len(posts)-1 {
			log.Printf("📝 User %d: %d posts seeded with %d comments each", user.Id, seedPostsPerUser, seedCommentsPerPost)
		}
	}

	return nil
}

func generateUsers(usersCount int) []*store.User {
	users := make([]*store.User, usersCount)

	for i := 0; i < usersCount; i++ {
		users[i] = &store.User{
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "12345",
		}
	}

	return users
}

func generatePosts(postsCount int, userId int64) []*store.Post {
	posts := make([]*store.Post, postsCount)

	for i := 0; i < postsCount; i++ {
		posts[i] = &store.Post{
			Content: fmt.Sprintf("Post %d", i),
			Title:   fmt.Sprintf("Post %d", i),
			Tags:    []string{fmt.Sprintf("tag%d", i)},
			UserId:  userId,
		}
	}

	return posts
}

func generateComments(commentsCount int, postId int64, userId int64) []*store.Comment {
	comments := make([]*store.Comment, commentsCount)

	for i := 0; i < commentsCount; i++ {
		comments[i] = &store.Comment{
			Content: fmt.Sprintf("Comment %d", i),
			PostId:  postId,
			UserId:  userId,
		}
	}

	return comments
}
