package scheduling

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	pb "github.com/jukeizu/treediagram/api/scheduling"
	"github.com/shawntoffel/services-core/transport"
)

type httpBinding struct {
	logger  log.Logger
	Service pb.SchedulingServer
}

func NewHttpBinding(logger log.Logger, s pb.SchedulingServer) httpBinding {
	return httpBinding{logger, s}
}

func (h httpBinding) NewServeMux() *http.ServeMux {
	handler := h.makeHandler()

	mux := http.NewServeMux()
	mux.Handle("/create", handler)
	mux.Handle("/jobs", handler)
	mux.Handle("/run", handler)
	mux.Handle("/disable", handler)

	return mux
}

func (h httpBinding) makeHandler() http.Handler {
	createJobHandler := transport.NewDefaultServer(
		h.logger,
		h.createJobEndpoint,
		DecodeCreateJobRequest,
	)

	jobsHandler := transport.NewDefaultServer(
		h.logger,
		h.jobsEndpoint,
		DecodeJobsRequest,
	)

	runJobsHandler := transport.NewDefaultServer(
		h.logger,
		h.runJobsEndpoint,
		DecodeRunJobsRequest,
	)

	disableJobHandler := transport.NewDefaultServer(
		h.logger,
		h.disableJobEndpoint,
		DecodeDisableJobRequest,
	)

	router := mux.NewRouter()

	router.Handle("/create", createJobHandler).Methods("POST")
	router.Handle("/jobs", jobsHandler).Methods("POST")
	router.Handle("/run", runJobsHandler).Methods("POST")
	router.Handle("/disable", disableJobHandler).Methods("POST")

	return router
}

func (h httpBinding) createJobEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.CreateJobRequest)
	return h.Service.Create(ctx, &request)
}

func (h httpBinding) jobsEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.JobsRequest)

	return h.Service.Jobs(ctx, &request)
}

func (h httpBinding) runJobsEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.RunJobsRequest)

	return h.Service.Run(ctx, &request)
}

func (h httpBinding) disableJobEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.DisableJobRequest)

	return h.Service.Disable(ctx, &request)
}

func DecodeCreateJobRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.CreateJobRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}

func DecodeJobsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.JobsRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}

func DecodeRunJobsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.RunJobsRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}

func DecodeDisableJobRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.DisableJobRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}
