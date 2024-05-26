package blog

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type BlogServer struct {
	UnimplementedBlogServiceServer
	posts map[string]*Post
	mu    sync.Mutex
}

func NewBlogServer() *BlogServer {
	return &BlogServer{
		posts: make(map[string]*Post),
	}
}

func (s *BlogServer) CreatePost(ctx context.Context, req *CreatePostRequest) (*CreatePostResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	post := &Post{
		PostId:          id,
		Title:           req.Title,
		Content:         req.Content,
		Author:          req.Author,
		PublicationDate: req.PublicationDate,
		Tags:            req.Tags,
	}

	s.posts[id] = post
	return &CreatePostResponse{Post: post}, nil
}

func (s *BlogServer) ReadPost(ctx context.Context, req *ReadPostRequest) (*ReadPostResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[req.PostId]
	if !exists {
		return &ReadPostResponse{Error: "Post not found"}, errors.New("Post not found")
	}

	return &ReadPostResponse{Post: post}, nil
}

func (s *BlogServer) UpdatePost(ctx context.Context, req *UpdatePostRequest) (*UpdatePostResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[req.PostId]
	if !exists {
		return &UpdatePostResponse{Error: "Post not found"}, errors.New("Post not found")
	}

	post.Title = req.Title
	post.Content = req.Content
	post.Author = req.Author
	post.Tags = req.Tags

	return &UpdatePostResponse{Post: post}, nil
}

func (s *BlogServer) DeletePost(ctx context.Context, req *DeletePostRequest) (*DeletePostResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.posts[req.PostId]
	if !exists {
		return &DeletePostResponse{Error: "Post not found"}, errors.New("Post not found")
	}

	delete(s.posts, req.PostId)
	return &DeletePostResponse{Message: "Post deleted"}, nil
}
