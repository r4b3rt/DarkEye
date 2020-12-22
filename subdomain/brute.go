package subdomain

import (
	"github.com/zsdevX/DarkEye/common"
	"golang.org/x/time/rate"
	"strconv"
	"sync"
	"time"
)

var (
	maxBruteWorker = 16
)

func (s *SubDomain) getBrute(query string) {
	pps, _ := strconv.Atoi(s.BruteRate)
	rateLimiter := rate.NewLimiter(rate.Every(1000000*time.Microsecond/time.Duration(pps)), 10)
	target := make([]string, 0)
	length, _ := strconv.Atoi(s.BruteLength)
	i := 0
	for length > 0 {
		target = genSource(target, common.LowCaseAlpha+"0123456789", i+1)
		length--
		i++
	}
	limiter := make(chan int, maxBruteWorker)
	wg := sync.WaitGroup{}
	wg.Add(len(target))
	for _, t := range target {
		if common.ShouldStop(&s.Stop) {
			wg.Done()
			continue
		}
		limiter <- 1
		go func() {
			rateWait(rateLimiter)
			s.try(t, query)
			<-limiter
			wg.Done()
		}()
	}
	wg.Wait()
}

func rateWait(r *rate.Limiter) {
	if r == nil {
		return
	}
	for {
		if r.Allow() {
			break
		} else {
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func genSource(target []string, source string, length int) []string {
	if length == 1 {
		for _, v := range source {
			target = append(target, string(v))
		}
		return target
	}
	t := target
	for _, v1 := range source {
		for _, v2 := range t {
			v := string(v1) + string(v2)
			if len(v) != length {
				continue
			}
			target = append(target, v)
		}
	}
	return target
}
