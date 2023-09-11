package stock

import (
	"fmt"
	"github.com/jackylong1/lltrade/pkg/zlog"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	stockCodeUrl = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeStockCount?node=%s"

	stockCodeLimitUrl = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=%d&num=%d&sort=symbol&asc=1&node=%s&symbol=&_s_r_a=init"
)

type Code struct {
	Symbol string `json:"symbol"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Trade  string `json:"trade"`
	//Pricechange   float64 `json:"pricechange"`
	//Changepercent float64 `json:"changepercent"`
	Buy           string  `json:"buy"`
	Sell          string  `json:"sell"`
	Settlement    string  `json:"settlement"`
	Open          string  `json:"open"`
	High          string  `json:"high"`
	Low           string  `json:"low"`
	Volume        int     `json:"volume"`
	Amount        int     `json:"amount"`
	Ticktime      string  `json:"ticktime"`
	Per           float64 `json:"per"`
	Pb            float64 `json:"pb"`
	Mktcap        float64 `json:"mktcap"`
	Nmc           float64 `json:"nmc"`
	Turnoverratio float64 `json:"turnoverratio"`
}

//https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=1&num=80&sort=symbol&asc=1&node=sz_a&symbol=&_s_r_a=init

/*
curl 'https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=1&num=80&sort=symbol&asc=1&node=sz_a&symbol=&_s_r_a=init' \
  -H 'authority: vip.stock.finance.sina.com.cn' \
  -H 'accept: */
//-H 'accept-language: zh-CN,zh;q=0.9' \
//-H 'content-type: application/x-www-form-urlencoded' \
//-H 'cookie: UOR=www.google.com,www.sina.com.cn,; SINAGLOBAL=120.239.154.105_1673515673.335828; SGUID=1673515673062_37867125; __bid_n=185a54ea722fa5a84a4207; Hm_lvt_b82ffdf7cbc70caaacee097b04128ac1=1675582156; FPTOKEN=x+v36ZT/x2buKodK8d2UsDFP2SH6psdc1j89jqtGQ4YRe1+tN5kfIQrxBxg6qniwdNGWpqLXwKelL11MTx9xFR+7k80rnOyaLgAW4ein1FJwF1yZ+1jhBiiE2Zs5X1v6KrmuCLG3Sl1tWGRSyO3S0up/lopM3VXEnKaAWLYExpqM//K2NCvatDXvqRqSDC4GbA7G7HMBy+PQoye0gXgyxQHRQGQe76t16ym6qMH94YPDLVQHGxKZnOPR/wGYdTqFMXsBoKrEZU9m8wyq3a7if4jmiUGjy5nBP5xYj2Vged4GXthHvfEEavRedM7Qj/ZF2k/kBm0c69ZbqsOn6ylHJxGrB7krynayW+6/Y5z89tIwcK9EP+vfn2kW2LU1BeTCAf2v70lXifvmAZk34hwjdA==|CXxB85kzWl/h+TLtNahmpc3RdEJNc5ttS2uOx35TnBI=|10|be814376ab38f507b9f31bc4f49250fb; SUB=_2AkMTGxpPf8NxqwFRmPEWymniaIt1wg3EieKlR-uUJRMyHRl-yD9vqmUHtRB6OJs0oHmpIK5HScAmFSMDhhgPxZMsjI1k; SR_SEL=1_511; MONEY-FINANCE-SINA-COM-CN-WEB5=; ULV=1694431425121:4:1:1::1688950174532; Apache=112.96.173.218_1694431425.507341' \
//-H 'referer: https://vip.stock.finance.sina.com.cn/mkt/' \
//-H 'sec-ch-ua: "Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"' \
//-H 'sec-ch-ua-mobile: ?0' \
//-H 'sec-ch-ua-platform: "macOS"' \
//-H 'sec-fetch-dest: empty' \
//-H 'sec-fetch-mode: cors' \
//-H 'sec-fetch-site: same-origin' \
//-H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36' \
//--compressed
//*/

func GetStockCodeCount(market string) (int, error) {
	c := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	markUrl := fmt.Sprintf(stockCodeUrl, market)
	req, err := http.NewRequest("GET", markUrl, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("referer", "https://vip.stock.finance.sina.com.cn/mkt/")
	rsp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()
	res, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}
	str := strings.Trim(string(res), "\"")
	count, err := strconv.Atoi(str)
	return count, err
}
func GetStockCodeList(market string) {
	c := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	total, err := GetStockCodeCount(market)
	if err != nil {
		zlog.Error("GetStockCodeList-total count failed", zap.Error(err))
		return
	}

	pageSize := 80
	pageTotal := total/pageSize + 1
	var outList []*Code
	for page := 1; page <= pageTotal; page++ {
		codeUrl := fmt.Sprintf(stockCodeLimitUrl, page, pageSize, market)
		req, err := http.NewRequest("GET", codeUrl, nil)
		if err != nil {
			zlog.Error("GetStockCodeList-new http request", zap.Error(err), zap.String("url", codeUrl))
			continue
		}
		req.Header.Add("referer", "https://vip.stock.finance.sina.com.cn/mkt/")
		rsp, err := c.Do(req)
		if err != nil {
			zlog.Error("GetStockCodeList-Http-Get failed", zap.Error(err), zap.String("url", codeUrl))
			continue
		}

		r := transform.NewReader(rsp.Body, simplifiedchinese.GBK.NewDecoder())
		res, err := io.ReadAll(r)
		if err != nil {
			zlog.Error("GetStockCodeList-Http-Result read failed", zap.Error(err), zap.String("url", codeUrl))
			rsp.Body.Close()
			continue
		}
		var list []*Code
		err = jsoniter.Unmarshal(res, &list)
		rsp.Body.Close()
		if err != nil {
			zlog.Error("GetStockCodeList-Unmarshal-Result failed", zap.Error(err), zap.String("res", string(res)))
			continue
		}
		outList = append(outList, list...)
	}
	writeFile(outList, market)
}

func writeFile(datas []*Code, market string) {
	fName := market + "_stock.txt"
	f, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	for _, v := range datas {
		bin, err := jsoniter.Marshal(v)
		if err != nil {
			zlog.Error("writeFile- jsoniter.Marshal failed", zap.Error(err))
			continue
		}
		f.WriteString(string(bin) + ",\n")

	}
	f.Sync()
}
