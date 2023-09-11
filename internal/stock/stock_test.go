package stock

import (
	"github.com/jackylong1/lltrade/pkg/zlog"
	"go.uber.org/zap"
	"testing"
)

func TestSinaStrToEntry(t *testing.T) {
	res, err := Get([]string{"sh601360", "sz000555"})
	if err != nil {
		panic(err)
	}
	list, err := SinaStrToEntry(res)
	if err != nil {
		panic(err)
	}
	zlog.Info("list", zap.Any("list", list))
}

func TestGetStockCode(t *testing.T) {
	//GetStockCodeCount()
	//GetStockCodeList()
	GetStockCodeList("sz_a")
	GetStockCodeList("sh_a")
}
