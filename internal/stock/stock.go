package stock

import (
	"fmt"
	"github.com/jackylong1/lltrade/pkg/zlog"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

const (
	elemSeparator = ";"
	infoSeparator = "="
	stockLen      = 8
	dateIndex     = 30
	timeIndex     = 31
	dateTimeKey   = "date_time"
)

var (
	sinaStockJsonKey = map[int]string{
		0:  "name",
		1:  "today_open_price",
		2:  "yesterday_close_price",
		3:  "current_price",
		4:  "today_highest_price",
		5:  "today_lowest_price",
		6:  "buy_price_1",
		7:  "sell_price_1",
		8:  "deal_num",
		9:  "deal_amount",
		10: "buy_1_num",
		11: "buy_1_price",
		12: "buy_2_num",
		13: "buy_2_price",
		14: "buy_3_num",
		15: "buy_3_price",
		16: "buy_4_num",
		17: "buy_4_price",
		18: "buy_5_num",
		19: "buy_5_price",
		20: "sell_1_num",
		21: "sell_1_price",
		22: "sell_2_num",
		23: "sell_2_price",
		24: "sell_3_num",
		25: "sell_3_price",
		26: "sell_4_num",
		27: "sell_4_price",
		28: "sell_5_num",
		29: "sell_5_price",
		30: "date",
		31: "time",
	}

	sinaStockIndexKeyEx = map[int]string{
		0:  "Name",
		1:  "TodayOpenPrice",
		2:  "YesterdayClosePrice",
		3:  "CurrentPrice",
		4:  "TodayHighestPrice",
		5:  "TodayLowestPrice",
		6:  "BuyPrice1",
		7:  "SellPrice1",
		8:  "DealNum",
		9:  "DealAmount",
		10: "Buy1Num",
		11: "Buy1Price",
		12: "Sell2Price",
		13: "Sell3Num",
		14: "Sell3Price",
		15: "Sell4Num",
		16: "Sell4Price",
		17: "Sell5Num",
		18: "Buy2Num",
		19: "Buy2Price",
		20: "Buy3Num",
		21: "Buy3Price",
		22: "Buy4Num",
		23: "Buy4Price",
		24: "Buy5Num",
		25: "Buy5Price",
		26: "Sell1Num",
		27: "Sell1Price",
		28: "Sell2Num",
		29: "Sell5Price",
		30: "Date",
		31: "Time",
	}
)

type Entry struct {
	Market              string  `json:"market,omitempty"`
	Code                string  `json:"code"`
	Name                string  `json:"name,omitempty"`
	TodayOpenPrice      float64 `json:"today_open_price,omitempty"`
	YesterdayClosePrice float64 `json:"yesterday_close_price,omitempty"`
	CurrentPrice        float64 `json:"current_price,omitempty"`
	TodayHighestPrice   float64 `json:"today_highest_price,omitempty"`
	TodayLowestPrice    float64 `json:"today_lowest_price,omitempty"`
	BuyPrice1           float64 `json:"buy_price_1,omitempty"`  // 买竞价 即买一报价
	SellPrice1          float64 `json:"sell_price_1,omitempty"` // 卖竞价 即卖一报价
	DealNum             int     `json:"deal_num,omitempty"`     // 成交量
	DealAmount          float64 `json:"deal_amount,omitempty"`  // 成交额
	Buy1Num             int     `json:"buy_1_num,omitempty"`
	Buy1Price           float64 `json:"buy_1_price,omitempty"`
	Buy2Num             int     `json:"buy_2_num,omitempty"`
	Buy2Price           float64 `json:"buy_2_price,omitempty"`
	Buy3Num             int     `json:"buy_3_num,omitempty"`
	Buy3Price           float64 `json:"buy_3_price,omitempty"`
	Buy4Num             int     `json:"buy_4_num,omitempty"`
	Buy4Price           float64 `json:"buy_4_price,omitempty"`
	Buy5Num             int     `json:"buy_5_num,omitempty"`
	Buy5Price           float64 `json:"buy_5_price,omitempty"`
	Sell1Num            int     `json:"sell_1_num,omitempty"`
	Sell1Price          float64 `json:"sell_1_price,omitempty"`
	Sell2Num            int     `json:"sell_2_num,omitempty"`
	Sell2Price          float64 `json:"sell_2_price,omitempty"`
	Sell3Num            int     `json:"sell_3_num,omitempty"`
	Sell3Price          float64 `json:"sell_3_price,omitempty"`
	Sell4Num            int     `json:"sell_4_num,omitempty"`
	Sell4Price          float64 `json:"sell_4_price,omitempty"`
	Sell5Num            int     `json:"sell_5_num,omitempty"`
	Sell5Price          float64 `json:"sell_5_price,omitempty"`
	DateTime            int64   `json:"date_time,omitempty"`
}

func SinaStrToEntry(src string) ([]*Entry, error) {
	list := strings.Split(src, elemSeparator)
	resList := make([]*Entry, 0)
	for _, v := range list {
		infos := strings.Split(v, infoSeparator)
		if len(infos) != 2 {
			zlog.Error("SinaStrToEntry element is not separated by ;", zap.String("src", v))
			continue
		}

		prefix := infos[0]
		datas := strings.Trim(infos[1], "\"")
		pindex := strings.LastIndex(prefix, "_")
		if pindex < 0 || pindex+1 >= len(prefix) {
			zlog.Error("SinaStrToEntry prefix info can not found", zap.String("prefix", prefix))
			continue
		}
		markStock := prefix[pindex+1:]
		if len(markStock) != stockLen {
			zlog.Error("SinaStrToEntry stock length is wrong", zap.String("stock", markStock))
			continue
		}

		market := markStock[:2]
		stockCode := markStock[2:]
		e := &Entry{}
		e.Market = market
		e.Code = stockCode

		dataList := strings.Split(datas, ",")
		builder := strings.Builder{}
		builder.WriteString("{")
		for i, v := range dataList {
			k, ok := sinaStockJsonKey[i]
			if !ok {
				continue
			}
			if i == 0 {
				builderWithString(&builder, k, v)
				continue
			}
			if i == dateIndex || i == timeIndex {
				continue
			}
			builder.WriteByte(',')
			builderWithNumber(&builder, k, v)
		}
		dateTime := dataList[dateIndex] + " " + dataList[timeIndex]
		t, err := time.ParseInLocation(time.DateTime, dateTime, time.Local)
		if err != nil {
			zlog.Error("SinaStrToEntry parse time error", zap.Error(err), zap.String("dataTime", dateTime))
			t = time.Now()
		}

		builder.WriteByte(',')
		builderWithNumber(&builder, dateTimeKey, fmt.Sprint(t.Unix()))
		builder.WriteString("}")
		bin := builder.String()
		err = jsoniter.Unmarshal([]byte(bin), e)
		if err != nil {
			zlog.Error("SinaStrToEntry unmarshal data failed", zap.Error(err), zap.String("src", bin))
			continue
		}
		resList = append(resList, e)
	}
	return resList, nil
}

func reflectKey(e *Entry, k, v string) error {
	vinfo := reflect.ValueOf(e)
	tinfo := reflect.TypeOf(e)
	field, ok := tinfo.Elem().FieldByName(k)
	if !ok {
		return fmt.Errorf("not found field %s", k)
	}

	var value interface{}
	switch field.Type.Kind() {
	case reflect.String:
		value = v
	case reflect.Int:
		value = cast.ToInt(v)
	case reflect.Float64:
		value = cast.ToFloat64(v)
	case reflect.Float32:
		value = cast.ToFloat32(v)
	case reflect.Int32:
		value = cast.ToInt32(v)
	default:
		if field.Name == "Date" {
			value, _ = time.ParseInLocation(time.DateOnly, v, time.Local)
		} else if field.Name == "Time" {
			value, _ = time.ParseInLocation(time.TimeOnly, v, time.Local)
		}
	}

	vinfo.Elem().FieldByName(k).Set(reflect.ValueOf(value))
	return nil
}

func builderWithString(builder *strings.Builder, k, v string) {
	builder.WriteByte('"')
	builder.WriteString(k)
	builder.WriteByte('"')
	builder.WriteString(":")
	builder.WriteByte('"')
	builder.WriteString(v)
	builder.WriteByte('"')
}
func builderWithNumber(builder *strings.Builder, k, v string) {
	builder.WriteByte('"')
	builder.WriteString(k)
	builder.WriteByte('"')
	builder.WriteString(":")
	builder.WriteString(v)
}
