package models

import (
	"sort"
)

type Ranking struct {
	Topics []Topic
}

// ランキング作成
func (r *Ranking) AddRanking(topic Topic) {
	i := sort.Search(len(r.Topics), func(i int) bool {
		return r.Topics[i].Likes < topic.Likes ||
			(r.Topics[i].Likes == topic.Likes && r.Topics[i].ID > topic.ID)
	})

	r.Topics = append(r.Topics, topic)
	copy(r.Topics[i+1:], r.Topics[i:])
	r.Topics[i] = topic

}
