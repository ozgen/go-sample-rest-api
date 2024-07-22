package camerametadata

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const base64Data = "iVBORw0KGgoAAAANSUhEUgAAAOoAAACUCAMAAACqYkXNAAAApVBMVEX////BEhwAAADABxXFNTjlsrOqqqq9AAD5 fn8/Pzx8fH19fXPz8/u7u7h4eE9PT2GhoZwcHDn5 fY2Nh3d3fBwcGPj49/f3/Hx8eXl5e1tbVjY2OioqK7u7teXl4JCQlRUVG AAlGRkYWFhb57e7ourscHBzNYmPy2NnAJyjGQkQuLi7cl5j14 QlJSXHSUvTdnjrxsfVgIHhpabKU1XajI7RbW8Dt1kvAAAP80lEQVR4nN1d6WKqSgzGYhFkEURBQCkuPV3scmx7 v6PdpMZKFQTXCpWb35RoUCYkC/J5BsUpTnRbX8SpWFfUe4/pq166dzdK4oTuhPf1hu8pYZE7zqj8SAdKsrib2eLpqDrE6ioeeOh3zV  873Fb3bt2J3NobN95a6VVW1dQtHBqFm9c1L0xXNN0gHPpjv3fZBhWH9fFYUKwwuz4Sl WYxbP7bRVPQ9R1GM/LGowsbVkNH880CU1Fu29vNF2X6sIIn5Iaxf1m6Cu beJaiPH/uNqgwrP8WijIaRJflmYyuY2lhJnwSqSk10GpHeKaBdlEmDIMaB2loK8rqdUoZK mpOq8ArvbAjX3buBRdwScNx56A1DdKp mj8kA9gc4TmHDsXZAJGwCpmptFsHmtEqaqdu6V w65AzyTEl4QuOY yeEgtfMPDiKHu3MH4Op7EYDrRaiaR4SaophPlD7qAx61eKBcU cd4ofxAMC1ewGBhGEipKYu3Ootrc5KHHdLvsUP4JkUz9UuAlz1rj MkgQiwmcyzO985sfRxv0XPJN1GSaMgyohVSchVX1Y5Efe02N a TgevZeWO9jmO91AVIfp5Qq18WRBvkopo/gmbpeAPHhmYMrmi/4pBEHqZ2P8tgF45/hLR8m46Fz3p4pN98Atq5J8209V46 pwJEAFfjEsBVQmrSB5/0QQ7ZU/VoGow6r/A295Mz90x6t49Z6oTTYvq6 Hb8goyQReaqQebqmOdrwmbfB0gdgPmuyOQlh9QvMW6nlAmrElzjMzZhMF A1MxnPc7b igxvusTUno/iyZn65kgzB9p4UzjIFXgyJrc84gUeeCZzlRXgNQ4SBK4ufs2qcDtpjnqpKOePsBD0VMEV/0cTVhEhF4PzNekI8I3k/gvuiDT Qu7rOxMwVVC6oyFVBnJb8qKtgB0YBAfnqVnkvMWGUSEC/rm3 nxYWDpEWDJzqL4DMvCElKXEBEaZOG3c7dg/vOeDjYwfx mZwiuwnzdLITNFVlLaa04QzSu6RQH7N0IB cHrgJSvaUNBvlKR4SUT5KyIM1g gi7HABX/7w8k4TUJc5b0HHtxyaklkJPSnbeYVecamdWfNFFRJjA1j1dIyQgtRQ6cxVZkIll4T5vEScXLPxG3osDm0xEyPkkKUxpBnNb/7zAVULqEiJChYbURxpSS2EqFrdw7jHGh2czrFhkiZIZPPoFndC8b/OhOl1GbYvM9YzmXMW8RfjHUjjz/azzSVKemfqhcl7xIZpvMAsMprarPtzucBISXAGN4fQyPmxcjR1ERIReDyBVoc33aZcRoTNXVXimJDqLYTUQUsfhHCJCuhWg87HNJ0lZkeA6fYJLTABcz0BXQ2Sp2QAMjJ5dm15vP4mQJxKR2/Cg qEovjSqxw4iCr/esg bZNODmPbfSej2JvVOORfPJML8ZWwwrQDTx9X2c RyS4Gr2oL40NSS389c0XyjJLV5SN3jZOS7rmJB1Rn8OrjKudS5iAiJIQFg3A6ppdDgqubg sueCedSg1mE7Ulk08rDrj5JyjsJrm3AZVvEh7 oKgwqQGqGPolOr//t XqRNbXpBzxKS8xs/JquosgSviCkvpHme7cbpJZyT81sqOqbkhdffs2ERZaauTYzLSyc557yRDk3kRk52JP3W8Oq25il9tAn0TnY2wEnpcN BNdR lvxoYEtk1hkgYv/owKd6evukFrKigpE1AewD12A66/oit11UerhXCp5e rT9nMQ8pc0EKxN b/VLSyy1EFPZKnk3X3uGhF F2bOVXimX2obQPN1e2OTQcNpe5cslRIycxWnsyPMXE8 YQU aQIRoc8Nw/QQnySFNxIr/AVwFRFh2Bsq HJRo1Bb K2Xe/7V17zTZ67YsB5kgc00rAuXebCQmatw6P0AM9eTqiohNbW4VgARoh8uH2SUifOzw8GJwdUwsWG9p4GDeCKdyNbCb73QwZcA1/HgpOBq5A3rDhe0tv798Ap8SO277injQ9mwnuFcKlMkOQxSK1d4pBMlHRveTxn267JhHV6da9onHQqppZANTQJcu1F4ssxVmG UeD5bNjgcUkthukFOyyYrOWDGP7L0QLQn7S9kw7s6FT15AzThE8RMggMWpAip5PzZIVkqJfScngBX90TgKiE1QUgljUz9PNKF6G4QnKlFcD2BCUtIzcaGYpC8VFGOP4o80 0j10iIDE/Q8C7nLWpo1dPDslRK6K4KSdU gWfaygH7OOKzpss4ODMSN  ZjBxSdW5CdHpIkYWTFfkw2ysDea4Ng6uRc8BYWvVBlTNe6G4QwSZrembDkBywuLYz8IhS063YMLjmkBryJJpjmi8K3zDTDV3NamxYDdmwjpBazwE7nuh0fPiWs8kaC/vRfLUwFRwwml5wjIjwu9BUHPU2B9eGTDjPUj2byVL3m0vd9aL0BPUrPNQ ssmaMWHDFpBawwFr4rILhpoEGeQkiZrJXOW8RRpARHhLE7yOFRF lxX5rmB7sYGeqYG2AQGpgeCA1bRlNyB0T55oGnfSRsrCxleWuqCvfWRILeX5kbzeu16wyY48rDCoSKsewPlXdbTqBsRg2vngfdG942euElK9mcM2Pr41F3vTDU2dT12wyYZ 96izOLJhPadVUwWu1vEhtRSGjCUyV0 zjmrCElKTBDz8M8MQP97FiMszHCXwDmZ2XHA1JKTOrFpqYoNCl3bESzOSDU3HGtecVu0iL5VEuYYgtZQVTX4Q66Ack02WZ6kzk8 qmq5UMlTtNlzX7AlwPc4dGCZCam/CwvlrU5BaChO2YCFrkmHYf5Sr5FmqhxekzbcxSK0IDa4ddPzpscDVkLxU5IAt9uaAHU8Yf4hssv6xTLiYS2UTqnaTkFpKzeIh2NDk/Px552F DzYXVPBA0qqbEIZN1oFdehbFPy  gKbgk7wbSGgMfmbhJMIUebAnz mJhqafPXMJqb1A4RxDq2lILYVMk0WeYRyBiiMh1bvCTdp8GVp1E8IsyoUm3J393DOJ9qQrX6krP59MasrssuH9B8sWySJLJnip9BM9BaSWck3huqBqI8/1J CaZ6lovsp0Sl3k6VQ SQpD1W7DLnsZDZ0fmDBEhJr7gt11pKfvvN77GiuTrmJYY2JHXF7AjgfLq0KywdgxlAl7wvGIA1eID/U4wRU8D1Q0z1Iz2NTJlUMBvwdXvLjKhN7h5xewZut7QlvJ BMuJwpZGFCnYFzd5HATFuY7HlxhJN3mIDWsUTVQRi/kjr40mWRzj6brKX/C2ZABV8Emy6nah5hwDqm4HCyz5iD4pBpVbyLF6pF7xKO3qT2WYteMamZh9yZhwqp6C4/uYKp2HhFm JRo88UsdYuqM3IPnp4ePUtxalTF7mOGTfYAj6 fBgeCq4DUFzw7Xe4QJJrw6uZmPv8DclNoiH/M5/ObP5Hi5/rczPEX E2IUHVU3P fl2Uvmy1fXv7Mr1IHVYUTijMWR8zzM17hOpAMuGKrqjFCKs7 sXDeCoDraZLz9Hl3nRUEY20ysvx oZXX931rGGuBZn6pCsnzZBJr4yhw3XAww9N7hVGO8gvqvpYgu0Nzo3E8tCzfyY9YanDC0TAeB0NxIM0mw7KwHYQHrIOSz6X20CfRNcJ1WnV 7zdu5Tcn9zyZs35 PddjNlrfU0p SM///rNOtpVOX8X1DpjZwCIL KSJwdKq1 OkQtWwqmr 42ztbktVe8Ptqlprv5MrSasPJbjuNaxi3iLIcD3NBT0TtkFhHNSouj4wcFPzwtXwOTWnKtPw/vhF1d7LMwlIlRwwsrl8upmlhpQB86p Odp5wJncDafq4pM0YWxPsJI9e/Lk6kkzDf7lndSUoFXXqjqy 0IcpxjEuAoisUMtRHnDPSeOGgGW1pXguvOwSkhNxXqaNIWR6K5za1Sdz7Ism/WWL3D3Y7lLX1Z0Bb892YzpeFWVN L5g/9YfDU07aqrnEvFeQvlkxxUaqWSOlWrkub7rLWgsac5a/dXo oz1d otnBmY5jsAa5yieos6jJLVKttila9r6pGPFvbkw2/e6gaVRlyGr5X/WiPRWjlvAVELsqCzJla/6hntq quGDJcm1f9K1IX6eq/pcEe2wbsHbvyZOQOhOQShWumJVKalXFoDDf9CpX0pK1fGBcHdc6VZV7irORg uubQPCfIPM7dfRZgmp88AvYRQFEBS6YTj4Fh6Z8SCrvrM3w8od1qrKULVxGPoiPtzBM8kii6Awkr3N6l 6yFKoulsIUYoP F3q6vU3Vd3AVXmfd6THxLsbeaL4sk3TfEFCjeWlciuVHKoqDO3EnX3pWtHrqlZVxmVKcN1lzjWH1EGf44A9cA3rdTHwkrnbQvThV0WiYsFbVGXAVXQLu8HE2gaueXvSiPNxU5YDVqg6qKqa7KYq1nFzxeINVZdc8rN4JVstER8mg61U7Qqk0hwwnladEKr6 Y8vW1U1g0LV8v4KVdnk55YO wVVeyu4lg3r9ArrLb49Ka1RdT7SddPsdru2DcGwg2 uPvqW0liFBVeGcKuqzGoDlW T8brKVoAMjMigeamv/LxFYYHfVM0K/RFsBNwA3gj3Pup5gRaPLMu3rFE8TvKk7qV0YEah6oS9JhOhv8P9i4Z3flhzSOU5YHW06sKJVkIErox2peUFl5tZlqaQCMy/9gRlvGQWqsabFyuEntwW4OrWrvBeNqwvSNNQ72pe9EKrpKoqXRy9ihVzTu7IKqhUVE9fNP6iOr0w6xebjPNMxUdL9LpSFSuFVmnlt9F6kJvLRHHI39NRZRT6harjmqvSs2YIrmbEf5vMMO2SA0b6pFoOWKHVrPLb8IbU6GqojIlfX1y/emOFqvOg7rLMEtPyw58x45lkK0A6wW4S0nw/aqfdwp6UagzsJz1KAHui2Zqe8zQa2t/vJz86q6m0cU1j06c6NlkeEQYsrbpVP5fqWFIqQayi xYlI1PpW8PJOHAHCYgXBuN45G90WRWH189N0FXqNtxs16XBdSsHTD16w7rRtUW1qW//iGvANHg y8zVsTd0lZCK3wEzyNojvKnX5yl0k496nX bbMOEiyy1y63532q1z1TIdRVa07Zkk22Cq7GFAyYe1JkKfbsCXEXm n1Yc1p1xNGqL1DyTxKvU3EAUq1JlA64iPAiRSwyZnsIrhVdjZJWvePnfi9Big9/Vj2ThNTUNbjmtssU0RRoYENTOay5T IjwgsVySb75plEmJ yH1K6XBEdkBHSl4uIC6OHBHmp9HfALlcE08lMo68k3TCxSOhyTeOXLKLZHpfJyy0YXtXYRfulP 91ySK67SPAmy9VrThEpPm/2S8IUjY09MEVVZMJ90XUi5YNVdGAo/8b1KAIuIHg8MuATfGRAIVZceKSRdBFsqB0SwJsUpbHeLkiVvER2U0RLhmi1IIFXPOp0/ntnOxo0pEcikR04yn/AVtHp0x6vI5IAAAAAElFTkSuQmCC"

func TestHandler_CreateCameraMetadata(t *testing.T) {
	t.Run("CreateCameraMetadata_withValidData_returnCreated", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		payload := types.CameraMetadataPayload{
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
		}
		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		expectedCamera := types.CameraMetadata{
			CamID:           uuid.New().String(),
			CameraName:      payload.CameraName,
			FirmwareVersion: payload.FirmwareVersion,
			CreatedAt:       nullTime,
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("CreateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)

		cameraData, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBuffer(cameraData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}

		// Use captured argument to assert values
		if capturedArg.CameraName != payload.CameraName {
			t.Errorf("expected CameraName %s, got %s", payload.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != payload.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", payload.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("CreateCameraMetadata_withMalformedData_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)
		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBufferString("{invalid json"))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("CreateCameraMetadata_withInvalidJsonData_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		payload := types.CameraMetadataPayload{
			CameraName:      "",
			FirmwareVersion: "",
		}

		cameraData, _ := json.Marshal(payload)

		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBuffer(cameraData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("CreateCameraMetadata_withDBError_returnInternalError", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		payload := types.CameraMetadataPayload{
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
		}

		mockCameraStore.On("CreateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Return(nil, fmt.Errorf("DB error"))

		cameraData, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBuffer(cameraData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})
}

func TestHandler_InitializeCameraMetaData(t *testing.T) {
	t.Run("InitializeCameraMetaData_withValidData_returnOk", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)
		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if capturedArg.CameraName != expectedCamera.CameraName {
			t.Errorf("expected CameraName %s, got %s", expectedCamera.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != expectedCamera.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withAlreadyInitializedData_returnConflict", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)

		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusConflict {
			t.Errorf("expected status code %d, got %d", http.StatusConflict, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withErrOnUpdate_returnInternalError", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Return(nil, fmt.Errorf("Failed"))
		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withWrongCamID_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("Not Found"))
		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withInvalidCamID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := "123"

		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("InitializeCameraMetaData_withEmptyCamID_returnStatusMovedPermanently", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata//init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusMovedPermanently {
			t.Errorf("expected status code %d, got %d", http.StatusMovedPermanently, rr.Code)
		}
	})
}

func TestHandler_GetCameraMetaData(t *testing.T) {
	t.Run("GetCameraMetaData_withValidData_returnOk", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}", handler.GetCameraMetaData).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("GetCameraMetaData_withCamId_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("Not Found"))

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}", handler.GetCameraMetaData).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("GetCameraMetaData_withMalformedCamId_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := "123"

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}", handler.GetCameraMetaData).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

func TestHandler_UploadImageHandler(t *testing.T) {

	t.Run("UploadImageHandler_withValidData_returnOk", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)
		mockAzureStorage.On("UploadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png", mock.AnythingOfType("[]uint8")).Return(nil)
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if capturedArg.CameraName != expectedCamera.CameraName {
			t.Errorf("expected CameraName %s, got %s", expectedCamera.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != expectedCamera.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.NameOfStoredPicture.String != imageID {
			t.Errorf("expected NameOfStoredPicture %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.ImageId.String != imageID {
			t.Errorf("expected ImageId %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
		mockAzureStorage.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withUploadError_returnInternalServerError", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)
		mockAzureStorage.On("UploadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png", mock.AnythingOfType("[]uint8")).Return(fmt.Errorf("upload err"))
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		if capturedArg.CameraName != expectedCamera.CameraName {
			t.Errorf("expected CameraName %s, got %s", expectedCamera.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != expectedCamera.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.NameOfStoredPicture.String != imageID {
			t.Errorf("expected NameOfStoredPicture %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.ImageId.String != imageID {
			t.Errorf("expected ImageId %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
		mockAzureStorage.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withUpdateCameraMetadataError_returnInternalServerError", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(nil, fmt.Errorf("update error"))
		mockAzureStorage.On("UploadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png", mock.AnythingOfType("[]uint8")).Return(fmt.Errorf("upload err"))
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		if capturedArg.CameraName != expectedCamera.CameraName {
			t.Errorf("expected CameraName %s, got %s", expectedCamera.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != expectedCamera.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.NameOfStoredPicture.String != imageID {
			t.Errorf("expected NameOfStoredPicture %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.ImageId.String != imageID {
			t.Errorf("expected ImageId %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withNotInitCamera_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withWrongCameraID_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("not found"))
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withMalformedBase64_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := uuid.New().String()

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + "123123"

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withMalformedImageID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := "12"

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withNullImageID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := ""

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withInvalidCamID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		imageID := uuid.New().String()
		camID := "123"

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})
}

func TestHandler_DownloadImageHandler(t *testing.T) {

	t.Run("DownloadImageHandler_withValidData_returnOk", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
			ImageId:         sql.NullString{String: imageID, Valid: true},
		}

		imageData, err := base64.StdEncoding.DecodeString(utils.NormalizeBase64(base64Data))
		if err != nil {
			t.Fatal(err)
		}
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockAzureStorage.On("DownloadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png").Return(imageData, nil)

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
		mockAzureStorage.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withWriteFailure_returnsInternalServerError", func(t *testing.T) {
		// Arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := uuid.New().String()
		imageData := []byte("image data")

		expectedCamera := types.CameraMetadata{
			CamID:   camID,
			ImageId: sql.NullString{String: imageID, Valid: true},
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockAzureStorage.On("DownloadImage", mock.Anything, imageID+".png").Return(imageData, nil)

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()                      // Normal recorder
		fw := &failWriter{ResponseWriter: rr, fail: true} // Wrap the recorder to simulate failure

		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(fw, req) // Use the failing writer here

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("DownloadImageHandler_withDownLoadError_returnsInternalServerError", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
			ImageId:         sql.NullString{String: imageID, Valid: true},
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockAzureStorage.On("DownloadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png").Return(nil, fmt.Errorf("download err"))

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withNullImageID_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withWrongCameraID_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("not found"))

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withMalformedCameraID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := "123"

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})
}

type mockCameraStore struct {
	mock.Mock
}

func (m *mockCameraStore) CreateCameraMetadata(c types.CameraMetadata) (*types.CameraMetadata, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.CameraMetadata), args.Error(1)
}

func (m *mockCameraStore) UpdateCameraMetadata(c types.CameraMetadata) (*types.CameraMetadata, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.CameraMetadata), args.Error(1)
}
func (m *mockCameraStore) GetCameraMetadataByID(c string) (*types.CameraMetadata, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.CameraMetadata), args.Error(1)
}

type mockAzureStorage struct {
	mock.Mock
}

func (m *mockAzureStorage) UploadImage(ctx context.Context, blobName string, imageData []byte) error {
	args := m.Called(ctx, blobName, imageData)
	return args.Error(0)
}

func (m *mockAzureStorage) DownloadImage(ctx context.Context, blobName string) ([]byte, error) {
	args := m.Called(ctx, blobName)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

type failWriter struct {
	http.ResponseWriter
	fail bool
}

func (fw *failWriter) Write(data []byte) (int, error) {
	if fw.fail {
		return 0, fmt.Errorf("simulated write error")
	}
	return fw.ResponseWriter.Write(data)
}
