package sshutil

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

func SCP(c *ssh.Client, src, dest string) error {
	fi, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fi.Close()
	stat, err := fi.Stat()
	if err != nil {
		return err
	}
	s, err := c.NewSession()
	if err != nil {
		return err
	}
	return scp(s, fi, stat.Size(), dest)
}

func scp(s *ssh.Session, r io.Reader, length int64, dest string) error {
	go func() {
		w, err := s.StdinPipe()
		if err != nil {
			return
		}
		defer w.Close()
		fmt.Fprintln(w, "C0644", length, dest)
		_, err = io.Copy(w, r)
		if err != nil {
			return
		}
		fmt.Fprint(w, "\x00")
	}()
	if err := s.Run("/usr/bin/scp -qrt ./"); err != nil {
		return err
	}
	return nil
}
