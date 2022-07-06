package domain

type JobOfferStore interface {
	GetRecommendations(user *User, jobOffers []*Post) ([]*Post, error)
}
