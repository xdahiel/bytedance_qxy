package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type Topic struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime int64  `json:"createTime"`
}

type Post struct {
	Id         int64  `json:"id"`
	ParentId   int64  `json:"parentId"`
	Content    string `json:"content"`
	CreateTime int64  `json:"createTime"`
}

type TopicDao struct {
}

type PostDao struct {
}

var (
	topicIndexMap map[int64]*Topic
	postIndexMap  map[int64][]*Post

	topicDao  *TopicDao
	topicOnce sync.Once
	postDao   *PostDao
	postOnce  sync.Once
)

func initTopicIndexMap(file string) error {
	open, err := os.Open(file + "topic")
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(open)
	topicTmpMap := make(map[int64]*Topic)
	for sc.Scan() {
		text := sc.Text()
		var topic Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil {
			return err
		}
		topicTmpMap[topic.Id] = &topic
	}
	topicIndexMap = topicTmpMap
	return nil
}

func initPostIndexMap(file string) error {
	open, err := os.Open(file + "post")
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(open)
	postTmpMap := make(map[int64][]*Post)
	for sc.Scan() {
		text := sc.Text()
		var post Post
		if err := json.Unmarshal([]byte(text), &post); err != nil {
			return err
		}
		if post.Id != 0 {
			postTmpMap[post.ParentId] = append(postTmpMap[post.ParentId], &post)
		}
	}
	postIndexMap = postTmpMap
	return nil
}

func (*TopicDao) QueryTopicById(id int64) *Topic {
	return topicIndexMap[id]
}

func (*PostDao) QueryTopicById(id int64) []*Post {
	return postIndexMap[id]
}

func NewTopicDaoInstance() *TopicDao {
	topicOnce.Do(func() {
		topicDao = &TopicDao{}
	})
	return topicDao
}

func NewPostDaoInstance() *PostDao {
	topicOnce.Do(func() {
		postDao = &PostDao{}
	})
	return postDao
}
