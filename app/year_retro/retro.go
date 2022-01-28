package year_retro

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	URL string = "https://emap.pcsc.com.tw/EMapSDK.aspx"
)

func Retro() {
	//start, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00+08:00")
	//end, _ := time.Parse(time.RFC3339, "2021-12-31T00:00:00+08:00")

	filter := bson.M{
		//"onSoldAt": bson.M{"$lt": t},
		//"status": bson.M{
		//"$in": []int32{1, 5},
		//},
		//"status.buyOrder.status": 7,
		//"status.sellOrder.status": 3,
		//"status.sellOrder.finishAt": bson.M{
		//"$gte": start,
		//"$lte": end,
		//},
		//"createdAt": bson.M{
		//"$gte": start,
		//"$lte": end,
		//},
	}

	//str := []string{}

	//for _, v := range buyOrder {
	//if v.Memo == "" {
	//continue
	//}
	//str = append(str, v.Memo)
	//}

	//j, _ := json.Marshal(str)
	//spew.Dump(string(j))
	//ioutil.WriteFile("拒絕理由", j, 0644)

	memberDeliveryData := MemberDeliverySearch(filter)
	//spew.Dump(memberDeliveryData)
	//sellorderData := generate_csv.SellOrderSearch(filter)
	//spew.Dump(len(sellorderData))

	del, g := Process(memberDeliveryData)
	//del, sellorderGroup := Process(sellorderData)
	//_ = SellOrderProcess(sellorderData, sellorderGroup)

	AgeCenter, MaxAge, MinAge := MemberProcess(g)

	spew.Dump(AgeCenter)
	spew.Dump(MaxAge)
	spew.Dump(MinAge)
	spew.Dump(del)

	//spew.Dump("---------------")

}

func BuyOrderProcess(b []*BuyOrder) map[string]int64 {
	BuyOrderMap := map[string]int64{}
	for _, v := range b {
		_, ok := BuyOrderMap[v.Memo]
		if !ok {
			BuyOrderMap[v.Memo] = 1
			continue
		}
		BuyOrderMap[v.Memo]++
	}
	return BuyOrderMap
}

// 計算各地區消費金額
func SellOrderProcess(s []*Delivery, g map[string][]string) map[string]float64 {

	sellOrderMap := map[string]int64{}
	for _, v := range s {
		sellOrderMap[v.No] = v.Price
	}

	sellOrderPrice := map[string]int64{}

	var total int64 = 0

	for address, nos := range g {

		if address == "亞太雲端" {
			address = "台中市"
		}
		for _, no := range nos {
			price, ok := sellOrderMap[no]
			if !ok {
				continue
			}
			total += price

			_, ok = sellOrderPrice[address]
			if !ok {
				sellOrderPrice[address] = price
				continue
			}

			sellOrderPrice[address] += price
		}

	}

	// 總額

	sellOrderRate := map[string]float64{}
	decTotal := decimal.NewFromInt(total)

	// 各地區消費佔比
	for address, price := range sellOrderPrice {
		decPrice := decimal.NewFromInt(price)

		decRate := decPrice.Div(decTotal)

		rate := decimal.NewFromInt(100)
		sellOrderRate[address], _ = decRate.Mul(rate).Round(2).Float64()
	}

	return sellOrderRate
}

// 計算各地區的平均年齡
func MemberProcess(memberGroup map[string][]string) (map[string]int, map[string]int, map[string]int) {

	var IDs []string
	for _, v := range memberGroup {
		IDs = append(IDs, v...)
	}

	start, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00+08:00")
	end, _ := time.Parse(time.RFC3339, "2021-12-31T00:00:00+08:00")
	// 會員原始資料
	memberRaw := MemberSearch(bson.M{
		"no": bson.M{
			"$in": IDs,
		},

		"createdAt": bson.M{
			"$gte": start,
			"$lte": end,
		},
	})

	// 轉成map
	memberMap := map[string]*Member{}
	for _, v := range memberRaw {
		memberMap[v.No] = v
	}

	// 地區下的會員資料
	memberNewGroup := map[string][]*Member{}

	// 遍歷各地區的會員編號
	for address, Ids := range memberGroup {
		// 取出的會員跟原始資料map找Member結構
		for _, ID := range Ids {
			member, _ := memberMap[ID]

			_, ok := memberNewGroup[address]
			if !ok {
				memberNewGroup[address] = []*Member{member}
				continue
			}
			memberNewGroup[address] = append(memberNewGroup[address], member)
		}
	}

	// 年齡中位數
	result := map[string]int{}
	// 最大年齡
	max := map[string]int{}
	// 最小年齡
	min := map[string]int{}
	// 計算年齡中位數

	now := time.Now().Year()
	// 計算平均
	for address, members := range memberNewGroup {
		if address == "" {
			continue
		}
		if address == "亞太雲端" {
			address = "台中市"
		}

		//sumAge = 0
		//total = 0
		var (
			// 年齡陣列
			AgeStr []int
		)
		for _, member := range members {
			if member == nil || member.Birthday.IsZero() {
				continue
			}
			//total++

			age := now - member.Birthday.Year()
			if age == 2 {
				continue
			}
			AgeStr = append(AgeStr, age)

		}

		var ans int
		if AgeStr == nil {
			continue
		}
		sortAgeStr := sort.IntSlice(AgeStr)
		sort.Sort(sortAgeStr)
		AgeStr = sortAgeStr
		if address == "台中市" {
			spew.Dump(sortAgeStr)
		}

		// 計算中位數
		if residue := len(AgeStr) % 2; residue == 0 {
			// 偶數
			if len(AgeStr) == 2 {
				ans = (AgeStr[0] + AgeStr[1]) / 2
			} else {
				// 0 1 2 3 4 5 6 7
				i := len(AgeStr) / 2
				ans = (AgeStr[i] + AgeStr[i+1]) / 2
			}
		} else {
			index := float64(len(AgeStr)/2) + 0.5

			ans = AgeStr[int(index)]
		}

		result[address] = ans
		max[address] = AgeStr[len(AgeStr)-1]
		min[address] = AgeStr[0]

	}

	return result, max, min
}

func Process(raw []*Delivery) (map[string]int64, map[string][]string) {
	add := map[string]int64{}
	group := map[string][]string{}

	//var trush []string
	for _, v := range raw {
		if v.Type == "面交" {
			v.Address.City = "台中市"
		}

		if v.Type == "超商取貨" {
			strings.TrimSpace(v.Address.Addrs)
			p0 := strings.TrimSpace(v.Address.Addrs)
			p1 := Replace(p0)
			p5 := strings.Split(p1, "/")
			p6 := p5[0]
			if len(p5) > 1 {

				for _, v := range p5 {
					_, err := cast.ToInt64E(v)
					if err != nil {
						p6 = v
					}

				}
			}

			p7 := strings.Split(p6, "／")
			p8 := p7[0]
			if len(p7) > 1 {
				for _, v := range p7 {
					_, err := cast.ToInt64E(v)
					if err != nil {
						p8 = v
					}

				}
			}
			p9 := strings.Split(p8, " ")
			p10 := p9[0]
			if len(p9) > 1 {
				for _, v := range p9 {
					_, err := cast.ToInt64E(v)
					if err != nil {
						p10 = v
					}

				}
			}

			if len([]rune(p10)) > 5 {
				_, err := cast.ToInt64E(string([]rune(p10)[0:6]))
				if err == nil {
					p10 = string([]rune(p10)[0:6])
				}
			}

			p99 := GetAddress(p10)
			if p99 == "" {
				//trush = append(trush, p10)
				continue
			}
			v.Address.City = p99

		}

		if v.Address.City == "臺中市" {
			v.Address.City = "台中市"
		}
		if v.Address.City == "臺南市" {
			v.Address.City = "台南市"
		}
		if v.Address.City == "臺東縣" {
			v.Address.City = "台東縣"
		}
		if v.Address.City == "臺北市" {
			v.Address.City = "台北市"
		}
		_, ok := add[v.Address.City]
		if !ok {
			add[v.Address.City] = 1
			group[v.Address.City] = []string{v.No}
			continue
		}

		add[v.Address.City]++
		group[v.Address.City] = append(group[v.Address.City], v.No)
	}

	return add, group
}

func Replace(raw string) string {
	var NeedReplace []string = []string{
		"711",
		"7-11",
		"店",
		"門市",
		"號：",
		"名：",
		"（",
		"）",
		"。",
		"宜蘭",
		"雲林",
		"新莊",
		"台中",
		"-",
		"台南",
		"宜蘭縣",
		"三重",
		"台中市",
		"屏東",
		"大同",
		"新莊",
		"高雄",
		"新竹縣",
		"湖口",
		"新竹縣湖口",
		"永和",
		"逢福",
		"蘆洲",
		"神岡",
		"（731台南市後壁區福安里下寮200號）",
		"高雄市仁武區",
		"新崙",
		"（高雄市鳳山區中崙路501號）",
		"淡水",
		"明義",
		"樹林",
		"新北泰山區",
		"三重",
		"高雄",
		"龍星",
		"桃園",
		"高雄",
		"新莊區中正路72號",
		"新竹",
		"縣",
		"市",
		"冬山鄉廣興路的",
		"仁武區",
		"區中崙路501號",
		"區義華路171號",
		"文忠號",
	}
	ans := raw
	for _, v := range NeedReplace {
		ans = strings.Replace(ans, v, "", -1)
	}
	return ans
}

func GetAddress(raw string) string {
	p1 := strings.TrimSpace(raw)
	p2 := GetShopAddress(p1)
	p3 := FindAddress(p2)

	return p3
}

func FindAddress(raw string) string {
	var p int
	for _, v := range TaiwanCity {

		p = strings.Index(raw, v)

		if p > -1 {
			break
		}
	}
	if p < 0 {
		return ""
	}

	address := []byte(raw)[p : p+9]
	return string(address)
}

func GetShopAddress(binder string) string {

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	raw := setBody(binder)

	for k, v := range raw {
		w.WriteField(k, v)
	}
	w.Close()

	req, _ := http.NewRequest("POST", URL, body)
	req.Header.Set("Content-Type", w.FormDataContentType())

	req.Header.Set("Cookie", "SET_EC_COOKIE=rd1378o00000000000000000000ffff0ac80808o443; citrix_ns_id=c2zWUEYAPyP/rBkq+DipmXLEPEg0001; ECMap=eshopparid=899,eshopid=899,url=/secure.rakuten.com.tw/checkout/callback,tempvar=,sid=1,storecategory=3,showtype=1,oStoreId=; ECMapLoginToken=Qu8XhfHDTTmZ2q9pggRiWoVdPnri7fHwzttaBOErlo20220103160510; _ga=GA1.3.208069652.1641195311; _gid=GA1.3.365931951.1641195311; citrix_ns_id=c2zWUEYAPyP/rBkq+DipmXLEPEg0001; ASP.NET_SessionId=4ya5eafa2ivk3gxlph3dkqnh")

	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(data)

}

func setBody(binder string) map[string]string {
	form := map[string]string{}

	form["commandid"] = "SearchStore"
	_, err := cast.ToInt64E(binder)
	//spew.Dump(err)
	if err == nil {

		form["ID"] = binder
		return form
	}

	form["StoreName"] = binder

	return form
}

func setHeader(r *resty.Request) {
	r.SetHeader("Cookie", "SET_EC_COOKIE=rd1378o00000000000000000000ffff0ac80808o443; citrix_ns_id=c2zWUEYAPyP/rBkq+DipmXLEPEg0001; ECMap=eshopparid=899,eshopid=899,url=/secure.rakuten.com.tw/checkout/callback,tempvar=,sid=1,storecategory=3,showtype=1,oStoreId=; ECMapLoginToken=Qu8XhfHDTTmZ2q9pggRiWoVdPnri7fHwzttaBOErlo20220103160510; _ga=GA1.3.208069652.1641195311; _gid=GA1.3.365931951.1641195311; citrix_ns_id=c2zWUEYAPyP/rBkq+DipmXLEPEg0001; ASP.NET_SessionId=4ya5eafa2ivk3gxlph3dkqnh")
	r.SetHeader("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
}

var TaiwanCity []string = []string{
	"台北市",
	"新北市",
	"桃園市",
	"台中市",
	"台南市",
	"高雄市",
	"新竹縣",
	"苗栗縣",
	"彰化縣",
	"南投縣",
	"雲林縣",
	"嘉義縣",
	"屏東縣",
	"宜蘭縣",
	"花蓮縣",
	"台東縣",
	"澎湖縣",
	"金門縣",
	"連江縣",
	"基隆市",
	"新竹市",
	"嘉義市",
}
