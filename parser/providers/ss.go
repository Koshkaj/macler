package providers

import (
	"context"
	"gitlab.com/koshkaj/macler/parser/types"
)

func parseSS() {}

func RunParseSS(ctx context.Context, seenProductIDs map[string]struct{}, ch chan<- types.DataNormalized) {
	finished := make(chan struct{})
	for i := 1; i < 5; i++ {
		go func(i int) {
			ch <- []interface{}{"wrote from", i}
			finished <- struct{}{}
		}(i)
	}
	<-finished
}
