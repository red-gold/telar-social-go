package config

type (
	// appError struct {
	// 	Error      string `json:"error"`
	// 	Message    string `json:"message"`
	// 	HttpStatus int    `json:"status"`
	// }
	// errorResource struct {
	// 	Data appError `json:"data"`
	// }
	Configuration struct {
		Server              *string
		Gateway             *string
		MongoDBHost         *string
		MongoPwdPath        *string
		Database            *string
		PublicKeyPath       *string
		RecaptchaKeyPath    *string
		RecaptchaSiteKey    *string
		HeaderCookieName    *string
		PayloadCookieName   *string
		SignatureCookieName *string
		SmtpEmail           *string
		RefEmail            *string
		RefEmailPassPath    *string
		Origin              *string
		PrivateKeyPath      *string
		AppName             *string
		PhoneSourceNumber   *string
		PhoneAuthTokenPath  *string
		PhoneAuthIdPath     *string
		OrgName             *string
		OrgAvatar           *string
		WebDomain           *string
		DBType              *string
		QueryPrettyURL      *bool
	}
)

const (
	DB_INMEMORY = "inmemory"
	DB_MONGO    = "mongo"
	DB_SQLITE   = "sqlite"
	DB_MYSQL    = "mysql"
)

// AppConfig holds the Configuration values from app-config.yml file
var AppConfig Configuration
