package types

import "fmt"

type StoreData struct {
	Attester string
	Attestee string
	Score    string
}

func (d StoreData) String() string {
	return fmt.Sprintf("Attester: %s, Attestee: %s, Score: %s", d.Attester, d.Attestee, d.Score)
}
