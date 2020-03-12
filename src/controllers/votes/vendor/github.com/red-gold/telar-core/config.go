package core

import (
	"log"
	"os"
	"strconv"

	"github.com/red-gold/telar-core/config"
)

// Initialize AppConfig
func InitConfig() {

	// Load config from environment values if exists
	loadConfigFromEnvironment()
}

// Load config from environment
func loadConfigFromEnvironment() {

	appName, ok := os.LookupEnv("app_name")
	if ok {
		config.AppConfig.AppName = &appName
		log.Printf("[INFO]: App Name information loaded from env.")
	}

	queryPrettyURL, ok := os.LookupEnv("query_pretty_url")
	if ok {
		parsedQueryPrettyURL, errParseDebug := strconv.ParseBool(queryPrettyURL)
		if errParseDebug != nil {
			log.Printf("[ERROR]: Query Pretty URL information loading error: %s", errParseDebug.Error())
		}
		config.AppConfig.QueryPrettyURL = &parsedQueryPrettyURL
		log.Printf("[INFO]: Query Pretty URL information loaded from env.")
	}
	gateway, ok := os.LookupEnv("gateway")
	if ok {
		config.AppConfig.Gateway = &gateway
		log.Printf("[INFO]: Gateway information loaded from env.")
	}

	webDomain, ok := os.LookupEnv("web_domain")
	if ok {
		config.AppConfig.WebDomain = &webDomain
		log.Printf("[INFO]: Web domain information loaded from env.")
	}

	orgName, ok := os.LookupEnv("org_name")
	if ok {
		config.AppConfig.OrgName = &orgName
		log.Printf("[INFO]: Organization Name information loaded from env.")
	}

	orgAvatar, ok := os.LookupEnv("org_avatar")
	if ok {
		config.AppConfig.OrgAvatar = &orgAvatar
		log.Printf("[INFO]: Organization Avatar information loaded from env.")
	}

	server, ok := os.LookupEnv("server")
	if ok {
		config.AppConfig.Server = &server
		log.Printf("[INFO]: Server information loaded from env.")
	}

	publicKeyPath, ok := os.LookupEnv("public_key_path")
	if ok {
		config.AppConfig.PublicKeyPath = &publicKeyPath
		log.Printf("[INFO]: Public key path information loaded from env.")
	}

	privateKeyPath, ok := os.LookupEnv("private_key_path")
	if ok {
		config.AppConfig.PrivateKeyPath = &privateKeyPath
		log.Printf("[INFO]: PrivateKeyPath information loaded from env.")
	}

	recaptchaKeyPath, ok := os.LookupEnv("recaptcha_key_path")
	if ok {
		config.AppConfig.RecaptchaKeyPath = &recaptchaKeyPath
		log.Printf("[INFO]: Recaptcha key path information loaded from env.")
	}

	recaptchaSiteKey, ok := os.LookupEnv("recaptcha_site_key")
	if ok {
		config.AppConfig.RecaptchaSiteKey = &recaptchaSiteKey
		log.Printf("[INFO]: Recaptcha site key information loaded from env.")
	}

	origin, ok := os.LookupEnv("origin")
	if ok {
		config.AppConfig.Origin = &origin
		log.Printf("[INFO]: Origin information loaded from env.")
	}

	headerCookieName, ok := os.LookupEnv("header_cookie_name")
	if ok {
		config.AppConfig.HeaderCookieName = &headerCookieName
		log.Printf("[INFO]: Header cookie name information loaded from env.")
	}

	payloadCookieName, ok := os.LookupEnv("payload_cookie_name")
	if ok {
		config.AppConfig.PayloadCookieName = &payloadCookieName
		log.Printf("[INFO]: Payload cookie name information loaded from env.")
	}

	signatureCookieName, ok := os.LookupEnv("signature_cookie_name")
	if ok {
		config.AppConfig.SignatureCookieName = &signatureCookieName
		log.Printf("[INFO]: Signature cookie name information loaded from env.")
	}

	mongodbHost, ok := os.LookupEnv("mongo_host")
	if ok {
		config.AppConfig.MongoDBHost = &mongodbHost
		log.Printf("[INFO]: MongoDB host information loaded from env.")
	}

	mongodbPwdPath, ok := os.LookupEnv("mongo_pwd_path")
	if ok {
		config.AppConfig.MongoPwdPath = &mongodbPwdPath
		log.Printf("[INFO]: MongoDB password path information loaded from env.")
	}

	database, ok := os.LookupEnv("mongo_database")
	if ok {
		config.AppConfig.Database = &database
		log.Printf("[INFO]: MongoDB database information loaded from env.")
	}

	smtpEmail, ok := os.LookupEnv("smtp_email")
	if ok {
		config.AppConfig.SmtpEmail = &smtpEmail
		log.Printf("[INFO]: SMTP Email information loaded from env.")
	}

	refEmail, ok := os.LookupEnv("ref_email")
	if ok {
		config.AppConfig.RefEmail = &refEmail
		log.Printf("[INFO]: Reference Email information loaded from env.")
	}

	phoneSourceNumebr, ok := os.LookupEnv("phone_source_number")
	if ok {
		config.AppConfig.PhoneSourceNumber = &phoneSourceNumebr
		log.Printf("[INFO]: Phone Source Number information loaded from env.")
	}

	phoneAuthToken, ok := os.LookupEnv("phone_auth_token_path")
	if ok {
		config.AppConfig.PhoneAuthTokenPath = &phoneAuthToken
		log.Printf("[INFO]: Phone Auth Token Path information loaded from env.")
	}

	phoneAuthId, ok := os.LookupEnv("phone_auth_id_path")
	if ok {
		config.AppConfig.PhoneAuthIdPath = &phoneAuthId
		log.Printf("[INFO]: Phone Auth Id Path information loaded from env.")
	}

	refEmailPassPath, ok := os.LookupEnv("ref_email_pass_path")
	if ok {
		config.AppConfig.RefEmailPassPath = &refEmailPassPath
		log.Printf("[INFO]: Reference Email Password Path information loaded from env.")
	}

	dbType, ok := os.LookupEnv("db_type")
	if ok {
		config.AppConfig.DBType = &dbType
		log.Printf("[INFO]: Database type information loaded from env.")
	}
}
