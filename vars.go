package main

const (
	VERSION                          = "1.31.2"
	DATABASE_DIR                     = "database"
	RELEASE_URL                      = "https://github.com/Seklfreak/discord-image-downloader-go/releases/latest"
	RELEASE_API_URL                  = "https://api.github.com/repos/Seklfreak/discord-image-downloader-go/releases/latest"
	IMGUR_CLIENT_ID                  = "a39473314df3f59"
	REGEXP_URL_TWITTER               = `^http(s?):\/\/pbs(-[0-9]+)?\.twimg\.com\/media\/[^\./]+\.(jpg|png)((\:[a-z]+)?)$`
	REGEXP_URL_TWITTER_STATUS        = `^http(s?):\/\/(www\.)?twitter\.com\/([A-Za-z0-9-_\.]+\/status\/|statuses\/|i\/web\/status\/)([0-9]+)$`
	REGEXP_URL_TISTORY               = `^http(s?):\/\/[a-z0-9]+\.uf\.tistory\.com\/(image|original)\/[A-Z0-9]+$`
	REGEXP_URL_TISTORY_WITH_CDN      = `^http(s)?:\/\/[0-9a-z]+.daumcdn.net\/[a-z]+\/[a-zA-Z0-9\.]+\/\?scode=mtistory&fname=http(s?)%3A%2F%2F[a-z0-9]+\.uf\.tistory\.com%2F(image|original)%2F[A-Z0-9]+$`
	REGEXP_URL_GFYCAT                = `^http(s?):\/\/gfycat\.com\/(gifs\/detail\/)?[A-Za-z]+$`
	REGEXP_URL_INSTAGRAM             = `^http(s?):\/\/(www\.)?instagram\.com\/p\/[^/]+\/(\?[^/]+)?$`
	REGEXP_URL_IMGUR_SINGLE          = `^http(s?):\/\/(i\.)?imgur\.com\/[A-Za-z0-9]+(\.gifv)?$`
	REGEXP_URL_IMGUR_ALBUM           = `^http(s?):\/\/imgur\.com\/(a\/|gallery\/|r\/[^\/]+\/)[A-Za-z0-9]+(#[A-Za-z0-9]+)?$`
	REGEXP_URL_GOOGLEDRIVE           = `^http(s?):\/\/drive\.google\.com\/file\/d\/[^/]+\/view$`
	REGEXP_URL_GOOGLEDRIVE_FOLDER    = `^http(s?):\/\/drive\.google\.com\/(drive\/folders\/|open\?id=)([^/]+)$`
	REGEXP_URL_POSSIBLE_TISTORY_SITE = `^http(s)?:\/\/[0-9a-zA-Z\.-]+\/(m\/)?(photo\/)?[0-9]+$`
	REGEXP_URL_FLICKR_PHOTO          = `^http(s)?:\/\/(www\.)?flickr\.com\/photos\/([0-9]+)@([A-Z0-9]+)\/([0-9]+)(\/)?(\/in\/album-([0-9]+)(\/)?)?$`
	REGEXP_URL_FLICKR_ALBUM          = `^http(s)?:\/\/(www\.)?flickr\.com\/photos\/(([0-9]+)@([A-Z0-9]+)|[A-Za-z0-9]+)\/(albums\/(with\/)?|(sets\/)?)([0-9]+)(\/)?$`
	REGEXP_URL_FLICKR_ALBUM_SHORT    = `^http(s)?:\/\/((www\.)?flickr\.com\/gp\/[0-9]+@[A-Z0-9]+\/[A-Za-z0-9]+|flic\.kr\/s\/[a-zA-Z0-9]+)$`
	REGEXP_URL_STREAMABLE            = `^http(s?):\/\/(www\.)?streamable\.com\/([0-9a-z]+)$`

	REGEXP_FILENAME = `^^[^/\\:*?"<>|]{1,150}\.[A-Za-z0-9]{2,4}$$`

	DEFAULT_CONFIG_FILE = "config.ini"
)
