package blog

import (
	"context"
	"reflect"
	"testing"
)

func TestBlogServer_CreatePost(t *testing.T) {

	type args struct {
		ctx context.Context
		req *CreatePostRequest
	}

	var (
		server = NewBlogServer()

		ctx       = context.Background()
		inputArgs = args{
			ctx: ctx,
			req: &CreatePostRequest{
				Title:           "Test Post",
				Content:         "This is a test post.",
				Author:          "Test Author",
				PublicationDate: "2024-01-01",
				Tags:            []string{"testTag1", "testTag2"},
			},
		}

		want = &CreatePostResponse{
			Post: &Post{
				Title:           "Test Post",
				Content:         "This is a test post.",
				Author:          "Test Author",
				PublicationDate: "2024-01-01",
				Tags:            []string{"testTag1", "testTag2"},
			},
		}
	)

	tests := []struct {
		name    string
		args    args
		want    *CreatePostResponse
		wantErr bool
	}{
		{
			name:    "positive",
			args:    inputArgs,
			want:    want,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := server.CreatePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BlogServer.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// postID is uuid, hence replacing with the what we got
			tt.want.Post.PostId = got.Post.PostId
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlogServer.CreatePost() \ngot= %v, \nwant= %v", got, tt.want)
			}
		})
	}
}

func TestBlogServer_ReadPost(t *testing.T) {
	type args struct {
		ctx    context.Context
		postId string
	}

	ctx := context.Background()
	server := NewBlogServer()

	// Create a post to read
	createReq := &CreatePostRequest{
		Title:           "Test Post",
		Content:         "This is a test post.",
		Author:          "Test Author",
		PublicationDate: "2024-01-01",
		Tags:            []string{"testTag1", "testTag2"},
	}

	createRes, _ := server.CreatePost(ctx, createReq)

	inputArgs := args{
		ctx:    ctx,
		postId: createRes.Post.PostId,
	}

	want := &ReadPostResponse{
		Post: &Post{
			PostId:          createRes.Post.PostId,
			Title:           "Test Post",
			Content:         "This is a test post.",
			Author:          "Test Author",
			PublicationDate: "2024-01-01",
			Tags:            []string{"testTag1", "testTag2"},
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *ReadPostResponse
		wantErr bool
	}{
		{
			name:    "positive",
			args:    inputArgs,
			want:    want,
			wantErr: false,
		},
		{
			name: "negative when post not found",
			args: args{
				ctx:    ctx,
				postId: "test post id",
			},
			want: &ReadPostResponse{
				Error: "Post not found",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.ReadPost(tt.args.ctx, &ReadPostRequest{PostId: tt.args.postId})
			if (err != nil) != tt.wantErr {
				t.Errorf("BlogServer.ReadPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.Post, tt.want.Post) {
				t.Errorf("BlogServer.ReadPost() \ngot= %v, \nwant= %v", got, tt.want)
			}
		})
	}
}

func TestBlogServer_UpdatePost(t *testing.T) {
	type args struct {
		ctx context.Context
		req *UpdatePostRequest
	}

	var (
		ctx    = context.Background()
		server = NewBlogServer()
	)

	// Create a post to update
	createReq := &CreatePostRequest{
		Title:           "Test Post",
		Content:         "This is a test post.",
		Author:          "Test Author",
		PublicationDate: "2024-01-01",
		Tags:            []string{"testTag1", "testTag2"},
	}
	createRes, _ := server.CreatePost(ctx, createReq)

	inputArgs := args{
		ctx: ctx,
		req: &UpdatePostRequest{
			PostId:  createRes.Post.PostId,
			Title:   "Updated Test Post",
			Content: "This is an updated test post.",
			Author:  "Updated Test Author",
			Tags:    []string{"updatedTestTag1", "updatedTestTag2"},
		},
	}

	want := &UpdatePostResponse{
		Post: &Post{
			PostId:          createRes.Post.PostId,
			Title:           "Updated Test Post",
			Content:         "This is an updated test post.",
			Author:          "Updated Test Author",
			PublicationDate: createRes.Post.PublicationDate,
			Tags:            []string{"updatedTestTag1", "updatedTestTag2"},
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *UpdatePostResponse
		wantErr bool
	}{
		{
			name:    "positive",
			args:    inputArgs,
			want:    want,
			wantErr: false,
		},
		{
			name: "negative when post not found",
			args: args{
				ctx: ctx,
				req: &UpdatePostRequest{
					PostId: "test post id",
				},
			},
			want:    &UpdatePostResponse{Error: "Post not found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.UpdatePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BlogServer.UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.Post, tt.want.Post) {
				t.Errorf("BlogServer.UpdatePost() \ngot= %v, \nwant= %v", got, tt.want)
			}
		})
	}
}

func TestBlogServer_DeletePost(t *testing.T) {
	type args struct {
		ctx    context.Context
		postId string
	}

	var (
		ctx    = context.Background()
		server = NewBlogServer()
	)

	// Create a post to delete
	createReq := &CreatePostRequest{
		Title:           "Test Post",
		Content:         "This is a test post.",
		Author:          "Test Author",
		PublicationDate: "2024-01-01",
		Tags:            []string{"testTag1", "testTag2"},
	}
	createRes, _ := server.CreatePost(ctx, createReq)

	inputArgs := args{
		ctx:    ctx,
		postId: createRes.Post.PostId,
	}

	want := &DeletePostResponse{
		Message: "Post deleted",
	}

	tests := []struct {
		name    string
		args    args
		want    *DeletePostResponse
		wantErr bool
	}{
		{
			name:    "positive",
			args:    inputArgs,
			want:    want,
			wantErr: false,
		},
		{
			name: "negative when post not found",
			args: args{
				ctx:    ctx,
				postId: "test post id",
			},
			want:    &DeletePostResponse{Error: "Post not found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.DeletePost(tt.args.ctx, &DeletePostRequest{PostId: tt.args.postId})
			if (err != nil) != tt.wantErr {
				t.Errorf("BlogServer.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlogServer.DeletePost() \ngot= %+v, \nwant= %+v", got, tt.want)
			}
		})
	}
}
