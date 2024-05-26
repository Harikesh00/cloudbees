package main

import (
	"context"
	"log"
	"time"

	"github.com/Harikesh00/cloudbees/blog"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := blog.NewBlogServiceClient(conn)

	// Create a new post
	postID := createPost(client)

	// Read the created post
	readPost(client, postID)

	// Update the created post
	updatePost(client, postID)

	// Read the updated post
	readPost(client, postID)

	// Delete the post
	deletePost(client, postID)

	// Get the post after post deletion
	readPost(client, postID)
}

func createPost(client blog.BlogServiceClient) string {
	log.Println("CREATE POST STARTED")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	post, err := client.CreatePost(ctx, &blog.CreatePostRequest{
		Title:           "First Post",
		Content:         "This is the content of the first post.",
		Author:          "Author Name",
		PublicationDate: time.Now().Format(time.RFC3339),
		Tags:            []string{"tag1", "tag2", "tag3"},
	})

	if err != nil {
		log.Fatalf("could not create post: %v", err)
		return ""
	}

	log.Printf("POST CREATED: %v", post.Post)
	return post.Post.PostId
}

func readPost(client blog.BlogServiceClient, postID string) {
	log.Println("READ POST STARTED")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	post, err := client.ReadPost(ctx, &blog.ReadPostRequest{PostId: postID})
	if err != nil {
		log.Fatalf("could not read post: %v", err)
	}

	if post.Error != "" {
		log.Fatalf("Error while reading post: %s", post.Error)
		return
	}

	log.Printf("POST: %v", post.Post)
}

func updatePost(client blog.BlogServiceClient, postID string) {
	log.Println("UPDATE POST STARTED")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	post, err := client.UpdatePost(ctx, &blog.UpdatePostRequest{
		PostId:  postID,
		Title:   "First Post Updated",
		Content: "Content of this log updated.",
		Author:  "Author Name",
		Tags:    []string{"updateTag1"},
	})

	if err != nil {
		log.Fatalf("could not update post: %v", err)
		return
	}

	if post.Error != "" {
		log.Printf("Error updating post: %v", post.Error)
		return
	}

	log.Printf("POST UPDATED: %v", post.Post)
}

func deletePost(client blog.BlogServiceClient, postID string) {
	log.Println("DELETE POST STARTED")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeletePost(ctx, &blog.DeletePostRequest{PostId: postID})
	if err != nil {
		log.Fatalf("could not delete post: %v", err)
	}

	if res.Error != "" {
		log.Printf("Error while deleting post: %v", res.Error)
		return
	}

	log.Printf("POST DELETED: %v", res.Message)
}
