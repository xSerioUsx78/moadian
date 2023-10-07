package service

import (
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"moadian/src/consts"
	"moadian/src/utils/encryption"
	"moadian/src/utils/format"
	"moadian/src/utils/helpers"
	jsonHttp "moadian/src/utils/http"
	baseSchemas "moadian/src/utils/schemas"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)


type Moadian interface {
	GetToken() (*baseSchemas.GetTokenResponseSchema, error)
	GenerateTaxID(datetime time.Time, internalInvoiceID uint) string
	SendInvoice(token *string, packetData *map[string]any) (*baseSchemas.SendInvoiceSchema, error)
	InquiryByReferenceNumber(token *string, referenceNumber *[]string) (*baseSchemas.InquiryResponseSchema, error)
}


type MoadiInfoStruct struct {
	FiscalID string
	PrivateKey []byte
	TaxServerPublicKey string
	TaxServerKeyID string
}


func (info *MoadiInfoStruct) GetToken() (*baseSchemas.GetTokenResponseSchema, error) {
	uuidStr := uuid.New().String()
	packet := map[string]any{
		"uid":        uuidStr,
		"packetType": "GET_TOKEN",
		"retry":      false,
		"data": map[string]any{
			"username": info.FiscalID,
		},
		"encryptionKeyId": "",
		"symmetricKey":    "",
		"iv":              "",
		"fiscalId":        info.FiscalID,
		"dataSignature":   "",
	}
	headers := map[string]string{
		"requestTraceId": uuidStr,
		"timestamp":      strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
	mergedMap := format.MergeMaps(packet, headers)
	normalizedData := helpers.GetNormalizedData(mergedMap)
	signature := b64.StdEncoding.EncodeToString(encryption.Encrypt(info.PrivateKey, &normalizedData))
	data := map[string]any{
		"time":      1,
		"packet":    packet,
		"signature": signature,
	}

	reqURL := consts.TaxAffairsPanelEndpoint + "sync/GET_TOKEN"
	res, err := jsonHttp.SendJsonRequest(
		"POST",
		reqURL,
		data,
		headers,
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	getTokenSchema := new(baseSchemas.GetTokenResponseSchema)
	err = json.Unmarshal(body, &getTokenSchema)
	if err != nil {
		return nil, err
	}
	return getTokenSchema, nil
}

func (info *MoadiInfoStruct) GenerateTaxID(datetime time.Time, internalInvoiceID uint) string {
	return helpers.GenerateTaxID(info.FiscalID, datetime, internalInvoiceID)
}

func (info *MoadiInfoStruct) SendInvoice(
	token *string,
	packetData *map[string]any,
) (
	*baseSchemas.SendInvoiceSchema, 
	error,
) {
	byteData, err := json.Marshal(packetData)
	stringData := string(byteData)
	if (err != nil) {
		return nil, err
	}
	mapData := make(map[string]any)
	decoder := json.NewDecoder(strings.NewReader(stringData))
    decoder.UseNumber()
	if err := decoder.Decode(&mapData); err != nil {
        panic(err)
    }
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
	  return nil, err
	}
	keyHex := hex.EncodeToString(key)
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
	  return nil, err
	}
	ivHex := hex.EncodeToString(iv)
	encryptAesKey, err := encryption.EncryptAesKey(info.TaxServerPublicKey, keyHex)
	if (err != nil) {
		return nil, err
	}
	symmetricKey := b64.StdEncoding.EncodeToString(encryptAesKey)
	encryptionKeyId := info.TaxServerKeyID
	uuidStr := uuid.New().String()
	packetDataNormilized := helpers.GetNormalizedData(mapData)
	dataSignature := b64.StdEncoding.EncodeToString(encryption.Encrypt(info.PrivateKey, &packetDataNormilized))
	packetDataEncrypted, err := encryption.XorAndEncryptData(stringData, keyHex, ivHex)
	if (err != nil) {
		return nil, err
	}
	packet := map[string]any{
		"uid":        uuidStr,
		"packetType": "INVOICE.V01",
		"retry":      false,
		"data": packetDataEncrypted,
		"encryptionKeyId": encryptionKeyId,
		"symmetricKey":    symmetricKey,
		"iv":              ivHex,
		"fiscalId":        info.FiscalID,
		"dataSignature":   dataSignature,
	}
	packetsData := []map[string]any{packet}
	headers := map[string]string{
		"requestTraceId": uuidStr,
		"timestamp":      strconv.FormatInt(time.Now().UnixMilli(), 10),
		"Authorization": *token,
	}
	mergedMap := format.MergeMaps(
		map[string]any{"packets": packetsData}, 
		headers,
	)
	normalizedFinalData := helpers.GetNormalizedData(mergedMap)
	signature := b64.StdEncoding.EncodeToString(encryption.Encrypt(info.PrivateKey, &normalizedFinalData))
	data := map[string]any{
		"time":      1,
		"packets":    packetsData,
		"signature": signature,
	}
	reqURL := consts.TaxAffairsPanelEndpoint + "async/normal-enqueue"
	headers["Authorization"] = "Bearer " + *token
	res, err := jsonHttp.SendJsonRequest(
		"POST",
		reqURL,
		data,
		headers,
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	sendInvoiceSchema := new(baseSchemas.SendInvoiceSchema)
	err = json.Unmarshal(body, &sendInvoiceSchema)
	if err != nil {
		return nil, err
	}
	return sendInvoiceSchema, nil
}

func (info *MoadiInfoStruct) InquiryByReferenceNumber(token *string, referenceNumber *[]string) (*baseSchemas.InquiryResponseSchema, error) {
	uuidStr := uuid.New().String()
	packet := map[string]any{
		"uid":        uuidStr,
		"packetType": "INQUIRY_BY_REFERENCE_NUMBER",
		"retry":      false,
		"data": map[string]any{
			"referenceNumber": *referenceNumber,
		},
		"encryptionKeyId": "",
		"symmetricKey":    "",
		"iv":              "",
		"fiscalId":        info.FiscalID,
		"dataSignature":   "",
	}
	headers := map[string]string{
		"requestTraceId": uuidStr,
		"timestamp":      strconv.FormatInt(time.Now().UnixMilli(), 10),
		"Authorization": *token,
	}
	mergedMap := format.MergeMaps(packet, headers)
	normalizedData := helpers.GetNormalizedData(mergedMap)
	signature := b64.StdEncoding.EncodeToString(encryption.Encrypt(info.PrivateKey, &normalizedData))
	data := map[string]any{
		"time":      1,
		"packet":    packet,
		"signature": signature,
	}

	reqURL := consts.TaxAffairsPanelEndpoint + "sync/INQUIRY_BY_REFERENCE_NUMBER"
	headers["Authorization"] = "Bearer " + *token
	res, err := jsonHttp.SendJsonRequest(
		"POST",
		reqURL,
		data,
		headers,
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	inquirySchema := new(baseSchemas.InquiryResponseSchema)
	err = json.Unmarshal(body, &inquirySchema)
	if err != nil {
		return nil, err
	}
	return inquirySchema, nil
}