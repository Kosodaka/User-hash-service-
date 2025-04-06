package unhasher

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mainHashService/internal/entity"
	"mainHashService/internal/repo"
	"mainHashService/internal/repo/postgres"
	"net/http"
)

func (uc *UnhasherUCImpl) UnhashData(ctx context.Context, data *postgres.Unhashdata) (entity.UnhashedData, error) {
	jsonData, err := uc.marshalHash(&postgres.UnhashRequest{
		HashSalt: data.HashSalt,
		Domain:   data.Domain,
	})
	if err != nil {
		uc.lg.Error(err.Error())
		return entity.UnhashedData{}, err
	}
	url := "http://" + uc.unhashEndpoint + "/unhash"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		uc.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}
	resp, err := uc.unhashClient.Do(req)
	if err != nil {
		uc.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}

	if resp.StatusCode != http.StatusOK {
		uc.lg.Logger.Error().Msgf("status: %d, url: %s", resp.StatusCode, url)
		resp.Body.Close()
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		uc.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}
	resp.Body.Close()
	unhashedData, err := uc.unmarshalHash(body)
	if err != nil {
		uc.lg.Error(err.Error())
		return entity.UnhashedData{}, repo.ErrRepoInternal
	}

	return unhashedData, nil
}

func (uc *UnhasherUCImpl) marshalHash(data *postgres.UnhashRequest) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (uc *UnhasherUCImpl) unmarshalHash(data []byte) (entity.UnhashedData, error) {
	var resp entity.UnhashedData
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return entity.UnhashedData{}, err
	}
	return resp, nil
}
