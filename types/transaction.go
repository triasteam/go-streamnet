package types

type Transaction struct {
	bytes []byte 
	address Hash 
	trunk Hash
	branch Hash
	obsoleteTag Hash
	value int64
	currentIndex int64
	lastIndex int64
	timestamp int64
	tag Hash
	attachmentTimestamp int64
	attachmentTimestampLowerBound int64
	attachmentTimestampUpperBound int64
	height int64
	sender string
	weightMagnitude int64
	nonce []byte
}
