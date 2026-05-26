package utils

import (
	"context"
	"log"
	"math/rand/v2"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

// JoinUrl joins a base URL with additional path segments safely
func JoinUrl(base string, paths ...string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	all := []string{u.Path}
	all = append(all, paths...)
	u.Path = path.Join(all...)

	return u.String(), nil
}

// Calls JoinUrl, but panics on errors, use with caution.
func JoinUrlOrPanic(base string, paths ...string) string {
	url, err := JoinUrl(base, paths...)
	if err != nil {
		panic(err)
	}
	return url
}

func DerefOrEmpty[T any](val *T) T {
	if val == nil {
		var empty T
		return empty
	}
	return *val
}

func IsNotNil[T any](val *T) bool {
	return val != nil
}

func WithSignalCancel(usecase string) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received shutdown signal, stopping " + usecase + "...")
		cancel()
	}()
	return ctx
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	sb := strings.Builder{}
	for range n {
		sb.WriteByte(letters[rand.IntN(len(letters))])
	}
	return sb.String()
}
