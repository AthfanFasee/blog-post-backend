package data

import (
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/dto"
)

var mockPost = Post{
	ID:        1,
	CreatedAt: time.Now(),
	Title:     "Mocked Post Title",
	PostText:  "Mocked Post PostText",
	Img:       "Mocked Post Img",
	ReadTime:  1,
	LikedBy:   []int64{1, 2},
	CreatedBy: 1,
	Version:   1,
}

var mockMetadata = Metadata{
	CurrentPage:  1,
	PageSize:     1,
	FirstPage:    1,
	LastPage:     1,
	TotalRecords: 1,
}

var mockPostResponseBody = &dto.PostResponseBody{
	ID:        mockPost.ID,
	Title:     mockPost.Title,
	PostText:  mockPost.PostText,
	Img:       mockPost.Img,
	ReadTime:  mockPost.ReadTime,
	LikedBy:   mockPost.LikedBy,
	CreatedBy: mockComment.CreatedBy,
	UserName:  "Mocked User",
}

var mockPostResponseBodyDifferentTitle = &dto.PostResponseBody{
	ID:        mockPost.ID,
	Title:     "Title",
	PostText:  mockPost.PostText,
	Img:       mockPost.Img,
	ReadTime:  mockPost.ReadTime,
	LikedBy:   mockPost.LikedBy,
	CreatedBy: mockComment.CreatedBy,
	UserName:  "Mocked User",
}

type MockPostModel struct{}

func (p MockPostModel) GetAll(title string, filters Filters) ([]*dto.PostResponseBody, Metadata, error) {
	switch {
	case title == "title":
		return []*dto.PostResponseBody{mockPostResponseBodyDifferentTitle}, mockMetadata, nil
	case title != "invalid" && filters.ID == 1:
		return []*dto.PostResponseBody{mockPostResponseBody}, mockMetadata, nil
	case title == "":
		return []*dto.PostResponseBody{mockPostResponseBody}, mockMetadata, nil
	default:
		return nil, Metadata{}, ErrRecordNotFound
	}
}

func (p MockPostModel) Get(id int64) (*Post, error) {
	switch id {
	case 1:
		return &mockPost, nil
	default:
		return nil, ErrRecordNotFound
	}
}

func (p MockPostModel) GetWithUserName(id int64) (*Post, *string, error) {
	switch id {
	case 1:
		return &mockPost, &mockPostResponseBody.UserName, nil
	default:
		return nil, nil, ErrRecordNotFound
	}
}

func (p MockPostModel) Insert(post *Post) error {
	return nil
}

// Consider testing race condition (the errEditConflict error case) in future
func (p MockPostModel) Update(post *Post) error {
	return nil
}

func (p MockPostModel) Delete(id int64) error {
	switch id {
	case 1:
		return nil
	default:
		return ErrRecordNotFound
	}
}

func (p MockPostModel) AddLike(post *Post, userID int64) error {
	return nil
}

func (p MockPostModel) RemoveLike(post *Post, userID int64) error {
	return nil
}
