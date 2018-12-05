package main

import (
	"regexp"
)

const (
	REGEXP_URL_TWITTER                 = `^http(s?):\/\/pbs(-[0-9]+)?\.twimg\.com\/media\/[^\./]+\.(jpg|png)((\:[a-z]+)?)$`
	REGEXP_URL_TWITTER_STATUS          = `^http(s?):\/\/(www\.)?twitter\.com\/([A-Za-z0-9-_\.]+\/status\/|statuses\/|i\/web\/status\/)([0-9]+)$`
	REGEXP_URL_TISTORY                 = `^http(s?):\/\/t[0-9]+\.daumcdn\.net\/cfile\/tistory\/([A-Z0-9]+?)(\?original)?$`
	REGEXP_URL_TISTORY_LEGACY          = `^http(s?):\/\/[a-z0-9]+\.uf\.tistory\.com\/(image|original)\/[A-Z0-9]+$`
	REGEXP_URL_TISTORY_LEGACY_WITH_CDN = `^http(s)?:\/\/[0-9a-z]+.daumcdn.net\/[a-z]+\/[a-zA-Z0-9\.]+\/\?scode=mtistory&fname=http(s?)%3A%2F%2F[a-z0-9]+\.uf\.tistory\.com%2F(image|original)%2F[A-Z0-9]+$`
	REGEXP_URL_GFYCAT                  = `^http(s?):\/\/gfycat\.com\/(gifs\/detail\/)?[A-Za-z]+$`
	REGEXP_URL_INSTAGRAM               = `^http(s?):\/\/(www\.)?instagram\.com\/p\/[^/]+\/(\?[^/]+)?$`
	REGEXP_URL_IMGUR_SINGLE            = `^http(s?):\/\/(i\.)?imgur\.com\/[A-Za-z0-9]+(\.gifv)?$`
	REGEXP_URL_IMGUR_ALBUM             = `^http(s?):\/\/imgur\.com\/(a\/|gallery\/|r\/[^\/]+\/)[A-Za-z0-9]+(#[A-Za-z0-9]+)?$`
	REGEXP_URL_GOOGLEDRIVE             = `^http(s?):\/\/drive\.google\.com\/file\/d\/[^/]+\/view$`
	REGEXP_URL_GOOGLEDRIVE_FOLDER      = `^http(s?):\/\/drive\.google\.com\/(drive\/folders\/|open\?id=)([^/]+)$`
	REGEXP_URL_POSSIBLE_TISTORY_SITE   = `^http(s)?:\/\/[0-9a-zA-Z\.-]+\/(m\/)?(photo\/)?[0-9]+$`
	REGEXP_URL_FLICKR_PHOTO            = `^http(s)?:\/\/(www\.)?flickr\.com\/photos\/([0-9]+)@([A-Z0-9]+)\/([0-9]+)(\/)?(\/in\/album-([0-9]+)(\/)?)?$`
	REGEXP_URL_FLICKR_ALBUM            = `^http(s)?:\/\/(www\.)?flickr\.com\/photos\/(([0-9]+)@([A-Z0-9]+)|[A-Za-z0-9]+)\/(albums\/(with\/)?|(sets\/)?)([0-9]+)(\/)?$`
	REGEXP_URL_FLICKR_ALBUM_SHORT      = `^http(s)?:\/\/((www\.)?flickr\.com\/gp\/[0-9]+@[A-Z0-9]+\/[A-Za-z0-9]+|flic\.kr\/s\/[a-zA-Z0-9]+)$`
	REGEXP_URL_STREAMABLE              = `^http(s?):\/\/(www\.)?streamable\.com\/([0-9a-z]+)$`

	REGEXP_FILENAME = `^^[^/\\:*?"<>|]{1,150}\.[A-Za-z0-9]{2,4}$$`
)

var (
	RegexpUrlTwitter              *regexp.Regexp
	RegexpUrlTwitterStatus        *regexp.Regexp
	RegexpUrlTistory              *regexp.Regexp
	RegexpUrlTistoryLegacy        *regexp.Regexp
	RegexpUrlTistoryLegacyWithCDN *regexp.Regexp
	RegexpUrlGfycat               *regexp.Regexp
	RegexpUrlInstagram            *regexp.Regexp
	RegexpUrlImgurSingle          *regexp.Regexp
	RegexpUrlImgurAlbum           *regexp.Regexp
	RegexpUrlGoogleDrive          *regexp.Regexp
	RegexpUrlGoogleDriveFolder    *regexp.Regexp
	RegexpUrlPossibleTistorySite  *regexp.Regexp
	RegexpUrlFlickrPhoto          *regexp.Regexp
	RegexpUrlFlickrAlbum          *regexp.Regexp
	RegexpUrlFlickrAlbumShort     *regexp.Regexp
	RegexpUrlStreamable           *regexp.Regexp
)

func initRegex() error {
	var err error
	RegexpUrlTwitter, err = regexp.Compile(REGEXP_URL_TWITTER)
	if err != nil {
		return err
	}
	RegexpUrlTwitterStatus, err = regexp.Compile(REGEXP_URL_TWITTER_STATUS)
	if err != nil {
		return err
	}
	RegexpUrlTistory, err = regexp.Compile(REGEXP_URL_TISTORY)
	if err != nil {
		return err
	}
	RegexpUrlTistoryLegacy, err = regexp.Compile(REGEXP_URL_TISTORY_LEGACY)
	if err != nil {
		return err
	}
	RegexpUrlTistoryLegacyWithCDN, err = regexp.Compile(REGEXP_URL_TISTORY_LEGACY_WITH_CDN)
	if err != nil {
		return err
	}
	RegexpUrlGfycat, err = regexp.Compile(REGEXP_URL_GFYCAT)
	if err != nil {
		return err
	}
	RegexpUrlInstagram, err = regexp.Compile(REGEXP_URL_INSTAGRAM)
	if err != nil {
		return err
	}
	RegexpUrlImgurSingle, err = regexp.Compile(REGEXP_URL_IMGUR_SINGLE)
	if err != nil {
		return err
	}
	RegexpUrlImgurAlbum, err = regexp.Compile(REGEXP_URL_IMGUR_ALBUM)
	if err != nil {
		return err
	}
	RegexpUrlGoogleDrive, err = regexp.Compile(REGEXP_URL_GOOGLEDRIVE)
	if err != nil {
		return err
	}
	RegexpUrlGoogleDriveFolder, err = regexp.Compile(REGEXP_URL_GOOGLEDRIVE_FOLDER)
	if err != nil {
		return err
	}
	RegexpUrlPossibleTistorySite, err = regexp.Compile(REGEXP_URL_POSSIBLE_TISTORY_SITE)
	if err != nil {
		return err
	}
	RegexpUrlFlickrPhoto, err = regexp.Compile(REGEXP_URL_FLICKR_PHOTO)
	if err != nil {
		return err
	}
	RegexpUrlFlickrAlbum, err = regexp.Compile(REGEXP_URL_FLICKR_ALBUM)
	if err != nil {
		return err
	}
	RegexpUrlStreamable, err = regexp.Compile(REGEXP_URL_STREAMABLE)
	if err != nil {
		return err
	}
	RegexpUrlFlickrAlbumShort, err = regexp.Compile(REGEXP_URL_FLICKR_ALBUM_SHORT)
	if err != nil {
		return err
	}
	RegexpFilename, err = regexp.Compile(REGEXP_FILENAME)
	if err != nil {
		return err
	}

	return nil
}
