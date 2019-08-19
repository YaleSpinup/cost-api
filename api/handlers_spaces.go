package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// SpaceGetHandler gets the cost for a space, grouped by the service.  By default,
// it pulls data from the start of the month until now.
func (s *server) SpaceGetHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	ceService, ok := s.costExplorerServices[account]
	if !ok {
		msg := fmt.Sprintf("cost explorer service not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cost explorer service %+v", ceService)

	spaceID := vars["space"]
	log.Debugf("getting costs for space %s", spaceID)

	y, m, d := time.Now().Date()

	// if it's the first day of the month, get todays usage thus far
	if d == 1 {
		d = 2
	}
	input := costexplorer.GetCostAndUsageInput{
		Filter: &costexplorer.Expression{
			And: []*costexplorer.Expression{
				&costexplorer.Expression{
					Tags: &costexplorer.TagValues{
						Key: aws.String("spinup:spaceid"),
						Values: []*string{
							aws.String(spaceID),
						},
					},
				},
				&costexplorer.Expression{
					Or: []*costexplorer.Expression{
						&costexplorer.Expression{
							Tags: &costexplorer.TagValues{
								Key: aws.String("yale:org"),
								Values: []*string{
									aws.String(Org),
								},
							},
						},
						&costexplorer.Expression{
							Tags: &costexplorer.TagValues{
								Key: aws.String("spinup:org"),
								Values: []*string{
									aws.String(Org),
								},
							},
						},
					},
				},
				&costexplorer.Expression{
					Not: &costexplorer.Expression{
						Or: []*costexplorer.Expression{
							&costexplorer.Expression{
								Tags: &costexplorer.TagValues{
									Key: aws.String("yale:subsidized"),
									Values: []*string{
										aws.String("true"),
									},
								},
							},
							&costexplorer.Expression{
								Tags: &costexplorer.TagValues{
									Key: aws.String("spinup:subsidized"),
									Values: []*string{
										aws.String("true"),
									},
								},
							},
						},
					},
				},
			},
		},
		Granularity: aws.String("MONTHLY"),
		Metrics: []*string{
			aws.String("BLENDED_COST"),
			aws.String("UNBLENDED_COST"),
			aws.String("USAGE_QUANTITY"),
		},
		TimePeriod: &costexplorer.DateInterval{
			End:   aws.String(fmt.Sprintf("%d-%02d-%02d", y, m, d)),
			Start: aws.String(fmt.Sprintf("%d-%02d-01", y, m)),
		},
	}

	out, err := ceService.GetCostAndUsage(r.Context(), &input)
	if err != nil {
		msg := fmt.Sprintf("failed to get costs for space %s: %s", spaceID, err.Error())
		handleError(w, errors.Wrap(err, msg))
		return
	}

	j, err := json.Marshal(out)
	if err != nil {
		log.Errorf("cannot marshal response (%v) into JSON: %s", out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
