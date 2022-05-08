package server

import (
	"testing"
	"time"

	"go.uber.org/goleak"
	"github.com/stretchr/testify/require"
)

func TestShutdown(t *testing.T) {
	defer goleak.VerifyNone(t)

	var (
		s = NewServer(DefaultSockFilePath)
	)

	ln, err := s.Listen()
	require.NoError(t, err)

	go func() {
		require.NoError(t, s.Serve(ln))
	}()

	time.Sleep(1 * time.Second)


	require.NoError(t, s.Shutdown(ln))
}

// func BenchmarkApply(b *testing.B) {
// 	var (
// 		now = time.Now()
// 		name = "name"
// 		email = "mail"
// 		role = Admin
// 	)
// 	for n := 0; n < b.N; n++ {
// 		user, err := NewBaseUser(name, email, role, now, now)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		email := user.Email()
// 		user.SetName(email)
// 	}
// }
