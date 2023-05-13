package nopieprofiletest

import (
	"log"
	"os"
	"testing"

	"github.com/orisano/nopieprofile"
)

func Main(m *testing.M) {
	code := m.Run()
	if err := nopieprofile.RewriteTestProfile(); err != nil {
		log.Printf("warn: failed to rewrite test profile: %v", err)
	}
	os.Exit(code)
}
