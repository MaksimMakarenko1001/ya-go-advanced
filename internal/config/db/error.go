package db

import "errors"

var (
	errScanRow   = errors.New("error scan row")
	errNoData    = errors.New("error no data")
	errUnmarshal = errors.New("error unmarshal")
)
