package Service

import "simpleForum/repository"

type PageInfo struct {
	Topic    *repository.Topic
	PostList []*repository.Topic
}
