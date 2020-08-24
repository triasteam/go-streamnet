package hash

type Sponge interface {
	absorb(trits []byte, offset int, length int)
	squeeze(offset int, length int) []byte
	reset()
}
