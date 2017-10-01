package main

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/ini.v1"
)

func init() {
	RegexpUrlTwitter, _ = regexp.Compile(REGEXP_URL_TWITTER)
	RegexpUrlTistory, _ = regexp.Compile(REGEXP_URL_TISTORY)
	RegexpUrlTistoryWithCDN, _ = regexp.Compile(REGEXP_URL_TISTORY_WITH_CDN)
	RegexpUrlGfycat, _ = regexp.Compile(REGEXP_URL_GFYCAT)
	RegexpUrlInstagram, _ = regexp.Compile(REGEXP_URL_INSTAGRAM)
	RegexpUrlImgurSingle, _ = regexp.Compile(REGEXP_URL_IMGUR_SINGLE)
	RegexpUrlImgurAlbum, _ = regexp.Compile(REGEXP_URL_IMGUR_ALBUM)
	RegexpUrlGoogleDrive, _ = regexp.Compile(REGEXP_URL_GOOGLEDRIVE)
	RegexpUrlPossibleTistorySite, _ = regexp.Compile(REGEXP_URL_POSSIBLE_TISTORY_SITE)
	RegexpUrlFlickrPhoto, _ = regexp.Compile(REGEXP_URL_FLICKR_PHOTO)
	RegexpUrlFlickrAlbum, _ = regexp.Compile(REGEXP_URL_FLICKR_ALBUM)
	RegexpUrlStreamable, _ = regexp.Compile(REGEXP_URL_STREAMABLE)
	flickrApiKey = os.Getenv("FLICKR_API_KEY")

	var err error
	cfg, err := ini.Load("config.ini")
	if err == nil {
		flickrApiKey = cfg.Section("flickr").Key("api key").MustString("yourflickrapikey")
	}
}

type urlsTestpair struct {
	value  string
	result map[string]string
}

var getTwitterUrlsTests = []urlsTestpair{
	{
		"https://pbs.twimg.com/media/CulDBM6VYAA-YhY.jpg:orig",
		map[string]string{"https://pbs.twimg.com/media/CulDBM6VYAA-YhY.jpg:orig": "CulDBM6VYAA-YhY.jpg"},
	},
	{
		"https://pbs.twimg.com/media/CulDBM6VYAA-YhY.jpg",
		map[string]string{"https://pbs.twimg.com/media/CulDBM6VYAA-YhY.jpg:orig": "CulDBM6VYAA-YhY.jpg"},
	},
	{
		"http://pbs.twimg.com/media/CulDBM6VYAA-YhY.jpg",
		map[string]string{"https://pbs.twimg.com/media/CulDBM6VYAA-YhY.jpg:orig": "CulDBM6VYAA-YhY.jpg"},
	},
}

func TestGetTwitterUrls(t *testing.T) {
	for _, pair := range getTwitterUrlsTests {
		v, err := getTwitterUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getTistoryUrlsTests = []urlsTestpair{
	{
		"http://cfile25.uf.tistory.com/original/235CA739582E86992EFC4E",
		map[string]string{"http://cfile25.uf.tistory.com/original/235CA739582E86992EFC4E": ""},
	},
	{
		"http://cfile25.uf.tistory.com/image/235CA739582E86992EFC4E",
		map[string]string{"http://cfile25.uf.tistory.com/original/235CA739582E86992EFC4E": ""},
	},
}

func TestGetTistoryUrls(t *testing.T) {
	for _, pair := range getTistoryUrlsTests {
		v, err := getTistoryUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getGfycatUrlsTests = []urlsTestpair{
	{
		"https://gfycat.com/SandyChiefBoubou",
		map[string]string{"https://fat.gfycat.com/SandyChiefBoubou.mp4": ""},
	},
}

func TestGetGfycatUrls(t *testing.T) {
	for _, pair := range getGfycatUrlsTests {
		v, err := getGfycatUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getInstagramUrlsPictureTests = []urlsTestpair{
	{
		"https://www.instagram.com/p/BHhDAmhAz33/?taken-by=s_sohye",
		map[string]string{"https://www.instagram.com/p/BHhDAmhAz33/media/?size=l&taken-by=s_sohye": "instagram s_sohye - BHhDAmhAz33.jpg"},
	},
	{
		"https://www.instagram.com/p/BHhDAmhAz33/",
		map[string]string{"https://www.instagram.com/p/BHhDAmhAz33/media/?size=l": "instagram s_sohye - BHhDAmhAz33.jpg"},
	},
}

var getInstagramUrlsVideoTests = []urlsTestpair{
	{
		"https://www.instagram.com/p/BL2_ZIHgYTp/?taken-by=s_sohye",
		map[string]string{"14811404_233311497085396_338650092456116224_n.mp4": "instagram s_sohye - BL2_ZIHgYTp.mp4"},
	},
	{
		"https://www.instagram.com/p/BL2_ZIHgYTp/",
		map[string]string{"14811404_233311497085396_338650092456116224_n.mp4": "instagram s_sohye - BL2_ZIHgYTp.mp4"},
	},
}

var getInstagramUrlsAlbumTests = []urlsTestpair{
	{
		"https://www.instagram.com/p/BRiCc0VjULk/?taken-by=gfriendofficial",
		map[string]string{
			"17265460_395888184109957_3500310922180689920_n.jpg":  "instagram gfriendofficial - BRiCc0VjULk",
			"17265456_267171360360765_8110946520456495104_n.jpg":  "instagram gfriendofficial - BRiCc0VjULk",
			"17265327_1394797493912862_2677004307588448256_n.jpg": "instagram gfriendofficial - BRiCc0VjULk"},
	},
	{
		"https://www.instagram.com/p/BRhheSPjaQ3/",
		map[string]string{
			"17125875_306909746390523_8184965703367917568_n.jpg": "instagram gfriendofficial - BRhheSPjaQ3",
			"17266053_188727064951899_2485556569865977856_n.jpg": "instagram gfriendofficial - BRhheSPjaQ3"},
	},
}

func TestGetInstagramUrls(t *testing.T) {
	for _, pair := range getInstagramUrlsPictureTests {
		v, err := getInstagramUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
	for _, pair := range getInstagramUrlsVideoTests {
		v, err := getInstagramUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		for keyResult, valueResult := range pair.result {
			for keyExpected, valueExpected := range v {
				if strings.Contains(keyResult, keyExpected) || valueResult != valueExpected { // CDN location can vary
					t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
				}
			}
		}
	}
	for _, pair := range getInstagramUrlsAlbumTests {
		v, err := getInstagramUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		for keyResult, valueResult := range pair.result {
			for keyExpected, valueExpected := range v {
				if strings.Contains(keyResult, keyExpected) || strings.Contains(valueResult, valueExpected) { // CDN location can vary
					t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
				}
			}
		}
	}
}

var getImgurSingleUrlsTests = []urlsTestpair{
	{
		"http://imgur.com/viZictl",
		map[string]string{"http://imgur.com/download/viZictl": ""},
	},
	{
		"https://imgur.com/viZictl",
		map[string]string{"https://imgur.com/download/viZictl": ""},
	},
	{
		"https://i.imgur.com/viZictl.jpg",
		map[string]string{"https://i.imgur.com/download/viZictl.jpg": ""},
	},
	{
		"http://imgur.com/uYwt2VV",
		map[string]string{"http://imgur.com/download/uYwt2VV": ""},
	},
	{
		"http://i.imgur.com/uYwt2VV.gifv",
		map[string]string{"http://i.imgur.com/download/uYwt2VV": ""},
	},
}

func TestGetImgurSingleUrls(t *testing.T) {
	for _, pair := range getImgurSingleUrlsTests {
		v, err := getImgurSingleUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getImgurAlbumUrlsTests = []urlsTestpair{
	{
		"http://imgur.com/a/ALTpi",
		map[string]string{
			"https://i.imgur.com/FKoguPh.jpg": "",
			"https://i.imgur.com/5FNL6Pe.jpg": "",
			"https://i.imgur.com/YA0V0g9.jpg": "",
			"https://i.imgur.com/Uc2iDhD.jpg": "",
			"https://i.imgur.com/J9JRSSJ.jpg": "",
			"https://i.imgur.com/Xrx0uyE.jpg": "",
			"https://i.imgur.com/3xDSq1O.jpg": "",
		},
	},
}

func TestGetImgurAlbumUrls(t *testing.T) {
	for _, pair := range getImgurAlbumUrlsTests {
		v, err := getImgurAlbumUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		for expectedLink, expectedName := range pair.result {
			linkFound := false
			for gotLink, gotName := range v {
				if expectedLink == gotLink && expectedName == gotName {
					linkFound = true
				}
			}
			if !linkFound {
				t.Errorf("For expected %s %s, got %s", expectedLink, expectedName, v)
			}
		}
	}
}

var getGoogleDriveUrlsTests = []urlsTestpair{
	{
		"https://drive.google.com/file/d/0B8TnwsJqlFllSUtvUEhoSU40WkE/view",
		map[string]string{"https://drive.google.com/uc?export=download&id=0B8TnwsJqlFllSUtvUEhoSU40WkE": ""},
	},
}

func TestGetGoogleDriveUrls(t *testing.T) {
	for _, pair := range getGoogleDriveUrlsTests {
		v, err := getGoogleDriveUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getTistoryWithCDNUrlsTests = []urlsTestpair{
	{
		"http://img1.daumcdn.net/thumb/R720x0.q80/?scode=mtistory&fname=http%3A%2F%2Fcfile24.uf.tistory.com%2Fimage%2F2658554B580BDC4C0924CA",
		map[string]string{"http://cfile24.uf.tistory.com/original/2658554B580BDC4C0924CA": ""},
	},
}

func TestGetTistoryWithCDNUrls(t *testing.T) {
	for _, pair := range getTistoryWithCDNUrlsTests {
		v, err := getTistoryWithCDNUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getPossibleTistorySiteUrlsTests = []urlsTestpair{
	{
		"http://soonduck.tistory.com/482",
		map[string]string{
			"a": "",
			"b": "",
			"c": "",
			"d": "",
			"e": "",
		},
	},
	{
		"http://soonduck.tistory.com/m/482",
		map[string]string{
			"a": "",
			"b": "",
			"c": "",
			"d": "",
			"e": "",
		},
	},
	{
		"http://slmn.de/123",
		map[string]string{},
	},
}

func TestGetPossibleTistorySiteUrls(t *testing.T) {
	for _, pair := range getPossibleTistorySiteUrlsTests {
		v, err := getPossibleTistorySiteUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if len(pair.result) != len(v) { // only check filenames, urls may vary
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getFlickrUrlFromPhotoIdTests = []map[string]string{
	{
		"value":  "31065043320",
		"result": "https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg",
	},
}

func TestGetFlickrUrlFromPhotoId(t *testing.T) {
	for _, pair := range getFlickrUrlFromPhotoIdTests {
		v := getFlickrUrlFromPhotoId(pair["value"])
		if v != pair["result"] {
			t.Errorf("For %s, expected %s, got %s", pair["value"], pair["result"], v)
		}
	}
}

var getFlickrPhotoUrlsTests = []urlsTestpair{
	{
		"https://www.flickr.com/photos/137385017@N08/31065043320/in/album-72157677350305446/",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
		},
	},
	{
		"https://www.flickr.com/photos/137385017@N08/31065043320/in/album-72157677350305446",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
		},
	},
	{
		"https://www.flickr.com/photos/137385017@N08/31065043320/",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
		},
	},
	{
		"https://www.flickr.com/photos/137385017@N08/31065043320",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
		},
	},
}

func TestGetFlickrPhotoUrls(t *testing.T) {
	for _, pair := range getFlickrPhotoUrlsTests {
		v, err := getFlickrPhotoUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getFlickrAlbumUrlsTests = []urlsTestpair{
	{
		"https://www.flickr.com/photos/137385017@N08/albums/72157677350305446/",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
			"https://farm6.staticflickr.com/5651/31434767515_49f88ee12e_b.jpg": "",
			"https://farm6.staticflickr.com/5750/31434766825_529fd08071_b.jpg": "",
			"https://farm6.staticflickr.com/5811/31319456971_37c8c4708a_b.jpg": "",
			"https://farm6.staticflickr.com/5494/30627074913_b7f810fc26_b.jpg": "",
			"https://farm6.staticflickr.com/5539/31065042720_d76f643b28_b.jpg": "",
			"https://farm6.staticflickr.com/5813/31434765285_94b85d5e8c_b.jpg": "",
			"https://farm6.staticflickr.com/5600/31065044090_eca63bd5a5_b.jpg": "",
			"https://farm6.staticflickr.com/5733/31434764435_350825477e_b.jpg": "",
			"https://farm6.staticflickr.com/5715/30627073573_b86e4b2c22_b.jpg": "",
			"https://farm6.staticflickr.com/5758/31289864222_5e3cca7e72_b.jpg": "",
			"https://farm6.staticflickr.com/5801/30627076673_5a32f3e562_b.jpg": "",
			"https://farm6.staticflickr.com/5538/31319458901_088858d7f1_b.jpg": "",
		},
	},
	{
		"https://www.flickr.com/photos/137385017@N08/albums/72157677350305446",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
			"https://farm6.staticflickr.com/5651/31434767515_49f88ee12e_b.jpg": "",
			"https://farm6.staticflickr.com/5750/31434766825_529fd08071_b.jpg": "",
			"https://farm6.staticflickr.com/5811/31319456971_37c8c4708a_b.jpg": "",
			"https://farm6.staticflickr.com/5494/30627074913_b7f810fc26_b.jpg": "",
			"https://farm6.staticflickr.com/5539/31065042720_d76f643b28_b.jpg": "",
			"https://farm6.staticflickr.com/5813/31434765285_94b85d5e8c_b.jpg": "",
			"https://farm6.staticflickr.com/5600/31065044090_eca63bd5a5_b.jpg": "",
			"https://farm6.staticflickr.com/5733/31434764435_350825477e_b.jpg": "",
			"https://farm6.staticflickr.com/5715/30627073573_b86e4b2c22_b.jpg": "",
			"https://farm6.staticflickr.com/5758/31289864222_5e3cca7e72_b.jpg": "",
			"https://farm6.staticflickr.com/5801/30627076673_5a32f3e562_b.jpg": "",
			"https://farm6.staticflickr.com/5538/31319458901_088858d7f1_b.jpg": "",
		},
	},
	{
		"https://www.flickr.com/photos/137385017@N08/albums/with/72157677350305446/",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
			"https://farm6.staticflickr.com/5651/31434767515_49f88ee12e_b.jpg": "",
			"https://farm6.staticflickr.com/5750/31434766825_529fd08071_b.jpg": "",
			"https://farm6.staticflickr.com/5811/31319456971_37c8c4708a_b.jpg": "",
			"https://farm6.staticflickr.com/5494/30627074913_b7f810fc26_b.jpg": "",
			"https://farm6.staticflickr.com/5539/31065042720_d76f643b28_b.jpg": "",
			"https://farm6.staticflickr.com/5813/31434765285_94b85d5e8c_b.jpg": "",
			"https://farm6.staticflickr.com/5600/31065044090_eca63bd5a5_b.jpg": "",
			"https://farm6.staticflickr.com/5733/31434764435_350825477e_b.jpg": "",
			"https://farm6.staticflickr.com/5715/30627073573_b86e4b2c22_b.jpg": "",
			"https://farm6.staticflickr.com/5758/31289864222_5e3cca7e72_b.jpg": "",
			"https://farm6.staticflickr.com/5801/30627076673_5a32f3e562_b.jpg": "",
			"https://farm6.staticflickr.com/5538/31319458901_088858d7f1_b.jpg": "",
		},
	},
	{
		"https://www.flickr.com/photos/137385017@N08/albums/with/72157677350305446",
		map[string]string{
			"https://farm6.staticflickr.com/5521/31065043320_cd03a9a448_b.jpg": "",
			"https://farm6.staticflickr.com/5651/31434767515_49f88ee12e_b.jpg": "",
			"https://farm6.staticflickr.com/5750/31434766825_529fd08071_b.jpg": "",
			"https://farm6.staticflickr.com/5811/31319456971_37c8c4708a_b.jpg": "",
			"https://farm6.staticflickr.com/5494/30627074913_b7f810fc26_b.jpg": "",
			"https://farm6.staticflickr.com/5539/31065042720_d76f643b28_b.jpg": "",
			"https://farm6.staticflickr.com/5813/31434765285_94b85d5e8c_b.jpg": "",
			"https://farm6.staticflickr.com/5600/31065044090_eca63bd5a5_b.jpg": "",
			"https://farm6.staticflickr.com/5733/31434764435_350825477e_b.jpg": "",
			"https://farm6.staticflickr.com/5715/30627073573_b86e4b2c22_b.jpg": "",
			"https://farm6.staticflickr.com/5758/31289864222_5e3cca7e72_b.jpg": "",
			"https://farm6.staticflickr.com/5801/30627076673_5a32f3e562_b.jpg": "",
			"https://farm6.staticflickr.com/5538/31319458901_088858d7f1_b.jpg": "",
		},
	},
}

func TestGetFlickrAlbumUrls(t *testing.T) {
	for _, pair := range getFlickrAlbumUrlsTests {
		v, err := getFlickrAlbumUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}
		if !reflect.DeepEqual(v, pair.result) {
			t.Errorf("For %s, expected %s, got %s", pair.value, pair.result, v)
		}
	}
}

var getStreamableUrlsTests = []urlsTestpair{
	{
		"http://streamable.com/41ajc",
		map[string]string{
			"streamablevideo.com/video/mp4/41ajc.mp4": "",
		},
	},
}

func TestGetStreamableUrls(t *testing.T) {
	for _, pair := range getStreamableUrlsTests {
		v, err := getStreamableUrls(pair.value)
		if err != nil {
			t.Errorf("For %v, expected %v, got %v", pair.value, nil, err)
		}

		for expectedLink, expectedName := range pair.result {
			linkFound := false
			for gotLink, _ := range v {
				if strings.Contains(gotLink, expectedLink) {
					linkFound = true
				}
			}
			if !linkFound {
				t.Errorf("For expected %s %s, got %s", expectedLink, expectedName, v)
			}
		}
	}
}
