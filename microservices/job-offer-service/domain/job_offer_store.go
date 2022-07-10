package domain

import "context"

type JobOfferStore interface {
	GetRecommendations(ctx context.Context, user *User, jobOffers []*Post) ([]*PostsID, error)
}
