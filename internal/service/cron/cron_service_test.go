package cron_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"base-code-go-gin-clean/internal/service/cron"
)

func TestCronService(t *testing.T) {
	service := cron.NewCronService()

	t.Run("add and run job", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		// Create a job that signals when it's done
		jobRan := false
		job := func() {
			jobRan = true
			wg.Done()
		}

		// Schedule the job to run every second
		_, err := service.AddJob("* * * * * *", job)
		assert.NoError(t, err)

		service.Start()
		defer service.Stop()

		// Wait for the job to run or timeout after 2 seconds
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Job ran successfully
		case <-time.After(2 * time.Second):
			t.Fatal("Job did not run within expected time")
		}

		assert.True(t, jobRan, "Expected job to have run")
	})

	t.Run("stop service", func(t *testing.T) {
		jobRan := false
		job := func() {
			jobRan = true
		}

		_, err := service.AddJob("* * * * * *", job)
		assert.NoError(t, err)

		service.Start()
		service.Stop()

		// Give some time for the job to potentially run
		time.Sleep(1100 * time.Millisecond)

		assert.False(t, jobRan, "Job should not have run after stop")
	})

	t.Run("invalid schedule", func(t *testing.T) {
		_, err := service.AddJob("invalid-schedule", func() {})
		assert.Error(t, err)
	})
}
