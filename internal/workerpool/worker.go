package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type VerificationJob struct {
	Name string
	Run  func(context.Context) error
}

type VerificationResult struct {
	Name  string
	Error error
}

type VerificationPool struct {
	Workers int
	Sem     chan struct{}
}

func NewVerificationPool(workers int, maxConcurrent int) *VerificationPool {
	return &VerificationPool{
		Workers: workers,
		Sem:     make(chan struct{}, maxConcurrent),
	}
}

func (vp *VerificationPool) Verify(ctx context.Context, jobs []VerificationJob) error {
	jobCh := make(chan VerificationJob)
	resultCh := make(chan VerificationResult)

	var wg sync.WaitGroup

	for i := 1; i <= vp.Workers; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobCh:
					if !ok {
						return
					}

					vp.Sem <- struct{}{}
					err := job.Run(ctx)

					<-vp.Sem

					resultCh <- VerificationResult{
						Name:  job.Name,
						Error: err,
					}
				}
			}
		}(i)
	}

	go func() {
		defer close(jobCh)

		for _, job := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobCh <- job:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		if result.Error != nil {
			return fmt.Errorf("%s verification failed: %w", result.Name, result.Error)
		}

		log.Printf("[Verifier] %s ok", result.Name)
	}

	return nil
}
