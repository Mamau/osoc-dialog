package healthcheck

import (
	"context"
	"sync"
)

const probingParallelism = 5

func CheckProbes(ctx context.Context, probes []ProbeFunc) ([]ProbeStatus, error) {
	if len(probes) == 0 {
		return nil, nil
	}

	probesCh := make(chan ProbeFunc, len(probes))
	for _, v := range probes {
		probesCh <- v
	}
	close(probesCh)

	result := make([]ProbeStatus, 0, len(probes))
	resultCh := make(chan ProbeStatus, len(probes))

	parallelism := probingParallelism
	if parallelism > len(probes) {
		parallelism = len(probes)
	}

	var wg sync.WaitGroup
	wg.Add(parallelism)

	for i := 0; i < parallelism; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-probesCh:
					if !ok {
						return
					}
					resultCh <- v(ctx)
				}
			}
		}()
	}

	wg.Wait()
	close(resultCh)

	for v := range resultCh {
		result = append(result, v)
	}

	return result, ctx.Err()
}
