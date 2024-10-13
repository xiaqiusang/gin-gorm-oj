package redis

//redis key

// redis key尽量使用命名空间的方式吧v
const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"   //zset:帖子发帖时间
	KeyPostScoreZSet   = "post:score"  //zset:帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" //zset:记录用户及投票类型,参数是Post id
	KeyCommunitySetPF  = "community:"  //set;保存每个分区下帖子的id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
