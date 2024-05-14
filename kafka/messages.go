package kafka

type KafkaMsg struct {
	Msg string `form:"msg" json:"msg"`
	Val int64  `form:"val" json:"val"`
}
