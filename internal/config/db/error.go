package db

import "errors"

var (
	errScanRow   = errors.New("error scan row")
	errNoData    = errors.New("error no data")
	errUnmarshal = errors.New("error unmarshal")

	errHostUndefined   = errors.New("host undefined")
	errPortUndefined   = errors.New("port undefined")
	errDBNameUndefined = errors.New("db name undefined")
	errUserUndefined   = errors.New("user undefined")
)
