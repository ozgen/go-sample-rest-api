package storage

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-sample-rest-api/customerrors"
	"go-sample-rest-api/logging"
	"io/ioutil"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type ImageStore interface {
	UploadImage(ctx context.Context, blobName string, imageData []byte) error
	DownloadImage(ctx context.Context, blobName string) ([]byte, error)
}

type AzureStorage struct {
	AccountName   string
	AccountKey    string
	ContainerName string
	ServiceURL    *azblob.ServiceURL
}

func NewAzureStorage(accountName, accountKey, containerName string) *AzureStorage {
	log := logging.GetLogger()
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil
	}
	p := azblob.NewPipeline(cred, azblob.PipelineOptions{})
	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName))
	if err != nil {
		return nil
	}
	serviceURL := azblob.NewServiceURL(*u, p)

	log.WithFields(logrus.Fields{
		"container": containerName,
	}).Info("Successfully connected to the Azure Blob!")

	return &AzureStorage{
		AccountName:   accountName,
		AccountKey:    accountKey,
		ContainerName: containerName,
		ServiceURL:    &serviceURL,
	}
}

func (az *AzureStorage) UploadImage(ctx context.Context, blobName string, imageData []byte) error {
	containerURL := az.ServiceURL.NewContainerURL(az.ContainerName)
	blobURL := containerURL.NewBlockBlobURL(blobName)

	_, err := azblob.UploadBufferToBlockBlob(ctx, imageData, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize: 4 * 1024 * 1024, // 4 MB
	})
	if err != nil {
		return &customerrors.AzureStorageError{Message: err.Error()}
	}
	return nil
}

func (az *AzureStorage) DownloadImage(ctx context.Context, blobName string) ([]byte, error) {
	containerURL := az.ServiceURL.NewContainerURL(az.ContainerName)
	blobURL := containerURL.NewBlockBlobURL(blobName)

	// Start downloading the blob
	downloadResponse, err := blobURL.Download(
		ctx,
		0,                                 // Start from the beginning
		azblob.CountToEnd,                 // Download the entire blob
		azblob.BlobAccessConditions{},     // No specific access conditions
		false,                             // No need for MD5 of the range
		azblob.ClientProvidedKeyOptions{}, // No customer-provided keys
	)
	if err != nil {
		return nil, &customerrors.AzureStorageError{Message: err.Error()}
	}

	// Read the downloaded data
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 3})
	defer bodyStream.Close()
	data, err := ioutil.ReadAll(bodyStream)
	if err != nil {
		return nil, &customerrors.AzureStorageError{Message: err.Error()}
	}

	return data, nil
}
