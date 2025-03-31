package unhasher

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mainHashService/app/entity"
	"mainHashService/app/repo"
	"mainHashService/app/repo/postgres"
	"net/http"
)

func (r *RepoImpl) UnhashData(ctx context.Context, data *postgres.Unhashdata) (entity.UnhashedData, error) {
	jsonData, err := r.marshalHash(&postgres.UnhashRequest{
		HashSalt: data.HashSalt,
		Domain:   data.Domain,
	})
	if err != nil {
		r.lg.Error(err.Error())
		return entity.UnhashedData{}, err
	}
	url := "http://" + r.unhashEndpoint + "/unhash"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		r.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}
	resp, err := r.unhashClient.Do(req)
	if err != nil {
		r.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}

	if resp.StatusCode != http.StatusOK {
		r.lg.Logger.Error().Msgf("status: %d, url: %s", resp.StatusCode, url)
		resp.Body.Close()
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}
	resp.Body.Close()
	unhashedData, err := r.unmarshalHash(body)
	if err != nil {
		r.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}

	return unhashedData, nil
}

func (r *RepoImpl) marshalHash(data *postgres.UnhashRequest) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (r *RepoImpl) unmarshalHash(data []byte) (entity.UnhashedData, error) {
	var resp entity.UnhashedData
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return entity.UnhashedData{}, err
	}
	return resp, nil
}
