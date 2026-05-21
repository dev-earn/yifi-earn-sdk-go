package yifi

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

// tokens

func encodeTokensRequest(req models.TokensRequest) url.Values {
	q := url.Values{}

	if req.Search != "" {
		q.Set("search", req.Search)
	}
	if len(req.Chains) > 0 {
		q.Set("chains", strings.Join(req.Chains, ","))
	}
	if len(req.Tags) > 0 {
		q.Set("tags", strings.Join(req.Tags, ","))
	}
	if req.Group != "" {
		q.Set("group", req.Group)
	}
	if req.Limit > 0 {
		q.Set("limit", strconv.FormatInt(req.Limit, 10))
	}
	if req.Offset != nil {
		q.Set("offset", strconv.FormatInt(*req.Offset, 10))
	}

	return q
}

func encodeFromTo(req models.FromTo) url.Values {
	q := url.Values{}

	if !req.From.IsZero() {
		q.Set("from", req.From.UTC().Format(time.RFC3339))
	}
	if !req.To.IsZero() {
		q.Set("to", req.To.UTC().Format(time.RFC3339))
	}

	return q
}

// markets

func encodeMarketsRequest(req models.MarketsRequest) url.Values {
	q := url.Values{}

	if req.Search != "" {
		q.Set("search", req.Search)
	}

	if req.IDs != "" {
		q.Set("ids", req.IDs)
	}

	if len(req.CategoryIDs) > 0 {
		q.Set("categoryIds", strings.Join(req.CategoryIDs, ","))
	}

	if len(req.PlatformIDs) > 0 {
		q.Set("platformIds", strings.Join(req.PlatformIDs, ","))
	}

	if len(req.ChainIDs) > 0 {
		q.Set("chainIds", strings.Join(req.ChainIDs, ","))
	}

	if len(req.TokenIDs) > 0 {
		q.Set("tokenIds", strings.Join(req.TokenIDs, ","))
	}

	if req.MinTVLUsd != nil {
		q.Set("minTvlUsd", req.MinTVLUsd.String())
	}

	if req.MaxTVLUsd != nil {
		q.Set("maxTvlUsd", req.MaxTVLUsd.String())
	}

	if req.MinAPR != nil {
		q.Set("minApr", req.MinAPR.String())
	}

	if req.MaxAPR != nil {
		q.Set("maxApr", req.MaxAPR.String())
	}

	if req.SortBy != "" {
		q.Set("sortBy", req.SortBy)
	}

	if req.SortDirection != "" {
		q.Set("sortDirection", req.SortDirection)
	}

	if req.Limit > 0 {
		q.Set("limit", strconv.FormatInt(req.Limit, 10))
	}

	if req.Offset != nil {
		q.Set("offset", strconv.FormatInt(int64(*req.Offset), 10))
	}

	return q
}

// strategies

func encodeStrategiesRequest(req models.StrategySearch) url.Values {
	q := url.Values{}

	if req.Search != nil && *req.Search != "" {
		q.Set("search", *req.Search)
	}

	if req.IDs != nil && *req.IDs != "" {
		q.Set("ids", *req.IDs)
	}

	if req.Status != nil {
		q.Set("status", string(*req.Status))
	}

	if req.Group != "" {
		q.Set("group", req.Group)
	}

	if len(req.CategoryIDs) > 0 {
		q.Set("category_ids", strings.Join(req.CategoryIDs, ","))
	}

	if len(req.ChainIDs) > 0 {
		q.Set("chain_ids", strings.Join(req.ChainIDs, ","))
	}

	if len(req.TokenIDs) > 0 {
		q.Set("token_ids", strings.Join(req.TokenIDs, ","))
	}

	if len(req.PlatformIDs) > 0 {
		q.Set("platform_ids", strings.Join(req.PlatformIDs, ","))
	}

	if req.MinTVLUsd != nil {
		q.Set("minTvlUsd", req.MinTVLUsd.String())
	}

	if req.MaxTVLUsd != nil {
		q.Set("maxTvlUsd", req.MaxTVLUsd.String())
	}

	if req.MinAPR != nil {
		q.Set("minApr", req.MinAPR.String())
	}

	if req.MaxAPR != nil {
		q.Set("maxApr", req.MaxAPR.String())
	}

	if req.SortBy != "" {
		q.Set("sortBy", req.SortBy)
	}

	if req.SortDirection != "" {
		q.Set("sortDirection", req.SortDirection)
	}

	if req.Limit > 0 {
		q.Set("limit", strconv.FormatInt(req.Limit, 10))
	}

	if req.Offset != nil {
		q.Set("offset", strconv.FormatInt(int64(*req.Offset), 10))
	}

	return q
}
