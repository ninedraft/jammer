package counter

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/ninedraft/gemax/gemax"
	"github.com/ninedraft/gemax/gemax/status"
)

type Counter struct {
	StatsPath string

	mu    sync.RWMutex
	start time.Time
	hll   *hyperloglog.Sketch
	pages map[string]uint64
}

func New(statsPath string) *Counter {
	return &Counter{
		start: time.Now(),
		hll:   hyperloglog.New(),
		pages: map[string]uint64{},
	}
}

func (counter *Counter) Middleware(next gemax.Handler) gemax.Handler {
	return func(ctx context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
		if counter.matchStatsPath(req.URL().Path) {
			counter.ServeStats(ctx, rw, req)
			return
		}
		counter.incr(req)
		next(ctx, rw, req)
	}
}

func (counter *Counter) ServeStats(ctx context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
	var values = counter.stats()
	sort.Slice(values.PerPage, func(i, j int) bool {
		return values.PerPage[i].hits > values.PerPage[j].hits
	})

	rw.WriteStatus(status.Success, "text/gemini")
	_, _ = io.WriteString(rw, "# Blog stats\n\n")
	fmt.Fprintf(rw, "- uptime: %v\n", values.Uptime)
	fmt.Fprintf(rw, "- unique clients estimated by hyperloglog: %d\n", values.NClients)

	_, _ = io.WriteString(rw, "## Requests per page\n")
	for _, pp := range values.PerPage {
		fmt.Fprintf(rw, "- %s: %d\n", pp.path, pp.hits)
	}
}

func (counter *Counter) incr(req gemax.IncomingRequest) {
	counter.mu.Lock()
	defer counter.mu.Unlock()

	counter.hll.Insert([]byte(req.RemoteAddr()))
	counter.pages[req.URL().Path]++
}

func (counter *Counter) matchStatsPath(p string) bool {
	return counter.StatsPath != "" &&
		strings.TrimSuffix(p, "/") == strings.TrimSuffix(counter.StatsPath, "/")
}

func (counter *Counter) stats() *statsValues {
	counter.mu.RLock()
	defer counter.mu.RUnlock()

	var values = &statsValues{
		Uptime:   time.Since(counter.start),
		NClients: counter.hll.Estimate(),
	}
	for page, hits := range counter.pages {
		values.PerPage = append(values.PerPage, pageStat{
			path: page,
			hits: hits,
		})
	}
	return values
}

type statsValues struct {
	Uptime   time.Duration
	NClients uint64
	PerPage  []pageStat
}

type pageStat struct {
	path string
	hits uint64
}
