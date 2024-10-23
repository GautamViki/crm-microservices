package controller

import (
	"crm-service/config"
	h "crm-service/helper"
	httpresponse "crm-service/helper/httpResponse"
	"crm-service/internals/dto"
	"crm-service/models"
	"crm-service/services"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"net/http"
)

type customerHandler struct {
}

func NewHandler() *customerHandler {
	return &customerHandler{}
}

// ListCustomers godoc
// @Summary List all customers
// @Description Retrieve a list of customers from the database or cache
// @Tags customers
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.CustomersResponse
// @Failure 400 {object} httpresponse.Response
// @Router /customers [get]
func (ch *customerHandler) ListCustomers(c *gin.Context) {
	dbConfig := config.ConnectDB()
	db := dbConfig.DB
	logger := config.GetLoggerInstance()
	customers, err := services.GetAllCustomersFromCache(db, dbConfig.Redis)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.CustomerFetchError, h.CustomerFetchErrorCode)
		return
	}

	logger.Log(h.CustomerFetchSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerFetchSuccess)
	custResp := dto.CustomersResponse{
		Response:  res,
		Total:     len(customers),
		Customers: customers,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}

// UploadCustomer godoc
// @Summary Upload customers from an Excel file
// @Description Upload and parse an Excel file to create multiple customer records
// @Tags customers
// @Accept  multipart/form-data
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "Excel file"
// @Success 201 {object} httpresponse.Response
// @Failure 400 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /customers/upload [post]
func (ch *customerHandler) UploadCustomer(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	file, err := c.FormFile("file")
	if err != nil {
		logger.Log(h.FileRetrieveFromFormDataError, err.Error(), h.FileRetrieveFromFormDataErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.FileRetrieveFromFormDataError, h.FileRetrieveFromFormDataErrorCode)
		return
	}

	if !strings.HasSuffix(file.Filename, h.XlsxFormat) {
		logger.Log(h.FileFormateInvalidError, "", h.FileFormateInvalidErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.FileFormateInvalidError, h.FileFormateInvalidErrorCode)
		return
	}

	customers, err := services.ParseExcel(file)
	if err != nil {
		logger.Log(h.ExcelFileParseError, err.Error(), h.ExcelFileParseErrorCode)
		h.RespondWithError(c, http.StatusInternalServerError, h.ExcelFileParseError, h.ExcelFileParseErrorCode)
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(customers))

	for _, customer := range customers {
		wg.Add(1)
		go func(cust models.Customer, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := services.AddCustomer(&cust, dbConfig); err != nil {
				errCh <- err
				return
			}
		}(customer, &wg)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		logger.Log(h.CustomerSaveError, err.Error(), h.CustomerSaveErrorCode)
		h.RespondWithError(c, http.StatusInternalServerError, h.CustomerSaveError, h.CustomerSaveErrorCode)
		return
	}
	logger.Log(h.CustomerSaveSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerSaveSuccess)
	h.RespondWithJSON(c, res, http.StatusCreated)
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update customer details by ID
// @Tags customers
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Customer ID"
// @Param customer body models.Customer true "Customer details to update"
// @Success 200 {object} dto.CustomerResponse
// @Failure 400 {object} httpresponse.Response
// @Router /customers/{id} [put]
func (ch *customerHandler) UpdateCustomer(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log(h.CustomerIdInvalidError, err.Error(), h.CustomerIdInvalidErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.CustomerIdInvalidError, h.CustomerIdInvalidErrorCode)
		return
	}

	requset := dto.CustomerRequest{}
	if err := c.ShouldBindJSON(&requset); err != nil {
		logger.Log(h.CustomerDataInvalidError, err.Error(), h.CustomerDataInvalidErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.CustomerDataInvalidError, h.CustomerDataInvalidErrorCode)
		return
	}

	_, err = services.FetchById(id, dbConfig.DB)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		h.RespondWithError(c, http.StatusInternalServerError, h.CustomerFetchError, h.CustomerFetchErrorCode)
		return
	}

	customer, err := services.UpdateCustomer(id, requset, dbConfig)
	if err != nil {
		logger.Log(h.CustomerUpdateError, err.Error(), h.CustomerUpdateErrorCode)
		h.RespondWithError(c, http.StatusInternalServerError, h.CustomerUpdateError, h.CustomerUpdateErrorCode)
		return
	}

	logger.Log(h.CustomerUpdateSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerUpdateSuccess)
	custResp := dto.CustomerResponse{
		Response: res,
		Customer: customer,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer by ID
// @Tags customers
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Customer ID"
// @Success 200 {object} dto.CustomerResponse
// @Failure 400 {object} httpresponse.Response
// @Router /customers/{id} [delete]
func (ch *customerHandler) DeleteCustomer(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log(h.CustomerIdInvalidError, err.Error(), h.CustomerIdInvalidErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.CustomerIdInvalidError, h.CustomerIdInvalidErrorCode)
		return
	}

	customer, err := services.FetchById(id, dbConfig.DB)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		h.RespondWithError(c, http.StatusInternalServerError, h.CustomerFetchError, h.CustomerFetchErrorCode)
		return
	}

	err = services.DeleteCustomer(id, dbConfig)
	if err != nil {
		logger.Log(h.CustomerDeleteError, err.Error(), h.CustomerDeleteErrorCode)
		h.RespondWithError(c, http.StatusInternalServerError, h.CustomerDeleteError, h.CustomerDeleteErrorCode)
		return
	}

	logger.Log(h.CustomerDeleteSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerDeleteSuccess)
	custResp := dto.CustomerResponse{
		Response: res,
		Customer: customer,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}

// GetAllCacheCustomers godoc
// @Summary Get all cached customers
// @Description Retrieve all customers from the cache
// @Tags customers
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.CustomersResponse
// @Failure 400 {object} httpresponse.Response
// @Router /customers/cache [get]
func (ch *customerHandler) GetAllCacheCustomers(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	db := dbConfig.DB
	redisClient := dbConfig.Redis

	customers, err := services.GetAllCustomersFromCache(db, redisClient)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.CustomerFetchError, h.CustomerFetchErrorCode)
		return
	}

	logger.Log(h.CustomerFetchSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerFetchSuccess)
	custResp := dto.CustomersResponse{
		Response:  res,
		Total:     len(customers),
		Customers: customers,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}

// GetUserById godoc
// @Summary Get customers by id
// @Description Get customers by id
// @Tags customers
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Customer ID"
// @Success 200 {object} dto.CustomersResponse
// @Failure 400 {object} httpresponse.Response
// @Router /customers/{id} [get]
func (ch *customerHandler) GetUserById(c *gin.Context) {
	dbConfig := config.ConnectDB()
	logger := config.GetLoggerInstance()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log(h.CustomerIdInvalidError, err.Error(), h.CustomerIdInvalidErrorCode)
		h.RespondWithError(c, http.StatusBadRequest, h.CustomerIdInvalidError, h.CustomerIdInvalidErrorCode)
		return
	}

	customer, err := services.FetchById(id, dbConfig.DB)
	if err != nil {
		logger.Log(h.CustomerFetchError, err.Error(), h.CustomerFetchErrorCode)
		h.RespondWithError(c, http.StatusInternalServerError, h.CustomerFetchError, h.CustomerFetchErrorCode)
		return
	}

	logger.Log(h.CustomerFetchSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.CustomerFetchSuccess)
	custResp := dto.CustomerResponse{
		Response: res,
		Customer: customer,
	}
	h.RespondWithJSON(c, custResp, http.StatusOK)
}
