package ports

import "sort"

// ClusterEntry holds a group of ports that frequently appear together.
type ClusterEntry struct {
	Ports []int
	Count int
}

// ClusterDetector groups ports that co-occur across scan diffs.
type ClusterDetector struct {
	minSupport int
	counts     map[string]int
	keys       map[string][]int
}

// NewClusterDetector returns a ClusterDetector that reports groups seen at
// least minSupport times.
func NewClusterDetector(minSupport int) *ClusterDetector {
	if minSupport < 1 {
		minSupport = 1
	}
	return &ClusterDetector{
		minSupport: minSupport,
		counts:     make(map[string]int),
		keys:       make(map[string][]int),
	}
}

// Record registers a set of ports as a co-occurrence event.
func (c *ClusterDetector) Record(ports []int) {
	if len(ports) < 2 {
		return
	}
	copy := append([]int(nil), ports...)
	sort.Ints(copy)
	k := clusterKey(copy)
	if _, ok := c.keys[k]; !ok {
		c.keys[k] = copy
	}
	c.counts[k]++
}

// Clusters returns all port groups that have reached minSupport, sorted by
// count descending.
func (c *ClusterDetector) Clusters() []ClusterEntry {
	var out []ClusterEntry
	for k, cnt := range c.counts {
		if cnt >= c.minSupport {
			out = append(out, ClusterEntry{Ports: c.keys[k], Count: cnt})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Count > out[j].Count
	})
	return out
}

// Reset clears all recorded co-occurrence data.
func (c *ClusterDetector) Reset() {
	c.counts = make(map[string]int)
	c.keys = make(map[string][]int)
}

func clusterKey(sorted []int) string {
	b := make([]byte, 0, len(sorted)*5)
	for i, p := range sorted {
		if i > 0 {
			b = append(b, ',')
		}
		b = appendInt(b, p)
	}
	return string(b)
}

func appendInt(b []byte, n int) []byte {
	if n == 0 {
		return append(b, '0')
	}
	var tmp [10]byte
	i := 10
	for n > 0 {
		i--
		tmp[i] = byte('0' + n%10)
		n /= 10
	}
	return append(b, tmp[i:]...)
}
