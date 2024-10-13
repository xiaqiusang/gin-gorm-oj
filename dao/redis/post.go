package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	//确定查询的索引起始
	start := (page - 1) * size
	end := start + size - 1

	//按分数从大到小查询指定数量的帖子
	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	// 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostScoreZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids { //返回value postID
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	//查找key为分数为1的数量，统计帖子赞成票的数量
	//	v1 := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}

	//使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		//	//查找key为分数为1的数量，统计帖子赞成票的数量
		pipeline.ZCount(key, "1", "1").Val()
	}
	//假设有三个结果5，10，15，cmders为
	//cmders := []redis.Cmder{
	//	&redis.IntCmd{  对应 KeyPostVotedZSetPF1 的 ZCount 结果
	//		val: 5,
	//	},
	//	&redis.IntCmd{  对应 KeyPostVotedZSetPF2 的 ZCount 结果
	//		val: 10,
	//	},
	//	&redis.IntCmd{  对应 KeyPostVotedZSetPF3 的 ZCount 结果
	//		val: 15,
	//	},
	//}
	cmders, err := pipeline.Exec() //返回
	if err != nil {
		return nil, err
	}
	data = make([]int64, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区根据ids查询每篇帖子的赞成票的数据
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderkey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderkey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore把分区的帖子set与帖子分数的zest生成一个新zset
	//对于新的zset按之前的逻辑取数据

	//社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(orderkey).Val() < 1 {
		//不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderkey)
		//使用 ZINTERSTORE 命令将社区集合 (cKey) 和排序集合 (orderkey) 进行交集操作，结果存储在 key 中。这里的交集操作是以最大值聚合 (Aggregate: "MAX")。
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在直接根据key查询ids
	return getIDsFromKey(key, p.Page, p.Size)
}
