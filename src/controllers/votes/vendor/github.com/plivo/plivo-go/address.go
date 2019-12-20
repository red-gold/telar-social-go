package plivo

import (
	"errors"
	"os"
	"path/filepath"
)

type AddressService struct {
	client         *Client
	SalutationType addressSalutationRegistry
	NumberType     addressNumberTypeRegistry
	ProofType      addressProofTypeRegistry
	FileType       addressFileTypeRegistry
}

type AddressCreateParams struct {
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
	CallBackUrl string `json:"callback_url,omitempty" url:"callback_url,omitempty"`
	Alias string `json:"alias,omitempty" url:"alias,omitempty"`
	File interface{} `json:"file,omitempty" url:"file,omitempty"`
	AutoCorrectAddress string `json:"autoCorrectAddress,omitempty" url:"autoCorrectAddress,omitempty"`
	AddressProofType string `json:"proof_type,omitempty" url:"proof_type,omitempty"`
	IdNumber string `json:"id_number,omitempty" url:"id_number,omitempty"`
	FiscalIdentificationCode string `json:"fiscal_identification_code,omitempty" url:"fiscal_identification_code,omitempty"`
	StreetCode string `json:"street_code,omitempty" url:"street_code,omitempty"`
	MunicipalCode string `json:"municipal_code,omitempty" url:"municipal_code,omitempty"`
}

type Address struct {
	Account string `json:"account,omitempty" url:"account,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty" url:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty" url:"address_line2,omitempty"`
	Alias string `json:"alias,omitempty" url:"alias,omitempty"`
	StreetCode string `json:"street_code,omitempty" url:"street_code,omitempty"`
	City string `json:"city,omitempty" url:"city,omitempty"`
	CountryIso string `json:"country_iso,omitempty" url:"country_iso,omitempty"`
	DocumentDetails map[string]interface{}  `json:"document_details,omitempty" url:"document_details,omitempty"`
	FirstName string `json:"first_name,omitempty" url:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty" url:"last_name,omitempty"`
	ID string `json:"id,omitempty" url:"id,omitempty"`
	ApiID        string `json:"api_id,omitempty" url:"api_id,omitempty"`
	Region string `json:"region,omitempty" url:"region,omitempty"`
	PostalCode string `json:"postal_code,omitempty" url:"postal_code,omitempty"`
	Salutation string `json:"salutation,omitempty" url:"salutation,omitempty"`
	FiscalIdentificationCode string `json:"fiscal_identification_code,omitempty" url:"fiscal_identification_code,omitempty"`
	Url string `json:"url,omitempty" url:"url,omitempty"`
	Subaccount string `json:"subaccount,omitempty" url:"subaccount,omitempty"`
	MunicipalCode string `json:"municipal_code,omitempty" url:"municipal_code,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty" url:"validation_status,omitempty"`
	AddressProofType string `json:"proof_type,omitempty" url:"proof_type,omitempty"`
	VerificationStatus string `json:"verification_status,omitempty" url:"verification_status,omitempty"`
}

type AddressUpdateParams struct {

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
	CallBackUrl string `json:"callback_url,omitempty" url:"callback_url,omitempty"`
	Alias string `json:"alias,omitempty" url:"alias,omitempty"`
	File interface{}  `json:"file,omitempty" url:"file,omitempty"`
	AutoCorrectAddress bool `json:"autoCorrectAddress,omitempty" url:"autoCorrectAddress,omitempty"`
	AddressProofType string `json:"proof_type,omitempty" url:"proof_type,omitempty"`
	FiscalIdentificationCode string `json:"fiscal_identification_code,omitempty" url:"fiscal_identification_code,omitempty"`
	StreetCode string `json:"street_code,omitempty" url:"street_code,omitempty"`
	MunicipalCode string `json:"municipal_code,omitempty" url:"municipal_code,omitempty"`

}

var AddressSalutationType = &addressSalutationRegistry{
	MR:   "Mr",
	MRS: "Mrs",
}

var AddressNumberType = &addressNumberTypeRegistry{
	LOCAL:    "local",
	NATIONAL: "national",
	MOBILE:   "mobile",
	TOLLFREE: "tollfree",
}

var AddressProofType =  &addressProofTypeRegistry{
	NATIONALID:   "national_id",
	PASSPORT: "passport",
	BUSINESSID: "business_id",
	NIF: "NIF",
	NIE :"NIE",
	DNI:"DNI",
}

var AddressFileType = &addressFileTypeRegistry{
	PNG:   "png",
	JPG: "jpg",
	PDF: "pdf",
}

type addressSalutationRegistry struct {
	MR   string
	MRS string
}

type addressNumberTypeRegistry struct {
	LOCAL    string
	NATIONAL string
	MOBILE   string
	TOLLFREE string
}

type addressProofTypeRegistry struct {
	NATIONALID   string
	PASSPORT string
	BUSINESSID   string
	NIF string
	NIE   string
	DNI string
}

type addressFileTypeRegistry struct {
	PNG   string
	JPG string
	PDF   string
}

type AddressListParams struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}


// Stores response for Create call
type AddressCreateResponseBody struct {
	Message string `json:"message" url:"message"`
	ApiID   string `json:"api_id" url:"api_id"`
}

type AddressList struct {
	BaseListResponse
	Objects []Address `json:"objects" url:"objects"`
}
type AddressUpdateResponse BaseResponse

func addressCreateParamsNotNull(params AddressCreateParams) bool{
	if len(params.PhoneNumberCountry) == 0 || len(params.NumberType)  == 0 || len(params.Salutation)  == 0 ||
		len(params.FirstName)  == 0 || len(params.LastName)  == 0 || len(params.AddressLine1)  == 0 ||
		len(params.AddressLine2)  == 0 || len(params.City)  == 0 || len(params.Region)  == 0 || len(params.PostalCode)  == 0 ||
		len(params.CountryIso)  == 0 {
		return false
	}
	return true
}



func (service *AddressService) Create(params AddressCreateParams)(response *AddressCreateResponseBody, err error){
	if !addressCreateParamsNotNull(params) {
		err = errors.New("PhoneNumberCountry, numberType, salutation, firstName, lastName," +
			"addressLine1, addressLine2, city, region, postalCode and countryIso must not be empty")
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
	request, err := service.client.NewRequest("POST", params, "Verification/Address")
	if err != nil {
		return
	}
	response = &AddressCreateResponseBody{}
	err = service.client.ExecuteRequest(request, response)
	return
}

func (service *AddressService) List(params AddressListParams) (response *AddressList, err error) {
	request, err := service.client.NewRequest("GET", params, "Verification/Address/")
	if err != nil {
		return
	}
	response = &AddressList{}
	err = service.client.ExecuteRequest(request, response)
	return
}


func (service *AddressService)  Get(documentID string) (response *Address, err error) {
	request, err := service.client.NewRequest("GET", nil, "Verification/Address/%s",documentID)
	if err != nil {
		return
	}
	response = &Address{}
	err = service.client.ExecuteRequest(request, response)
	return
}


func (service *AddressService)  Delete(documentID string) (err error) {
	request, err := service.client.NewRequest("DELETE", nil, "Verification/Address/%s",documentID)
	if err != nil {
		return
	}
	err = service.client.ExecuteRequest(request, nil)
	return
}

func (service *AddressService)  Update(documentID string, params AddressUpdateParams) (response *AddressUpdateResponse, err error) {
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
	request, err := service.client.NewRequest("POST", params, "Verification/Address/%s",documentID)
	if err != nil {
		return
	}
	response = &AddressUpdateResponse{}
	err = service.client.ExecuteRequest(request, response)
	return
}