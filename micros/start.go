package function

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	core "github.com/red-gold/telar-core"
	"github.com/red-gold/telar-core/config"
	coreSetting "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/data/mongodb"
	coreUtils "github.com/red-gold/telar-core/utils"
)

const (
	basePath                = "/var/openfaas/secrets/"
	mongoHostSecretKey      = "mongo-host"
	mongoDatabaseSecretKey  = "mongo-database"
	phoneAuthIDSecretKey    = "phone-auth-id"
	phoneAuthTokenSecretKey = "phone-auth-token"
	privateKeySecretKey     = "key"
	publicKeySecretKey      = "key.pub"
	recaptchaSecretKey      = "recaptcha-key"
	refEmailPassSecretKey   = "ref-email-pass"
	payloadSecretKey        = "payload-secret"
)

var secretKeys = []string{mongoHostSecretKey, mongoDatabaseSecretKey,
	phoneAuthIDSecretKey, phoneAuthTokenSecretKey, privateKeySecretKey,
	publicKeySecretKey, recaptchaSecretKey, refEmailPassSecretKey, payloadSecretKey}

// Initiailze micros configurations
func InitConfig() {
	coreConfig := getAllConfiguration()
	core.InitConfigFromData(*coreConfig)
}

// Start run startup operations
func Start(ctx context.Context) (interface{}, error) {
	coreConfig := coreSetting.AppConfig

	switch *coreConfig.DBType {
	case config.DB_MONGO:
		mongoClient, err := mongodb.NewMongoClient(ctx, *coreConfig.MongoDBHost, *coreConfig.Database)
		if err != nil {
			return nil, err
		}
		return mongoClient, nil
	}

	return nil, fmt.Errorf("Please set valid database type in confing file!")
}

// getAllConfigFromFile get all config from files
func getAllConfigFromFile() map[string][]byte {
	filePaths := []string{}
	for _, v := range secretKeys {
		filePaths = append(filePaths, basePath+v)
	}
	return coreUtils.GetFilesContents(filePaths...)
}

// getAllConfiguration get all configuration
func getAllConfiguration() *config.Configuration {
	var newCoreConfig config.Configuration

	loadSecretMode, ok := os.LookupEnv("load_secret_mode")
	if ok {
		log.Printf("[INFO]: Load secret mode information loaded from env.")
		if loadSecretMode == "env" {
			loadSecretsFromEnv(&newCoreConfig)
		}
	} else {
		log.Printf("[INFO]: No secret mode in env. Secrets are loading from file.")
		loadSecretsFromFile(&newCoreConfig)
	}

	// Load from environment //

	appName, ok := os.LookupEnv("app_name")
	if ok {
		newCoreConfig.AppName = &appName
		log.Printf("[INFO]: App Name information loaded from env.")
	}

	queryPrettyURL, ok := os.LookupEnv("query_pretty_url")
	if ok {
		parsedQueryPrettyURL, errParseDebug := strconv.ParseBool(queryPrettyURL)
		if errParseDebug != nil {
			log.Printf("[ERROR]: Query Pretty URL information loading error: %s", errParseDebug.Error())
		}
		newCoreConfig.QueryPrettyURL = &parsedQueryPrettyURL
		log.Printf("[INFO]: Query Pretty URL information loaded from env.")
	}

	debug, ok := os.LookupEnv("debug")
	if ok {
		parsedDebug, errParseDebug := strconv.ParseBool(debug)
		if errParseDebug != nil {
			log.Printf("[ERROR]: Debug information loading error: %s", errParseDebug.Error())
		}
		newCoreConfig.Debug = &parsedDebug
		log.Printf("[INFO]: Debug information loaded from env.")
	}

	gateway, ok := os.LookupEnv("gateway")
	if ok {
		newCoreConfig.Gateway = &gateway
		log.Printf("[INFO]: Gateway information loaded from env.")
	}

	internalGateway, ok := os.LookupEnv("internal_gateway")
	if ok {
		newCoreConfig.InternalGateway = &internalGateway
		log.Printf("[INFO]: Internal gateway information loaded from env. | %s |", internalGateway)
	}

	webDomain, ok := os.LookupEnv("web_domain")
	if ok {
		newCoreConfig.WebDomain = &webDomain
		log.Printf("[INFO]: Web domain information loaded from env.")
	}

	orgName, ok := os.LookupEnv("org_name")
	if ok {
		newCoreConfig.OrgName = &orgName
		log.Printf("[INFO]: Organization Name information loaded from env.")
	}

	orgAvatar, ok := os.LookupEnv("org_avatar")
	if ok {
		newCoreConfig.OrgAvatar = &orgAvatar
		log.Printf("[INFO]: Organization Avatar information loaded from env.")
	}

	server, ok := os.LookupEnv("server")
	if ok {
		newCoreConfig.Server = &server
		log.Printf("[INFO]: Server information loaded from env.")
	}

	recaptchaSiteKey, ok := os.LookupEnv("recaptcha_site_key")
	if ok {
		newCoreConfig.RecaptchaSiteKey = &recaptchaSiteKey
		log.Printf("[INFO]: Recaptcha site key information loaded from env.")
	}

	origin, ok := os.LookupEnv("origin")
	if ok {
		newCoreConfig.Origin = &origin
		log.Printf("[INFO]: Origin information loaded from env.")
	}

	headerCookieName, ok := os.LookupEnv("header_cookie_name")
	if ok {
		newCoreConfig.HeaderCookieName = &headerCookieName
		log.Printf("[INFO]: Header cookie name information loaded from env.")
	}

	payloadCookieName, ok := os.LookupEnv("payload_cookie_name")
	if ok {
		newCoreConfig.PayloadCookieName = &payloadCookieName
		log.Printf("[INFO]: Payload cookie name information loaded from env.")
	}

	signatureCookieName, ok := os.LookupEnv("signature_cookie_name")
	if ok {
		newCoreConfig.SignatureCookieName = &signatureCookieName
		log.Printf("[INFO]: Signature cookie name information loaded from env.")
	}

	baseRoute, ok := os.LookupEnv("base_route_domain")
	if ok {
		newCoreConfig.BaseRoute = &baseRoute
		log.Printf("[INFO]: Base route information loaded from env.")
	}

	smtpEmail, ok := os.LookupEnv("smtp_email")
	if ok {
		newCoreConfig.SmtpEmail = &smtpEmail
		log.Printf("[INFO]: SMTP Email information loaded from env.")
	}

	refEmail, ok := os.LookupEnv("ref_email")
	if ok {
		newCoreConfig.RefEmail = &refEmail
		log.Printf("[INFO]: Reference Email information loaded from env.")
	}

	phoneSourceNumebr, ok := os.LookupEnv("phone_source_number")
	if ok {
		newCoreConfig.PhoneSourceNumber = &phoneSourceNumebr
		log.Printf("[INFO]: Phone Source Number information loaded from env.")
	}

	phoneAuthToken, ok := os.LookupEnv("phone_auth_token")
	if ok {
		newCoreConfig.PhoneAuthToken = &phoneAuthToken
		log.Printf("[INFO]: Phone Auth Token  information loaded from env.")
	}

	phoneAuthId, ok := os.LookupEnv("phone_auth_id")
	if ok {
		newCoreConfig.PhoneAuthId = &phoneAuthId
		log.Printf("[INFO]: Phone Auth Id  information loaded from env.")
	}

	dbType, ok := os.LookupEnv("db_type")
	if ok {
		newCoreConfig.DBType = &dbType
		log.Printf("[INFO]: Database type information loaded from env.")
	}
	return &newCoreConfig
}

// loadSecretsFromFile Load secrets from file
func loadSecretsFromFile(newCoreConfig *config.Configuration) {
	filesConfig := getAllConfigFromFile()

	// Load from files //
	if filesConfig[basePath+payloadSecretKey] != nil {
		payloadSecret := string(filesConfig[basePath+payloadSecretKey])
		newCoreConfig.PayloadSecret = &payloadSecret
		log.Printf("[INFO]: Payload secret information loaded from env.")
	}

	if filesConfig[basePath+publicKeySecretKey] != nil {
		publicKey := string(filesConfig[basePath+publicKeySecretKey])
		newCoreConfig.PublicKey = &publicKey
		log.Printf("[INFO]: Public key information loaded from env.")
	}

	if filesConfig[basePath+privateKeySecretKey] != nil {
		privateKey := string(filesConfig[basePath+privateKeySecretKey])
		newCoreConfig.PrivateKey = &privateKey
		log.Printf("[INFO]: Private key information loaded from env.")
	}

	if filesConfig[basePath+recaptchaSecretKey] != nil {
		recaptchaKey := string(filesConfig[basePath+recaptchaSecretKey])
		newCoreConfig.RecaptchaKey = &recaptchaKey
		log.Printf("[INFO]: Recaptcha key information loaded from env.")
	}

	if filesConfig[basePath+mongoHostSecretKey] != nil {
		mongoDBHost := string(filesConfig[basePath+mongoHostSecretKey])
		newCoreConfig.MongoDBHost = &mongoDBHost
		log.Printf("[INFO]: MongoDB host information loaded from env.")
	}

	if filesConfig[basePath+mongoDatabaseSecretKey] != nil {
		mongoDB := string(filesConfig[basePath+mongoDatabaseSecretKey])
		newCoreConfig.Database = &mongoDB
		log.Printf("[INFO]: Database name information loaded from env.")
	}

	if filesConfig[basePath+refEmailPassSecretKey] != nil {
		refEmailPass := string(filesConfig[basePath+refEmailPassSecretKey])
		newCoreConfig.RefEmailPass = &refEmailPass
		log.Printf("[INFO]: Ref email password information loaded from env.")
	}
}

// loadSecretsFromEnv Load secrets from environment variables
func loadSecretsFromEnv(newCoreConfig *config.Configuration) {

	payloadSecret, ok := os.LookupEnv("payload_secret")
	if ok {
		payloadSecret = decodeBase64(payloadSecret)
		newCoreConfig.PayloadSecret = &payloadSecret
		log.Printf("[INFO]: Payload secret information loaded from env.")
	}

	publicKey, ok := os.LookupEnv("key_pub")
	if ok {
		publicKey = decodeBase64(publicKey)
		newCoreConfig.PublicKey = &publicKey
		log.Printf("[INFO]: Public key information loaded from env.")
	}

	privateKey, ok := os.LookupEnv("key")
	if ok {
		privateKey = decodeBase64(privateKey)
		newCoreConfig.PrivateKey = &privateKey
		log.Printf("[INFO]: Private key information loaded from env.")
	}

	recaptchaKey, ok := os.LookupEnv("recaptcha_key")
	if ok {
		recaptchaKey = decodeBase64(recaptchaKey)
		newCoreConfig.RecaptchaKey = &recaptchaKey
		log.Printf("[INFO]: Recaptcha key information loaded from env.")
	}

	mongoDBHost, ok := os.LookupEnv("mongo_host")
	if ok {
		mongoDBHost = decodeBase64(mongoDBHost)
		newCoreConfig.MongoDBHost = &mongoDBHost
		log.Printf("[INFO]: MongoDB host information loaded from env.")
	}

	mongoDB, ok := os.LookupEnv("mongo_database")
	if ok {
		mongoDB = decodeBase64(mongoDB)
		newCoreConfig.Database = &mongoDB
		log.Printf("[INFO]: Database name information loaded from env.")
	}

	refEmailPass, ok := os.LookupEnv("ref_email_pass")
	if ok {
		refEmailPass = decodeBase64(refEmailPass)
		newCoreConfig.RefEmailPass = &refEmailPass
		log.Printf("[INFO]: Ref email password information loaded from env.")
	}
}

// decodeBase64 Decode base64 string
func decodeBase64(encodedString string) string {
	base64Value, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		fmt.Println("[ERROR] decode secret base64 value with value:  ", encodedString, " - ", err.Error())
		panic(err)
	}
	return string(base64Value)
}
