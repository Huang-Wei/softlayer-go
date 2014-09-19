package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	services "github.com/maximilien/softlayer-go/services"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

const (
	SOFTLAYER_API_URL  = "api.softlayer.com/rest/v3"
	TEMPLATE_ROOT_PATH = "templates"
)

type softLayerClient struct {
	username string
	apiKey   string

	templatePath string

	httpClient *http.Client

	softLayerServices map[string]softlayer.Service
}

func NewSoftLayerClient(username, apiKey string) *softLayerClient {
	pwd, _ := os.Getwd()
	slc := &softLayerClient{
		username: username,
		apiKey:   apiKey,

		templatePath: filepath.Join(pwd, TEMPLATE_ROOT_PATH),

		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},

		softLayerServices: map[string]softlayer.Service{},
	}

	slc.initSoftLayerServices()

	return slc
}

//softlayer.Client interface methods

func (slc *softLayerClient) GetService(serviceName string) (softlayer.Service, error) {
	slService, ok := slc.softLayerServices[serviceName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("softlayer-go does not support service '%s'", serviceName))
	}

	return slService, nil
}

func (slc *softLayerClient) GetSoftLayer_Account_Service() (softlayer.SoftLayer_Account_Service, error) {
	slService, err := slc.GetService("SoftLayer_Account")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Account_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Virtual_Guest_Service() (softlayer.SoftLayer_Virtual_Guest_Service, error) {
	slService, err := slc.GetService("SoftLayer_Virtual_Guest")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Virtual_Guest_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Virtual_Disk_Image_Service() (softlayer.SoftLayer_Virtual_Disk_Image_Service, error) {
	slService, err := slc.GetService("SoftLayer_Virtual_Disk_Image")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Virtual_Disk_Image_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Security_Ssh_Key_Service() (softlayer.SoftLayer_Security_Ssh_Key_Service, error) {
	slService, err := slc.GetService("SoftLayer_Ssh_Key")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Security_Ssh_Key_Service), nil
}

//Public methods

func (slc *softLayerClient) DoRawHttpRequest(path string, requestType string, requestBody *bytes.Buffer) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s@%s/%s", slc.username, slc.apiKey, SOFTLAYER_API_URL, path)

	var lastResponse http.Response
	switch requestType {
	case "POST", "DELETE":
		req, err := http.NewRequest(requestType, url, requestBody)

		if err != nil {
			return nil, err
		}
		resp, err := slc.httpClient.Do(req)

		if err != nil {
			return nil, err
		} else {
			lastResponse = *resp
		}
	case "GET":
		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		} else {
			lastResponse = *resp
		}
	default:
		return nil, errors.New(fmt.Sprintf("Undefined request type '%s', only GET/POST/DELETE are available!", requestType))
	}

	responseBody, err := ioutil.ReadAll(lastResponse.Body)
	lastResponse.Body.Close()
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (slc *softLayerClient) GenerateRequestBody(templateData interface{}) (*bytes.Buffer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	bodyTemplate := template.Must(template.ParseFiles(filepath.Join(cwd, slc.templatePath)))
	body := new(bytes.Buffer)
	bodyTemplate.Execute(body, templateData)

	return body, nil
}

func (slc *softLayerClient) HasErrors(body map[string]interface{}) error {
	if errString, ok := body["error"]; !ok {
		return nil
	} else {
		return errors.New(errString.(string))
	}
}

func (slc *softLayerClient) CheckForHttpResponseErrors(data []byte) error {
	var decodedResponse map[string]interface{}
	err := json.Unmarshal(data, &decodedResponse)
	if err != nil {
		return err
	}

	if err := slc.HasErrors(decodedResponse); err != nil {
		return err
	}

	return nil
}

//Private methods

func (slc *softLayerClient) initSoftLayerServices() {
	slc.softLayerServices["SoftLayer_Account"] = services.NewSoftLayer_Account_Service(slc)
	slc.softLayerServices["SoftLayer_Virtual_Guest"] = services.NewSoftLayer_Virtual_Guest_Service(slc)
	slc.softLayerServices["SoftLayer_Virtual_Disk_Image"] = services.NewSoftLayer_Virtual_Disk_Image_Service(slc)
	slc.softLayerServices["SoftLayer_Ssh_Key"] = services.NewSoftLayer_Ssh_Key_Service(slc)
}
