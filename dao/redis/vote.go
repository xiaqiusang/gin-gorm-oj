package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"time"
)

//投票功能:
//投票几种情况
/*
direction=1时，两种情况：
1、之前没有投过票，现在投赞成票 差值绝对值：1 +432
2、之前投反对票，现在投赞成票		2 +432*2
direction=0，两种情况：
1、之前投赞成票，现在取消投票		1	-432
2、之前投反对票，现在取消投票		1  +432
direction=-1，两种情况：
1、之前没投过票，现在投反对票		1	-432
2、之前投赞成票，现在投反对票		2	-432*2

投票限制：
每个帖子自发表之日起一个星期内允许用户投票
1.到期后将redis中保存的赞成和反对票数保存到mysql中
2.到期之后删除保存的keyPostVotedZset

*/

const (
	oneWeekInSecondes = 7 * 24 * 3600
	scorePerVote      = 432 //每一票的分数
)

var (
	ErrVoteTimeExpired = errors.New("投票时间已过")
	ErrVoteRepeated    = errors.New("不允许重复投票")
)

func CreatPost(postID int64) error {
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis获取贴子发布时间
	//获取postID 在有序集合 KeyPostTimeZSet 中的分数
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSecondes {
		return ErrVoteTimeExpired
	}
	//2.更新分数 2和3要一起实现
	//先查当前用户给前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	var op float64
	//如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value)                                                  // 计算两次投票的差值
	pipeline := client.TxPipeline()                                               // 事务操作
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID) // 更新分数
	//3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID).Result()
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
