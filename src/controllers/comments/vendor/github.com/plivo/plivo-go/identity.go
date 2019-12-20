package plivo

import (
	"errors"
	"os"
	"path/filepath"
)

type IdentityService struct {
	client         *Client
	IdentitySalutationType identitySalutationRegistry
	IdentityNumberType     identityNumberTypeRegistry
	IdentityProofType      identityProofTypeRegistry
	IdentityFileType       identityFileTypeRegistry
}
var IdentitySalutationType = &identitySalutationRegistry{
	MR:   "Mr",
	MRS: "Mrs",
}

var IdentityNumberType = &identityNumberTypeRegistry{
	LOCAL:    "local",
	NATIONAL: "national",
	MOBILE:   "mobile",
	TOLLFREE: "tollfree",
}

var IdentityProofType =  &identityProofTypeRegistry{
	NATIONALID:   "national_id",
	PASSPORT: "passport",
	BUSINESSID: "business_id",
	NIF: "NIF",
	NIE :"NIE",
	DNI:"DNI",
}

var IdentityFileType = &identityFileTypeRegistry{
	PNG:   "png",
	JPG: "jpg",
	PDF: "pdf",
}

type identitySalutationRegistry struct {
	MR   string
	MRS string
}

type identityNumberTypeRegistry struct {
	LOCAL    string
	NATIONAL string
	MOBILE   string
	TOLLFREE string
}

type identityProofTypeRegistry struct {
	NATIONALID   string
	PASSPORT string
	BUSINESSID   string
	NIF string
	NIE   string
	DNI string
}

type identityFileTypeRegistry struct {
	PNG   string
	JPG string
	PDF   string
}

type IdentityCreateParams struct {
	PhoneNumberCountry string `json:"phone_number_country,omitempty" url:"phone_number_country,omitempty"`
	NumberType string `json:"number_type,omitempty" url:"number_type,omitempty"`
	Salutation string `json:"salutation,omitempty" url:"salutation,omitempty"`
	FirstName string `json:"first_name,omitempty" url:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty" url:"last_name,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty" url:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty" url:"address_line2,omitempty"`
	City string `json:"city,omitempty" url:"city,omitempty"`
	Region string `json:"region,omitempty" url:"region,omitempty"`
	PostalCode string `json:"postal_code,omitempty" url:"postal_code,omitempty"`
	CountryIso string `json:"country_iso,omitempty" url:"country_iso,omitempty"`
	AddressProofType string `json:"proof_type,omitempty" url:"proof_type,omitempty"`
	IdNumber string `json:"id_number,omitempty" url:"id_number,omitempty"`
	Nationality string `json:"nationality,omitempty" url:"nationality,omitempty"`


	CallBackUrl string `json:"callback_url,omitempty" url:"callback_url,omitempty"`
	Alias string `json:"alias,omitempty" url:"alias,omitempty"`
	File interface{} `json:"file,omitempty" url:"file,omitempty"`
	IdNationality string `json:"id_nationality,omitempty" url:"id_nationality,omitempty"`
	BirthPlace string `json:"birth_place,omitempty" url:"birth_place,omitempty"`
	BirthDate string `json:"birth_date,omitempty" url:"birth_date,omitempty"`
	IdIssueDate string `json:"id_issue_date,omitempty" url:"id_issue_date,omitempty"`
	BusinessName string `json:"business_name,omitempty" url:"business_name,omitempty"`
	FiscalIdentificationCode string `json:"fiscal_identification_code,omitempty" url:"fiscal_identification_code,omitempty"`
	StreetCode string `json:"street_code,omitempty" url:"street_code,omitempty"`
	MunicipalCode string `json:"municipal_code,omitempty" url:"municipal_code,omitempty"`
	AutoCorrectAddress string `json:"autoCorrectAddress,omitempty" url:"autoCorrectAddress,omitempty"`

}

type IdentityUpdateParams struct {
	PhoneNumberCountry string `json:"phone_number_country,omitempty" url:"phone_number_country,omitempty"`
	NumberType string `json:"number_type,omitempty" url:"number_type,omitempty"`
	Salutation string `json:"salutation,omitempty" url:"salutation,omitempty"`
	FirstName string `json:"first_name,omitempty" url:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty" url:"last_name,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty" url:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty" url:"address_line2,omitempty"`
	City string `json:"city,omitempty" url:"city,omitempty"`
	Region string `json:"region,omitempty" url:"region,omitempty"`
	PostalCode string `json:"postal_code,omitempty" url:"postal_code,omitempty"`
	CountryIso string `json:"country_iso,omitempty" url:"country_iso,omitempty"`
	AddressProofType string `json:"proof_type,omitempty" url:"proof_type,omitempty"`
	IdNumber string `json:"id_number,omitempty" url:"id_number,omitempty"`
	Nationality string `json:"nationality,omitempty" url:"nationality,omitempty"`


	CallBackUrl string `json:"callback_url,omitempty" url:"callback_url,omitempty"`
	Alias string `json:"alias,omitempty" url:"alias,omitempty"`
	File interface{} `json:"file,omitempty" url:"file,omitempty"`
	IDNationality string `json:"id_nationality,omitempty" url:"id_nationality,omitempty"`
	BirthPlace string `json:"birth_place,omitempty" url:"birth_place,omitempty"`
	BirthDate string `json:"birth_date,omitempty" url:"birth_date,omitempty"`
	IdIssueDate string `json:"id_issue_date,omitempty" url:"id_issue_date,omitempty"`
	BusinessName string `json:"business_name,omitempty" url:"business_name,omitempty"`
	FiscalIdentificationCode string `json:"fiscal_identification_code,omitempty" url:"fiscal_identification_code,omitempty"`
	StreetCode string `json:"street_code,omitempty" url:"street_code,omitempty"`
	MunicipalCode string `json:"municipal_code,omitempty" url:"municipal_code,omitempty"`
	AutoCorrectAddress string `json:"autoCorrectAddress,omitempty" url:"autoCorrectAddress,omitempty"`

}

type Identity struct {
	Account string `json:"account,omitempty" url:"account,omitempty"`
	ID string `json:"id,omitempty" url:"id,omitempty"`
	CountryIso string `json:"country_iso,omitempty" url:"country_iso,omitempty"`
	Alias string `json:"alias,omitempty" url:"alias,omitempty"`
	Salutation string `json:"salutation,omitempty" url:"salutation,omitempty"`
	FirstName string `json:"first_name,omitempty" url:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty" url:"last_name,omitempty"`
	BirthPlace string `json:"birth_place,omitempty" url:"birth_place,omitempty"`
	BirthDate string `json:"birth_date,omitempty" url:"birth_date,omitempty"`
	Nationality string `json:"nationality,omitempty" url:"nationality,omitempty"`
	IDNationality string `json:"id_nationality,omitempty" url:"id_nationality,omitempty"`
	IdIssueDate string `json:"id_issue_date,omitempty" url:"id_issue_date,omitempty"`
	BusinessName string `json:"business_name,omitempty" url:"business_name,omitempty"`
	IDType string `json:"id_type,omitempty" url:"id_type,omitempty"`
	IdNumber string `json:"id_number,omitempty" url:"id_number,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty" url:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty" url:"address_line2,omitempty"`
	City string `json:"city,omitempty" url:"city,omitempty"`
	Region string `json:"region,omitempty" url:"region,omitempty"`
	PostalCode string `json:"postal_code,omitempty" url:"postal_code,omitempty"`
	FiscalIdentificationCode string `json:"fiscal_identification_code,omitempty" url:"fiscal_identification_code,omitempty"`
	StreetCode string `json:"street_code,omitempty" url:"street_code,omitempty"`
	MunicipalCode string `json:"municipal_code,omitempty" url:"municipal_code,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty" url:"validation_status,omitempty"`
	VerificationStatus string `json:"verification_status,omitempty" url:"verification_status,omitempty"`
	Subaccount string `json:"subaccount,omitempty" url:"subaccount,omitempty"`
	Url string `json:"url,omitempty" url:"url,omitempty"`
	DocumentDetails map[string]interface{}  `json:"document_details,omitempty" url:"document_details,omitempty"`

}


type IdentityListParams struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

// Stores response for Create call
type IdentityCreateResponseBody struct {
	Message string `json:"message" url:"message"`
	ApiID   string `json:"api_id" url:"api_id"`
}

type IdentityList struct {
	BaseListResponse
	Objects []Address `json:"objects" url:"objects"`
}

type IdentityUpdateResponse BaseResponse

func identityCreateParamsNotNull(params IdentityCreateParams) bool{
	if len(params.PhoneNumberCountry) == 0 || len(params.NumberType)  == 0 || len(params.Salutation)  == 0 || len(params.FirstName)  == 0 ||
			len(params.LastName)  == 0 || len(params.AddressLine1)  == 0 || len(params.AddressLine2)  == 0 || len(params.City)  == 0 ||
			len(params.Region)  == 0 || len(params.PostalCode)  == 0 || len(params.CountryIso)  == 0 || len(params.AddressProofType)  == 0  ||
		len(params.IdNumber) == 0 || len(params.Nationality) == 0{
		return false
	}
	return true
}

func (service *IdentityService) Create(params IdentityCreateParams)(response *IdentityCreateResponseBody, err error){
	if !identityCreateParamsNotNull(params) {
		err = errors.New("PhoneNumberCountry, numberType, salutation, firstName, lastName," +
			"addressLine1, addressLine2, city, region, postalCode, countryIso, proofType, idNumber and nationality must not be empty")
		return
	}
	err = IsValidCountryParams(params.CountryIso,params.FiscalIdentificationCode,params.StreetCode,params.MunicipalCode)
	if err != nil {
		return
	}
	if (params.File !=  nil){
		err = IsValidFile(params.File.(string))
		if err != nil {
			return
		}
		targetFilePath := filepath.Join(params.File.(string), "./")
		file, _ := os.Stat(targetFilePath)
		params.File = file
	}
	request, err := service.client.NewRequest("POST", params, "Verification/Identity")
	if err != nil {
		return
	}
	response = &IdentityCreateResponseBody{}
	err = service.client.ExecuteRequest(request, response)
	return
}

func (service *IdentityService) List(params IdentityListParams) (response *IdentityList, err error) {
	request, err := service.client.NewRequest("GET", params, "Verification/Identity/")
	if err != nil {
		return
	}
	response = &IdentityList{}
	err = service.client.ExecuteRequest(request, response)
	return
}

func (service *IdentityService) Get(documentID string) (response *Identity, err error) {
	request, err := service.client.NewRequest("GET", nil, "Verification/Identity/%s",documentID)
	if err != nil {
		return
	}
	response = &Identity{}
	err = service.client.ExecuteRequest(request, response)
	return
}

func (service *IdentityService)  Delete(documentID string) (err error) {
	request, err := service.client.NewRequest("DELETE", nil, "Verification/Identity/%s",documentID)
	if err != nil {
		return
	}
	err = service.client.ExecuteRequest(request, nil)
	return
}

func (service *IdentityService)  Update(documentID string, params IdentityUpdateParams) (response *IdentityUpdateResponse, err error) {
	err = IsValidCountryParams(params.CountryIso,params.FiscalIdentificationCode,params.StreetCode,params.MunicipalCode)
	if err != nil {
		return
	}
	if (params.File !=  nil){
		err = IsValidFile(params.File.(string))
		if err != nil {
			return
		}
		targetFilePath := filepath.Join(params.File.(string), "./")
		file, _ := os.Stat(targetFilePath)
		params.File = file
	}
	request, err := service.client.NewRequest("POST", params, "Verification/Identity/%s",documentID)
	if err != nil {
		return
	}
	response = &IdentityUpdateResponse{}
	err = service.client.ExecuteRequest(request, response)
	return
}