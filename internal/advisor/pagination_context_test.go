package advisor

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type cancelingMetadataGetter struct {
	cancel context.CancelFunc
	calls  int
}

func (f *cancelingMetadataGetter) Get(context.Context, string) ([]byte, error) {
	f.calls++
	if f.calls > 6 {
		return nil, fmt.Errorf("fixture runaway tripwire")
	}
	if f.calls == 3 {
		f.cancel()
	}
	return []byte(fmt.Sprintf(`{"value":[],"nextLink":"https://management.azure.com/page?next=%d"}`, f.calls)), nil
}
func TestDistinctMetadataContinuationsHonorCancellationBetweenPages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	getter := &cancelingMetadataGetter{cancel: cancel}
	rows, err := NewMetadataClient(getter, "https://management.azure.com").RecommendationTypes(ctx)
	if !errors.Is(err, context.Canceled) || rows != nil || getter.calls != 3 {
		t.Fatalf("error=%v rows=%v calls=%d", err, rows, getter.calls)
	}
}
