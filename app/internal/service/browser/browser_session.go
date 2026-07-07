package browser

import (
	"example-wikipedia-scraper/internal/logger"
	types "example-wikipedia-scraper/internal/types/browser"
	"example-wikipedia-scraper/pkg/helpers"
	"context"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type BrowserSession struct {
	options         types.FetchOptions
	mu              sync.RWMutex
	body            []byte
	startTime       time.Time
	endTime         time.Time
	requestID       network.RequestID
	userAgent       string
	context         context.Context
	url             string
	cancel          context.CancelFunc
	networkResponse *network.Response
	headers         map[string]string
	saveBody        bool
	saveHeaders     bool
	allowAssets     bool
	bodyGetSleep    time.Duration
}

func NewBrowserSession(ctx context.Context, cancel context.CancelFunc) *BrowserSession {
	return &BrowserSession{
		context: ctx,
		cancel:  cancel,
		headers: make(map[string]string),
	}
}

func (s *BrowserSession) GetContext() context.Context {
	return s.context
}

func (s *BrowserSession) ResetContext(ctx context.Context, cancel context.CancelFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.context = ctx
	s.cancel = cancel
}

func (s *BrowserSession) GetCancelFunc() context.CancelFunc {
	return s.cancel
}

func (s *BrowserSession) SetURL(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.url = url
}

func (s *BrowserSession) GetURL() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.url
}

func (s *BrowserSession) SetRequestID(id network.RequestID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestID = id
}

func (s *BrowserSession) GetRequestID() network.RequestID {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.requestID
}

func (s *BrowserSession) SetStartTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.startTime = t
}

func (s *BrowserSession) SetEndTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endTime = t
}

func (s *BrowserSession) SetBody(body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.body = body
}

func (s *BrowserSession) SetUserAgent(ua string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userAgent = ua
}

func (s *BrowserSession) SetOptions(options types.FetchOptions) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.options = options
	s.saveBody = options.SaveBody
	s.saveHeaders = options.SaveHeaders
	s.allowAssets = options.AllowAssets
}

func (s *BrowserSession) GetOptions() types.FetchOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.options
}

func (s *BrowserSession) ShouldSaveBody() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.saveBody
}

func (s *BrowserSession) ShouldSaveHeaders() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.saveHeaders
}

func (s *BrowserSession) ShouldAllowAssets() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.allowAssets
}

func (s *BrowserSession) SetNetworkResponse(response *network.Response) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.networkResponse = response
}

func (s *BrowserSession) GetBrowserResponse() types.BrowserResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.networkResponse == nil {
		return types.BrowserResponse{}
	}
	response := types.BrowserResponse{
		StatusCode: int(s.networkResponse.Status),
		URL:        s.url,
		Response:   s.networkResponse,
		UserAgent:  s.userAgent,
		Body:       string(s.body),
		Headers:    make(map[string]string),
	}
	if s.saveHeaders {
		s.extractHeaders(&response)
	}
	if !s.startTime.IsZero() && !s.endTime.IsZero() {
		response.ResponseTime = s.endTime.Sub(s.startTime).Milliseconds()
	}
	return response
}

func (s *BrowserSession) GetBodyGetSleep() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bodyGetSleep
}

func (s *BrowserSession) SetBodyGetSleep(sleep time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bodyGetSleep = sleep
}

func (s *BrowserSession) extractHeaders(response *types.BrowserResponse) {
	if s.networkResponse != nil && s.networkResponse.Headers != nil {
		for k, v := range s.networkResponse.Headers {
			if str, ok := v.(string); ok {
				response.Headers[k] = str
			}
		}
	}
}

func (s *BrowserSession) SetupResponseListener(wg *sync.WaitGroup) {
	ctx := s.GetContext()
	var requestID network.RequestID
	getBrowserResponseOnce := false
	s.endTime = time.Time{}
	s.body = nil
	s.networkResponse = nil
	s.headers = make(map[string]string)
	s.requestID = ""
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			if helpers.UrlsEqualNormalized(e.Request.URL, s.GetURL()) {
				if !getBrowserResponseOnce {
					s.SetRequestID(e.RequestID)
					requestID = e.RequestID
					s.SetStartTime(time.Now())
					getBrowserResponseOnce = true
				}
			}
		case *network.EventLoadingFinished:
			if requestID == e.RequestID {
				s.SetEndTime(time.Now())
				if s.ShouldSaveBody() {
					go func() {
						defer wg.Done()
						time.Sleep(s.GetBodyGetSleep())
						body := s.getBody()
						s.SetBody(body)
					}()
				} else {
					wg.Done()
				}
			}
		case *fetch.EventRequestPaused:
			go s.blockedResources(e)
		}
	})
}

func (s *BrowserSession) blockedResources(event *fetch.EventRequestPaused) {
	c := chromedp.FromContext(s.GetContext())
	ctx := cdp.WithExecutor(s.GetContext(), c.Target)
	allowAssets := s.ShouldAllowAssets()
	if !allowAssets {
		if BlockedResources[event.ResourceType] {
			fetch.FailRequest(event.RequestID, network.ErrorReasonBlockedByClient).Do(ctx)
			return
		}
	}
	fetch.ContinueRequest(event.RequestID).Do(ctx)
}

func (s *BrowserSession) getBody() []byte {
	var html string
	err := chromedp.Run(s.GetContext(), chromedp.Evaluate(`document.documentElement.outerHTML`, &html))
	if err != nil {
		logger.Error("error fetching document HTML", "err", err)
		return nil
	}
	return []byte(html)
}
