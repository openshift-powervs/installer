package powervs

import (
	"strings"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_tenants_ssh_keys"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/pkg/errors"
)

const sshKeyTypeName = "sshKey"

// listSSHKeys lists images in the vpc
func (o *ClusterUninstaller) listSSHKeys() (cloudResources, error) {
	o.Logger.Debugf("Listing images")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listSSHKeys: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	var tenantID string = o.piSession.UserAccount

	params := p_cloud_tenants_ssh_keys.NewPcloudTenantsSshkeysGetallParamsWithTimeout(helpers.PIGetTimeOut).WithTenantID(tenantID)
	resp, err := o.piSession.Power.PCloudTenantsSSHKeys.PcloudTenantsSshkeysGetall(params, ibmpisession.NewAuth(o.piSession, o.ServiceGUID))
	if err != nil {
		return nil, errors.Wrapf(err, "listSSHKeys: PcloudTenantsSshkeysGetall: %v", err)
	}

	var sshKeys *models.SSHKeys = resp.Payload
	var foundOne = false

	result := []cloudResource{}
	for _, sshKey := range sshKeys.SSHKeys {
		if strings.Contains(*sshKey.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listSSHKeys: FOUND: %v", *sshKey.Name)
			result = append(result, cloudResource{
				key:      *sshKey.Name,
				name:     *sshKey.Name,
				status:   "",
				typeName: sshKeyTypeName,
				id:       *sshKey.Name,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listSSHKeys: NO matching sshKey against: %s", o.InfraID)
		for _, sshKey := range sshKeys.SSHKeys {
			o.Logger.Debugf("listSSHKeys: sshKey: %s", *sshKey.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) destroySSHKey(item cloudResource) error {
	var tenantID string = o.piSession.UserAccount
	var getOptions *p_cloud_tenants_ssh_keys.PcloudTenantsSshkeysGetParams
	var err error

	getOptions = p_cloud_tenants_ssh_keys.NewPcloudTenantsSshkeysGetParamsWithTimeout(helpers.PIGetTimeOut).WithTenantID(tenantID).WithSshkeyName(item.id)

	_, err = o.piSession.Power.PCloudTenantsSSHKeys.PcloudTenantsSshkeysGet(getOptions, ibmpisession.NewAuth(o.piSession, o.ServiceGUID))

	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted sshKey %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting sshKey %q", item.name)

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("destroySSHKey: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	params := p_cloud_tenants_ssh_keys.NewPcloudTenantsSshkeysDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithTenantID(tenantID).WithSshkeyName(item.name)
	_, err = o.piSession.Power.PCloudTenantsSSHKeys.PcloudTenantsSshkeysDelete(params, ibmpisession.NewAuth(o.piSession, o.ServiceGUID))
	if err != nil {
		return errors.Wrapf(err, "Failed to delete sshKey %s", item.name)
	}

	return nil
}

// destroySSHKeys removes all image resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySSHKeys() error {
	found, err := o.listSSHKeys()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(sshKeyTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroySSHKeys: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted sshKey %q", item.name)
				continue
			}
			err := o.destroySSHKey(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(sshKeyTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(sshKeyTypeName); len(items) > 0 {
		return errors.Errorf("Error: destroySSHKeys: %d items pending", len(items))
	}
	return nil
}
