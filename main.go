// ssh-echo project main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func cfgWithBanner(ctx ssh.Context) *gossh.ServerConfig {
	c := &gossh.ServerConfig{}
	c.BannerCallback = func(conn gossh.ConnMetadata) string {
		return fmt.Sprintf("---- Password echo server -----\nUser: %s\n\nOpen a shell to see your password\n", conn.User())
	}
	return c
}

func pwdDisplayShell(s ssh.Session) {
	io.WriteString(s, fmt.Sprintf("\"%s\" is your password\n", s.Context().Value("password")))

	tick := time.NewTicker(time.Second * 30)
	select {
	case <-tick.C:
		{
			s.Close()
		}
	case <-s.Context().Done():
		{
			tick.Stop()
			s.Close()
		}
	}
}

func main() {
	fmt.Println("SSH password echo server v1.0")
	fmt.Println("Please note that this is not a authentic SSH server, it only helps you reveal your forgotten password.")
	addr := flag.String("l", ":22", "listen to this address and port")
	flag.Parse()

	srv := &ssh.Server{
		Addr:                 *addr,
		Handler:              pwdDisplayShell,
		ServerConfigCallback: cfgWithBanner,
	}
	srv.SetOption(ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		ctx.SetValue("password", password)
		return true
	}))

	if e := srv.ListenAndServe(); e != nil {
		fmt.Println("Error:", e)
	}
}
