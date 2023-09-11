package stock

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"strings"
)

var (
	sinaPriceUrl = "https://hq.sinajs.cn/list="
)

func Get(markCodes []string) (string, error) {
	c := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	list := strings.Join(markCodes, ",")
	req, err := http.NewRequest("GET", sinaPriceUrl+list, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Referer", "https://finance.sina.com.cn")
	rsp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	r := transform.NewReader(rsp.Body, simplifiedchinese.GBK.NewDecoder())
	res, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

//var hq_str_sz002603="以岭药业,21.970,21.970,22.400,22.550,21.950,22.400,22.410,17909232,400172176.920,132200,22.400,17100,22.390,45600,22.380,59400,22.370,43700,22.360,9200,22.410,36900,22.420,25500,22.430,19700,22.440,54180,22.450,2023-09-11,14:34:42,00";
//var hq_str_sz002604="龙力退,0.000,0.000,0.000,0.000,0.000,0.000,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,2023-09-08,11:45:00,-

// https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeStockCount?node=sz_a
