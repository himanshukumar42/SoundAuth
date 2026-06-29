package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	sdk "github.com/himanshukumar42/soundauth/SDK"
	"github.com/himanshukumar42/soundauth/internal/models"
)

const (
	RequestTimeout = 5 * time.Second
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)

	signal.Notify(
		signalCh,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		sig := <-signalCh

		log.Printf("[System] Received signal %v", sig)
		cancel()
	}()

	// SDK
	auth := sdk.NewAuthenticationSDK()

	requests := []models.AuthRequest{
		{
			TenantID:   "google",
			Provider:   models.ProviderPasskey,
			Credential: "credential-1",
			DeviceID:   "DEVICE-101",
		},
		{
			TenantID:   "google",
			Provider:   models.ProviderPasskey,
			Credential: "credential-2",
			DeviceID:   "DEVICE-102",
		},
		{
			TenantID:   "google",
			Provider:   models.ProviderPasskey,
			Credential: "credential-3",
			DeviceID:   "DEVICE-103",
		},
	}

	// Concurrent Authentication

	var wg sync.WaitGroup
	for _, req := range requests {
		wg.Add(1)

		go func(req models.AuthRequest) {
			defer wg.Done()

			reqCtx, cancel := context.WithTimeout(ctx, RequestTimeout)
			defer cancel()

			resp, err := auth.Authenticate(reqCtx, req)
			if err != nil {
				log.Printf("[ERROR] %v\n", err)
				return
			}
			log.Printf("[SUCCESS] User=%s Provider=%s Token=%s\n", resp.UserID, resp.Provider, resp.AccessToken)
		}(req)
	}

	wg.Wait()
	log.Println("[System] Authentication Completed")
}
