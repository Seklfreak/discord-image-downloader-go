package main

import (
	"regexp"
)

var (
	RegexpUrlTwitter             *regexp.Regexp
	RegexpUrlTwitterStatus       *regexp.Regexp
	RegexpUrlTistory             *regexp.Regexp
	RegexpUrlTistoryWithCDN      *regexp.Regexp
	RegexpUrlGfycat              *regexp.Regexp
	RegexpUrlInstagram           *regexp.Regexp
	RegexpUrlImgurSingle         *regexp.Regexp
	RegexpUrlImgurAlbum          *regexp.Regexp
	RegexpUrlGoogleDrive         *regexp.Regexp
	RegexpUrlGoogleDriveFolder   *regexp.Regexp
	RegexpUrlPossibleTistorySite *regexp.Regexp
	RegexpUrlFlickrPhoto         *regexp.Regexp
	RegexpUrlFlickrAlbum         *regexp.Regexp
	RegexpUrlFlickrAlbumShort    *regexp.Regexp
	RegexpUrlStreamable          *regexp.Regexp
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
	RegexpUrlTistoryWithCDN, err = regexp.Compile(REGEXP_URL_TISTORY_WITH_CDN)
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
