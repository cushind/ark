/*
Copyright 2017 the Heptio Ark contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azure

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/pkg/errors"
	"github.com/satori/uuid"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/heptio/ark/pkg/cloudprovider"
	"github.com/heptio/ark/pkg/util/collections"
)

const (
	azureClientIDKey         = "AZURE_CLIENT_ID"
	azureClientSecretKey     = "AZURE_CLIENT_SECRET"
	azureSubscriptionIDKey   = "AZURE_SUBSCRIPTION_ID"
	azureTenantIDKey         = "AZURE_TENANT_ID"
	azureStorageAccountIDKey = "AZURE_STORAGE_ACCOUNT_ID"
	azureStorageKeyKey       = "AZURE_STORAGE_KEY"
	azureResourceGroupKey    = "AZURE_RESOURCE_GROUP"
	apiTimeoutKey            = "apiTimeout"
	snapshotsResource        = "snapshots"
	disksResource            = "disks"
)

type blockStore struct {
	disks         *disk.DisksClient
	snaps         *disk.SnapshotsClient
	subscription  string
	resourceGroup string
	apiTimeout    time.Duration
}

type snapshotIdentifier struct {
	subscription  string
	resourceGroup string
	name          string
}

func getConfig() map[string]string {
	cfg := map[string]string{
		azureClientIDKey:         "",
		azureClientSecretKey:     "",
		azureSubscriptionIDKey:   "",
		azureTenantIDKey:         "",
		azureStorageAccountIDKey: "",
		azureStorageKeyKey:       "",
		azureResourceGroupKey:    "",
	}

	for key := range cfg {
		cfg[key] = os.Getenv(key)
	}

	return cfg
}

func NewBlockStore() cloudprovider.BlockStore {
	return &blockStore{}
}

func (b *blockStore) Init(config map[string]string) error {
	var (
		apiTimeoutVal = config[apiTimeoutKey]
		apiTimeout    time.Duration
		err           error
	)

	if apiTimeout, err = time.ParseDuration(apiTimeoutVal); err != nil {
		return errors.Wrapf(err, "could not parse %s (expected time.Duration)", apiTimeoutKey)
	}

	if apiTimeout == 0 {
		apiTimeout = 2 * time.Minute
	}

	cfg := getConfig()

	spt, err := helpers.NewServicePrincipalTokenFromCredentials(cfg, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return errors.Wrap(err, "error creating new service principal token")
	}

	disksClient := disk.NewDisksClient(cfg[azureSubscriptionIDKey])
	snapsClient := disk.NewSnapshotsClient(cfg[azureSubscriptionIDKey])

	authorizer := autorest.NewBearerAuthorizer(spt)
	disksClient.Authorizer = authorizer
	snapsClient.Authorizer = authorizer

	b.disks = &disksClient
	b.snaps = &snapsClient
	b.subscription = cfg[azureSubscriptionIDKey]
	b.resourceGroup = cfg[azureResourceGroupKey]
	b.apiTimeout = apiTimeout

	return nil
}

func (b *blockStore) CreateVolumeFromSnapshot(snapshotID, volumeType, volumeAZ string, iops *int64) (string, error) {
	snapshotIdentifier, err := parseFullSnapshotName(snapshotID)
	if err != nil {
		return "", err
	}

	// Lookup snapshot info for its Location
	snapshotInfo, err := b.snaps.Get(snapshotIdentifier.resourceGroup, snapshotIdentifier.name)
	if err != nil {
		return "", errors.WithStack(err)
	}

	diskName := "restore-" + uuid.NewV4().String()

	disk := disk.Model{
		Name:     &diskName,
		Location: snapshotInfo.Location,
		Properties: &disk.Properties{
			CreationData: &disk.CreationData{
				CreateOption:     disk.Copy,
				SourceResourceID: &snapshotID,
			},
			AccountType: disk.StorageAccountTypes(volumeType),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.apiTimeout)
	defer cancel()

	_, errChan := b.disks.CreateOrUpdate(b.resourceGroup, *disk.Name, disk, ctx.Done())

	err = <-errChan

	if err != nil {
		return "", errors.WithStack(err)
	}
	return diskName, nil
}

func (b *blockStore) GetVolumeInfo(volumeID, volumeAZ string) (string, *int64, error) {
	res, err := b.disks.Get(b.resourceGroup, volumeID)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return string(res.AccountType), nil, nil
}

func (b *blockStore) IsVolumeReady(volumeID, volumeAZ string) (ready bool, err error) {
	res, err := b.disks.Get(b.resourceGroup, volumeID)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if res.ProvisioningState == nil {
		return false, errors.New("nil ProvisioningState returned from Get call")
	}

	return *res.ProvisioningState == "Succeeded", nil
}

func (b *blockStore) CreateSnapshot(volumeID, volumeAZ string, tags map[string]string) (string, error) {
	// Lookup disk info for its Location
	diskInfo, err := b.disks.Get(b.resourceGroup, volumeID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	fullDiskName := getComputeResourceName(b.subscription, b.resourceGroup, disksResource, volumeID)
	// snapshot names must be <= 80 characters long
	var snapshotName string
	suffix := "-" + uuid.NewV4().String()

	if len(volumeID) <= (80 - len(suffix)) {
		snapshotName = volumeID + suffix
	} else {
		snapshotName = volumeID[0:80-len(suffix)] + suffix
	}

	snap := disk.Snapshot{
		Name: &snapshotName,
		Properties: &disk.Properties{
			CreationData: &disk.CreationData{
				CreateOption:     disk.Copy,
				SourceResourceID: &fullDiskName,
			},
		},
		Tags:     &map[string]*string{},
		Location: diskInfo.Location,
	}

	for k, v := range tags {
		val := v
		(*snap.Tags)[k] = &val
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.apiTimeout)
	defer cancel()

	_, errChan := b.snaps.CreateOrUpdate(b.resourceGroup, *snap.Name, snap, ctx.Done())
	err = <-errChan

	if err != nil {
		return "", errors.WithStack(err)
	}

	return getComputeResourceName(b.subscription, b.resourceGroup, snapshotsResource, snapshotName), nil
}

func (b *blockStore) DeleteSnapshot(snapshotID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.apiTimeout)
	defer cancel()

	_, errChan := b.snaps.Delete(b.resourceGroup, snapshotID, ctx.Done())

	err := <-errChan

	return errors.WithStack(err)
}

func getComputeResourceName(subscription, resourceGroup, resource, name string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/%s/%s", subscription, resourceGroup, resource, name)
}

var snapshotURIRegexp = regexp.MustCompile(
	`^\/subscriptions\/(?P<subscription>.*)\/resourceGroups\/(?P<resourceGroup>.*)\/providers\/Microsoft.Compute\/snapshots\/(?P<snapshotName>.*)$`)

// parseFullSnapshotName takes a snapshot URI and returns a snapshot identifier
// or an error if the URI does not match the regexp.
func parseFullSnapshotName(name string) (*snapshotIdentifier, error) {
	submatches := snapshotURIRegexp.FindStringSubmatch(name)
	if len(submatches) != len(snapshotURIRegexp.SubexpNames()) {
		return nil, errors.New("snapshot URI could not be parsed")
	}

	snapshotID := &snapshotIdentifier{}

	// capture names start at index 1 to line up with the corresponding indexes
	// of submatches (see godoc on SubexpNames())
	for i, names := 1, snapshotURIRegexp.SubexpNames(); i < len(names); i++ {
		switch names[i] {
		case "subscription":
			snapshotID.subscription = submatches[i]
		case "resourceGroup":
			snapshotID.resourceGroup = submatches[i]
		case "snapshotName":
			snapshotID.name = submatches[i]
		default:
			return nil, errors.New("unexpected named capture from snapshot URI regex")
		}
	}

	return snapshotID, nil
}

func (b *blockStore) GetVolumeID(pv runtime.Unstructured) (string, error) {
	if !collections.Exists(pv.UnstructuredContent(), "spec.azureDisk") {
		return "", nil
	}

	volumeID, err := collections.GetString(pv.UnstructuredContent(), "spec.azureDisk.diskName")
	if err != nil {
		return "", err
	}

	return volumeID, nil
}

func (b *blockStore) SetVolumeID(pv runtime.Unstructured, volumeID string) (runtime.Unstructured, error) {
	azure, err := collections.GetMap(pv.UnstructuredContent(), "spec.azureDisk")
	if err != nil {
		return nil, err
	}

	azure["diskName"] = volumeID
	azure["diskURI"] = getComputeResourceName(b.subscription, b.resourceGroup, disksResource, volumeID)

	return pv, nil
}
