package domain

import (
	"reflect"
	"testing"
	"time"

	"blog/pkg/ddd"
)

func TestNewComment(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		postID      PostID
		commenterID UserID
		content     string
		want        *Comment
		wantErr     bool
	}{
		{
			name:        "Test Proper Inputs",
			postID:      "1",
			commenterID: "2",
			content:     "3",
			want: &Comment{
				postID:      "1",
				commenterID: "2",
				content:     "3",
			},
			wantErr: false,
		},
		{
			name:        "Test Empty Content Fails",
			postID:      "1",
			commenterID: "2",
			content:     "",
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewComment(tt.postID, tt.commenterID, tt.content)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewComment() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewComment() succeeded unexpectedly")
			}
			if got.postID != tt.want.postID || got.commenterID != tt.want.commenterID ||
				got.content != tt.want.content || got.lastUpdatedAt != nil || got.archivedAt != nil {
				t.Errorf("NewComment() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestComment_SetID(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		postID      PostID
		commenterID UserID
		content     string
		// Named input parameters for target function.
		id CommentID
	}{
		{
			name:        "Test Comment Set ID",
			postID:      "1",
			commenterID: "2",
			content:     "3",
			id:          "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewComment(tt.postID, tt.commenterID, tt.content)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			a.SetID(tt.id)

			if a.GetID() != tt.id {
				t.Errorf("SetID() did not set the comment ID")
			}
		})
	}
}

func TestComment_Edit(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cpostID      PostID
		ccommenterID UserID
		ccontent     string
		// Named input parameters for target function.
		content string
		wantErr bool
	}{
		{
			name:         "Test Successful Comment Edit",
			cpostID:      "1",
			ccommenterID: "2",
			ccontent:     "3",
			content:      "abcd",
			wantErr:      false,
		},
		{
			name:         "Test Failed Comment Edit",
			cpostID:      "1",
			ccommenterID: "2",
			ccontent:     "3",
			content:      "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewComment(tt.cpostID, tt.ccommenterID, tt.ccontent)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := a.Edit(tt.content)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Edit() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Edit() succeeded unexpectedly")
			}

			if a.Content() != tt.content {
				t.Errorf("Edit() did not write the value properly")
			}
		})
	}
}

func TestComment_Archive(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		postID      PostID
		commenterID UserID
		content     string
	}{
		{
			name:        "Test Comment Archive",
			postID:      "1",
			commenterID: "2",
			content:     "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewComment(tt.postID, tt.commenterID, tt.content)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			a.Archive()

			if a.Archived() != true {
				t.Errorf("Archive() failed to archive comment")
			}
		})
	}
}

func TestRebuildComment(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour * 1)

	desiredComment := &Comment{
		AggregateBase: &ddd.AggregateBase{},
		postID:        "2",
		commenterID:   "3",
		content:       "4",
		createdAt:     now,
		lastUpdatedAt: &later,
		archivedAt:    &time.Time{},
	}
	desiredComment.SetID("1")

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id            CommentID
		postID        PostID
		commenterID   UserID
		content       string
		createdAt     time.Time
		lastUpdatedAt *time.Time
		archivedAt    *time.Time
		want          Comment
	}{
		{
			name:          "Test RebuildComment",
			id:            "1",
			postID:        "2",
			commenterID:   "3",
			content:       "4",
			createdAt:     now,
			lastUpdatedAt: &later,
			archivedAt:    &time.Time{},
			want:          *desiredComment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RebuildComment(
				tt.id,
				tt.postID,
				tt.commenterID,
				tt.content,
				tt.createdAt,
				tt.lastUpdatedAt,
				tt.archivedAt,
			)
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("RebuildComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
