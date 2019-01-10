package mailgun

import (
	"context"
	"time"
)

const iso8601date = "2006-01-02"

type Accepted struct {
	Incoming int `json:"incoming"`
	Outgoing int `json:"outgoing"`
	Total    int `json:"total"`
}

type Delivered struct {
	Smtp  int `json:"smtp"`
	Http  int `json:"http"`
	Total int `json:"total"`
}

type Temporary struct {
	Espblock int `json:"espblock"`
}

type Permanent struct {
	SuppressBounce      int `json:"suppress-bounce"`
	SuppressUnsubscribe int `json:"suppress-unsubscribe"`
	SuppressComplaint   int `json:"suppress-complaint"`
	Bounce              int `json:"bounce"`
	DelayedBounce       int `json:"delayed-bounce"`
	Total               int `json:"total"`
}

type Failed struct {
	Temporary Temporary `json:"temporary"`
	Permanent Permanent `json:"permanent"`
}

type Total struct {
	Total int `json:"total"`
}

type Stats struct {
	Time         string    `json:"time"`
	Accepted     Accepted  `json:"accepted"`
	Delivered    Delivered `json:"delivered"`
	Failed       Failed    `json:"failed"`
	Stored       Total     `json:"stored"`
	Opened       Total     `json:"opened"`
	Clicked      Total     `json:"clicked"`
	Unsubscribed Total     `json:"unsubscribed"`
	Complained   Total     `json:"complained"`
}

type statsTotalResponse struct {
	End        string  `json:"end"`
	Resolution string  `json:"resolution"`
	Start      string  `json:"start"`
	Stats      []Stats `json:"stats"`
}

type Resolution string

const (
	ResolutionHour = Resolution("hour")
	ResolutionDay = Resolution("day")
	ResolutionMonth = Resolution("month")
)

type ListStatOptions struct {
	Resolution Resolution
	Duration string
	Start time.Time
	End time.Time
	Limit int
	Skip int
}

// Returns total stats for a given domain for the specified time period
func (mg *MailgunImpl) ListStats(ctx context.Context, events []string, opts *ListStatOptions) ([]Stats, error) {
	// TODO: Test this
	r := newHTTPRequest(generateApiUrl(mg, statsTotalEndpoint))

	if opts != nil {
		if !opts.Start.IsZero() {
			r.addParameter("start", opts.Start.Format(iso8601date))
		}
		if !opts.End.IsZero() {
			r.addParameter("end", opts.End.Format(iso8601date))
		}
		if opts.Resolution != "" {
			r.addParameter("resolution", string(opts.Resolution))
		}
		if opts.Duration != "" {
			r.addParameter("duration", opts.Duration)
		}
	}

	for _, e := range events {
		r.addParameter("event", e)
	}

	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var res statsTotalResponse
	err := getResponseFromJSON(ctx, r, &res)
	if err != nil {
		return nil, err
	} else {
		return res.Stats, nil
	}
}
