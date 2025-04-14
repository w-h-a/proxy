package fault

import (
	"context"
	"math"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/w-h-a/pkg/telemetry/log"
)

func ShouldApply(ctx context.Context, rules []Rule) bool {
	for _, rule := range rules {
		if shouldApply(ctx, rule) {
			return true
		}
	}

	return false
}

func shouldApply(ctx context.Context, rule Rule) bool {
	log.Infof("matching rule: %+v", rule)

	// http only
	var req *http.Request
	if r, ok := HttpRequestFromCtx(ctx); ok {
		req = r
	}

	if req != nil {
		if rule.Endpoint != "" {
			log.Infof("matching endpoint '%s' against '%s'", rule.Endpoint, req.URL.Path)
			if match, _ := regexp.MatchString(rule.Endpoint, req.URL.Path); !match {
				return false
			}
		}

		if rule.HttpMethod != "" {
			log.Infof("matching http method '%s' against '%s'", rule.HttpMethod, req.Method)
			if match, _ := regexp.MatchString(rule.HttpMethod, req.Method); !match {
				return false
			}
		}
	}

	// grpc

	// tcp, etc

	// all protocols
	if rule.Percentage > 0 {
		randomInteger := rand.Intn(100)
		log.Infof("matching percentage %.2f against draw %d", rule.Percentage, randomInteger)
		if randomInteger > int(math.Min(rule.Percentage, 100)) {
			return false
		}
	}

	return true
}
