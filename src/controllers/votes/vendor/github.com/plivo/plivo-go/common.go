package plivo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Meta struct {
	Previous *string
	Next     *string

	TotalCount int64
	Offset     int64
	Limit      int64
}

type BaseListResponse struct {
	ApiID string `json:"api_id" url:"api_id"`
	Meta  Meta   `json:"meta" url:"meta"`
}

type BaseResponse struct {
	ApiId   string `json:"api_id" url:"api_id"`
	Message string `json:"message" url:"message"`
}

func (self Application) ID() string {
	return self.AppID
}

func (self Address) DocumentID() string {
	return self.ID
}

func (self Identity) DocumentID() string {
	return self.ID
}

func (self Account) ID() string {
	return self.AuthID
}

func (self Subaccount) ID() string {
	return self.AuthID
}

func (self Call) ID() string {
	return self.CallUUID
}

func (self LiveCall) ID() string {
	return self.CallUUID
}

func (self Conference) ID() string {
	return self.ConferenceName
}

func (self Endpoint) ID() string {
	return self.EndpointID
}

func (self Message) ID() string {
	return self.MessageUUID
}

func (self Number) ID() string {
	return self.Number
}

func (self PhoneNumber) ID() string {
	return self.Number
}

func (self Pricing) ID() string {
	return self.CountryISO
}

func (self Recording) ID() string {
	return self.RecordingID
}

func IsValidFile(fileName string) error{
	fileParts := strings.Split(fileName, ".")
	if len(fileParts) < 2 {
		err := errors.New("Invalid file specified")
		return err
	}
	extension := strings.ToLower(fileParts[len(fileParts)-1])
	fileExtensionToLowerCase := strings.ToLower(extension)
	if (!(fileExtensionToLowerCase == AddressFileType.JPG || fileExtensionToLowerCase == AddressFileType.PDF || fileExtensionToLowerCase == AddressFileType.PNG)){
		err := errors.New("Only jpg, png and pdf extensions are supported")
		return err
	}

	targetFilePath := filepath.Join(fileName, "./")

	if file, err := os.Stat(targetFilePath); err == nil {
		sizeInMB := file.Size()/ (1024 * 1024);
		if(sizeInMB > 5) {
			err := errors.New("File size exeeds 5 MB limit")
			return err
		}
	} else if os.IsNotExist(err) {
		err := errors.New("File not found in the specified path")
		return err
	}
	return  nil
}

func IsValidCountryParams(countryISO string, fiscalIdentificationCode string, streetCode string , municipalCode string) error{
	if(countryISO == "ES") {
		if(len(fiscalIdentificationCode) == 0) {
			err := errors.New("The parameter fiscal_identification_code is required for Spain numbers")
			return err
		}
	}
	if(countryISO == "DK") {
		if(len(streetCode) == 0 || len(municipalCode) == 0){
			err := errors.New("The parameters street_code and municipal_code are required for Denmark numbers")
			return err
		}
	}
	return nil
}