package repository

import (
	"context"
	"fmt"
	"social-backend/internal/domain/post"
	"social-backend/internal/infrastructure/dto/response"
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

func (r *PostRepository) GetUserPosts(ctx context.Context, userId uuid.UUID, offset int) ([]response.GetPostByUserId, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT
		    p.id,
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

	var posts []response.GetPostByUserId
	for rows.Next() {
		var targetPost response.GetPostByUserId
		if err = rows.Scan(
			&targetPost.Id,
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

func (r *PostRepository) GetPostCountsById(ctx context.Context, postId uuid.UUID) (response.GetPostCountsById, error) {
	var res response.GetPostCountsById
	if err := r.conn.QueryRow(ctx, `
		SELECT likes_count, comments_count
		FROM posts
		WHERE id = $1
	`, postId).Scan(&res.LikesCount, &res.CommentsCount); err != nil {
		return response.GetPostCountsById{}, err
	}

	return res, nil
}

func (r *PostRepository) GetFeedPosts(ctx context.Context, userId uuid.UUID, offset int) ([]response.GetFeedPostByUserId, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT 
		    p.id,
		    p.description,
		    ARRAY (
		    	SELECT ARRAY_AGG(url order by position)
		        FROM images
		    	WHERE post_id = p.id
		    	LIMIT 3
		    ) AS images_urls,
		    (
		        SELECT 
		        COUNT(*)
		        FROM images 
		            WHERE post_id = p.id
		    ) AS images_count,
		    likes_count,
		    comments_count,
		    json_build_object (
		    	'id', u.id,
		    	'username', u.username,
		    	'iconUrl', u.icon_url
		    ) as author,
		    p.created_at,
		    EXISTS (
		        SELECT 1
		        FROM post_likes pl
		        WHERE pl.post_id = p.id AND pl.author_id = $1
		    ) AS if_current_user_liked
		FROM posts p
		JOIN users u ON u.id = p.author_id
		JOIN followings f ON f.follow_to_id = p.author_id
		WHERE f.follower_id = $1
		ORDER BY created_at DESC
		LIMIT 5 OFFSET $2
	`, userId, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []response.GetFeedPostByUserId
	for rows.Next() {
		var targetPost response.GetFeedPostByUserId
		if err = rows.Scan(
			&targetPost.Id,
			&targetPost.Description,
			&targetPost.ImagesUrls,
			&targetPost.ImagesCount,
			&targetPost.LikesCount,
			&targetPost.CommentsCount,
			&targetPost.Author,
			&targetPost.CreatedAt,
			&targetPost.IfCurrentUserLiked,
		); err != nil {
			return nil, err
		}

		posts = append(posts, targetPost)
	}

	return posts, nil
}

func (r *PostRepository) LikePostTx(ctx context.Context, exec execer.Execer, postId, userId uuid.UUID) error {
	fmt.Println("postId"+postId.String(), userId.String())
	_, err := exec.Exec(ctx, `
		INSERT INTO post_likes (post_id, author_id)
		VALUES ($1, $2)
	`, postId, userId)
	return err
}

func (r *PostRepository) RemoveLikePostTx(ctx context.Context, exec execer.Execer, postId, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		DELETE FROM post_likes
		WHERE post_id = $1 AND author_id = $2 
	`, postId, userId)
	return err
}

func (r *PostRepository) UpdatePostLikesCountTx(ctx context.Context, exec execer.Execer, count int, postId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE posts
		SET likes_count = likes_count + $1
		WHERE id = $2
	`, count, postId)
	return err
}
