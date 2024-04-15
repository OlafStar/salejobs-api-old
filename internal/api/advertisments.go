package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/olafstar/salejobs-api/internal/types"
)

func (s *APIServer) handleAdvertisments(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.getAdvertisements(w, r)
	}
	if r.Method == "POST"{
		return s.createAdvertisements(w, r)
	}

	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

func (s *APIServer) getAdvertisements(w http.ResponseWriter, r *http.Request) error {
	defaultParams := types.GetAdvertismentBody{
			Page:  1,
			Limit: 10,
	}

	var params types.GetAdvertismentBody
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil && err != io.EOF {
			params = defaultParams
	}

	if params.Page < 1 {
			params.Page = defaultParams.Page
	}
	if params.Limit < 1 || params.Limit > 100 {
			params.Limit = defaultParams.Limit
	}

	totalAds, err := s.store.CountAdvertisements()
	if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to fetch total advertisement count"}
	}

	advertisements, err := s.store.QueryAdvertisements(params.Page, params.Limit)
	if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to fetch advertisements"}
	}

	lastPage := (totalAds + int64(params.Limit) - 1) / int64(params.Limit) // Calculating the number of the last page

	response := types.GetAdvertismentResponse{
			CurrentPage: int64(params.Page),
			Total:       totalAds,
			Last:        lastPage,
			Advertisments: advertisements,
	}

	return WriteJSON(w, http.StatusOK, response)
}

func (s *APIServer) createAdvertisements(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var body types.CreateAdvertisementBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Invalid request body"}
	}
	defer r.Body.Close()

	if err := validateAdvertismentBody(body); err != nil {
		return err 
	}

	err := s.store.CreateAdvertisement(body)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, body)
}