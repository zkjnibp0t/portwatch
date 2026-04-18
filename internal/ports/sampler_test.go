package ports

import (
	"testing"
	"time"
)

func makeSampler(window time.Duration) *Sampler {
	s := NewSampler(window)
	return s
}

func TestSamplerLatestNilWhenEmpty(t *testing.T) {
	s := makeSampler(time.Minute)
	if s.Latest() != nil {
		t.Fatal("expected nil for empty sampler")
	}
}

func TestSamplerAddAndLatest(t *testing.T) {
	s := makeSampler(time.Minute)
	ports := map[int]struct{}{80: {}, 443: {}}
	s.Add(ports)
	l := s.Latest()
	if l == nil {
		t.Fatal("expected non-nil latest")
	}
	if _, ok := l.Ports[80]; !ok {
		t.Error("expected port 80 in latest sample")
	}
}

func TestSamplerPrunesOldSamples(t *testing.T) {
	s := makeSampler(2 * time.Second)
	now := time.Now()
	base := now.Add(-5 * time.Second)

	// inject old sample manually
	s.samples = append(s.samples, Sample{Timestamp: base, Ports: map[int]struct{}{22: {}}})

	// add a fresh sample which triggers prune
	s.clock = func() time.Time { return now }
	s.Add(map[int]struct{}{80: {}})

	if s.Len() != 1 {
		t.Fatalf("expected 1 sample after pruning, got %d", s.Len())
	}
	if _, ok := s.Latest().Ports[80]; !ok {
		t.Error("expected port 80 in remaining sample")
	}
}

func TestSamplerAllReturnsCopy(t *testing.T) {
	s := makeSampler(time.Minute)
	s.Add(map[int]struct{}{8080: {}})
	s.Add(map[int]struct{}{9090: {}})
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 samples, got %d", len(all))
	}
}

func TestSamplerLenMatchesAdded(t *testing.T) {
	s := makeSampler(time.Minute)
	for i := 0; i < 5; i++ {
		s.Add(map[int]struct{}{i: {}})
	}
	if s.Len() != 5 {
		t.Fatalf("expected 5, got %d", s.Len())
	}
}
