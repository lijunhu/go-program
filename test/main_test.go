package test

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"testing"
)

func TestLeetCode(t *testing.T) {
	req := httplib.Post("http://10.172.20.17/hotelListV5/getHotelListV5Kiwi")

	req.Header("latitude", "39.986971")
	req.Header("phonemodel", "iPhone11,6")
	req.Header("positioning", "0")
	req.Header("version", "9.73.0")
	req.Header("innerfrom", "10000")
	req.Header("smdeviceid", "201909061057579f8be63e9015f4f9429052a95005d593015b4ff936b09408")
	req.Header("cloudtype", "eLong")
	req.Header("chid", "ewiphone")
	req.Header("network", "Wifi")
	req.Header("appfrom", "1")
	req.Header("channelid", "ewiphone")
	req.Header("deviceid", "FF975ADB-0956-49CB-8E38-B88B3E8736F9")
	req.Header("coorsys", "0")
	req.Header("usertraceid", "0239EE6A-3D14-4BEC-A5CA-05596B58A4BA")
	req.Header("dimension", "1242*2688")
	req.Header("longitude", "116.507834")
	req.Header("saviortraceid", "1604646388739")
	req.Header("idfv", "68CA7C00-6175-403F-9EB4-BFF81DD510CF")
	req.Header("appname", "com.elong.app")
	req.Header("channel", "hotelgeneral")
	req.Header("priority", "0")
	req.Header("clienttype", "1")
	req.Header("hotelgroup", "1")
	req.Header("localtime", "1604646388799")
	req.Header("elongdebugnetworkrequestid", "1604646388802_hotel/hotelListV4")
	req.Header("phonebrand", "iPhone")
	req.Header("traceid", "6A613DB1-DB0D-45D5-844C-3E623B649C92")
	req.Header("interceptaction", "0")
	req.Header("outerfrom", "20000")
	req.Header("osversion", "iphone_14.0.1")
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.SetHost("dueros.hd.elong.com")
	req.SetUserAgent("ElongClient/2103 CFNetwork/1197 Darwin/20.0.0")
	req.Param("req","%7b%22controlTag%22%3a17179869248%2c%22dataVersion%22%3a%222.0%22%2c%22userPropertyCtripPromotion%22%3a0%2c%22AreaType%22%3a0%2c%22StarCode%22%3a%22-1%22%2c%22HighestPrice%22%3a-1%2c%22IsAtCurrentCity%22%3atrue%2c%22SearchTraceID%22%3a%225745A9F5-B48E-4361-BB2E-9F44D09358AD%22%2c%22CheckOutDate%22%3a%222020-11-07%22%2c%22imageMode%22%3a0%2c%22AreaName%22%3anull%2c%22businessTrip%22%3afalse%2c%22SearchType%22%3a0%2c%22Radius%22%3a0%2c%22businessTravelPop%22%3afalse%2c%22IsPositioning%22%3afalse%2c%22HotelName%22%3a%22%22%2c%22HotelBrandID%22%3a%220%22%2c%22Longitude%22%3a0%2c%22PageIndex%22%3a0%2c%22CityID%22%3a%220101%22%2c%22Latitude%22%3a0%2c%22hotelFilterDatas%22%3a%5b%7b%22filterId%22%3a1%2c%22typeId%22%3a8888%7d%5d%2c%22imageFlag%22%3a22%2c%22GuestGPS%22%3a%7b%22Latitude%22%3a39.986971009747933%2c%22Longtitude%22%3a116.50783354264671%2c%22LocationType%22%3a0%7d%2c%22Currency%22%3a%22RMB%22%2c%22IsSearchAgain%22%3atrue%2c%22searchActivityId%22%3a%22%22%2c%22timeZone%22%3a8%2c%22searchEntranceId%22%3a%22home2search%22%2c%22hotelFilterFlag%22%3a%22%22%2c%22IsApartment%22%3afalse%2c%22OrderBy%22%3a0%2c%22hotelNum%22%3a1%2c%22PriceLevel%22%3a0%2c%22Filter%22%3a0%2c%22IntelligentSearchText%22%3a%22%22%2c%22TalentRecomendImageSize%22%3a23%2c%22bigOperatingTipCacheInfos%22%3a%5b%5d%2c%22IsShowPsgHotel%22%3atrue%2c%22MutilpleFilter%22%3a%221460%22%2c%22CheckInDate%22%3a%222020-11-06%22%2c%22ehActivityId%22%3a%221110%22%2c%22HasHongbao%22%3atrue%2c%22IsAroundSale%22%3afalse%2c%22MemberLevel%22%3a0%2c%22AreaId%22%3anull%2c%22RequestFrom%22%3a0%2c%22CityName%22%3a%22%e5%8c%97%e4%ba%ac%22%2c%22PageSize%22%3a20%2c%22inter%22%3a0%2c%22LowestPrice%22%3a-1%7d")
	resp, _ := req.Response()
	fmt.Println(resp.Status)

	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}
