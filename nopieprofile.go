package nopieprofile

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	_ "unsafe"

	"github.com/google/pprof/profile"
)

type moduledata struct {
	pcHeader     uintptr
	funcnametab  []byte
	cutab        []uint32
	filetab      []byte
	pctab        []byte
	pclntable    []byte
	ftab         []uintptr
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext uintptr
}

//go:linkname activeModules runtime.activeModules
func activeModules() []*moduledata

// Rewrite rewrites the sample held by the profilePath's profile from a runtime address to a static address.
func Rewrite(profilePath string) error {
	f, err := os.Open(profilePath)
	if err != nil {
		return fmt.Errorf("open profile: %w", err)
	}
	defer f.Close()
	p, err := profile.Parse(f)
	if err != nil {
		return fmt.Errorf("parse profile: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close profile: %w", err)
	}
	text := uint64(activeModules()[0].text)
	for _, s := range p.Sample {
		for _, l := range s.Location {
			if text > l.Address {
				continue
			}
			const runtimeText = 0x100001000
			l.Address = runtimeText + (l.Address - text)
		}
	}
	f2, err := os.Create(profilePath)
	if err != nil {
		return fmt.Errorf("create new profile: %w", err)
	}
	defer f2.Close()
	if err := p.Write(f2); err != nil {
		return fmt.Errorf("write new profile: %w", err)
	}
	if err := f2.Close(); err != nil {
		return fmt.Errorf("close new profile: %w", err)
	}
	return nil
}

// RewriteTestProfile rewrites the CPU profile generated during testing, making sure
// the profile corresponds to static addresses. The profile is located in the directory specified
// by the "test.outputdir" test flag and has the file name specified by the "test.cpuprofile" flag.
//
// This function returns an error if the profile was not successfully rewritten.
func RewriteTestProfile() error {
	dir := flag.Lookup("test.outputdir").Value.String()
	cpuprofile := flag.Lookup("test.cpuprofile").Value.String()
	profilePath := filepath.Join(dir, cpuprofile)
	if profilePath == "" {
		return nil
	}
	if err := Rewrite(profilePath); err != nil {
		return fmt.Errorf("rewrite profile: %w", err)
	}
	return nil
}
