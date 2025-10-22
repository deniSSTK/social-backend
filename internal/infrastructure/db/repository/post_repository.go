package repository

import (
	"context"
	"social-backend/internal/domain/post"
	"social-backend/internal/infrastructure/execer"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	conn *pgxpool.Pool
}

func NewPostRepository(conn *pgxpool.Pool) *PostRepository {
	return &PostRepository{conn}
}

func (r *PostRepository) InsertTx(ctx context.Context, exec execer.Execer, post post.Post) (uuid.UUID, error) {
	id := uuid.New()
	_, err := exec.Exec(ctx, `
		INSERT INTO posts	
		(id, description, author_id, close_friends, pinned)
		VALUES ($1, $2, $3, $4, $5)
	`, id, post.Description, post.AuthorId, post.CloseFriends, post.Pinned)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *PostRepository) GetById(ctx context.Context, postId uuid.UUID) (post.Post, error) {
	var targetPost post.Post
	if err := r.conn.QueryRow(ctx, `
		SELECT 
		    id,
		    description,
		    author_id,
		    created_at,
		    likes_count,
		    comments_count,
		    close_friends,
		    pinned
		FROM posts
		WHERE id = $1
	`, postId).Scan(&targetPost); err != nil {
		return post.Post{}, err
	}

	return targetPost, nil
}

func (r *PostRepository) GetUserPosts(ctx context.Context, userId uuid.UUID, offset int) ([]post.Post, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT
		    p.id,
		    description,
		    author_id,
		    created_at,
			close_friends,
			pinned,
			i.url AS first_image
		FROM posts p
		LEFT JOIN LATERAL (
		    SELECT url
		    FROM images
		    WHERE post_id = p.id
		    ORDER BY position
		    LIMIT 1
		) i ON true
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT 10 OFFSET $2
	`, userId, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var targetPost post.Post
		if err = rows.Scan(
			&targetPost.Id,
			&targetPost.Description,
			&targetPost.AuthorId,
			&targetPost.CreatedAt,
			&targetPost.CloseFriends,
			&targetPost.Pinned,
			&targetPost.FirstImage,
		); err != nil {
			return nil, err
		}

		posts = append(posts, targetPost)
	}

	return posts, nil
}

func (r *PostRepository) InsertHashtagTx(ctx context.Context, exec execer.Execer, hashtag post.Hashtag) error {
	_, err := exec.Exec(ctx, `
		INSERT INTO post_hashtags
		(post_id, hashtag_id, position)
		VALUES ($1, $2, $3)
	`, hashtag.PostId, hashtag.HashtagId, hashtag.Position)
	return err
}
