/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"path/filepath"

	"github.com/casbin/casbin/v2"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/infra/casbin/adapter"
	"github.com/jacklv111/optimus/infra/casbin/config"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
	"github.com/jacklv111/optimus/pkg/iam/permission/repo"
	pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"
)

type PermissionServiceImpl struct {
}

func (psvc *PermissionServiceImpl) CreatePermission(permission pvb.Permission) error {
	act := filepath.Join(permission.Action...)
	return adapter.Adapter.AddPolicy("p", "p", []string{
		permission.Domain,
		iamconst.ADMIN,
		permission.ResourceType,
		permission.ResourceId,
		act,
		permission.Effect,
	})
}

func (psvc *PermissionServiceImpl) DeletePermission(permission pvb.Permission) error {
	act := filepath.Join(permission.Action...)
	return adapter.Adapter.RemovePolicy("p", "p", []string{
		permission.Domain,
		iamconst.ADMIN,
		permission.ResourceType,
		permission.ResourceId,
		act,
		permission.Effect,
	})
}

func (psvc *PermissionServiceImpl) Enforce(permission pvb.PermissionEnforce) (bool, error) {
	enforcer, err := psvc.getEnforcer(permission)
	if err != nil {
		return false, err
	}
	return enforcer.Enforce(permission.Domain, permission.Name, permission.ResourceType, permission.ResourceId, permission.Action)
}

func (psvc *PermissionServiceImpl) AddRoleForUserInDomain(domain, role, name string) error {
	return adapter.Adapter.AddPolicy("g", "g", []string{domain, role, name})
}

func (psvc *PermissionServiceImpl) DeleteRoleForUserInDomain(domain, role, name string) error {
	return adapter.Adapter.RemoveFilteredPolicy("g", "g", 0, domain, role, name)
}

func (psvc *PermissionServiceImpl) getEnforcer(permission pvb.PermissionEnforce) (*casbin.Enforcer, error) {
	enforcer, err := casbin.NewEnforcer(config.CasbinConfig.ModelPath, adapter.Adapter)
	if err != nil {
		return nil, err
	}
	enforcer.EnableAutoSave(false)

	rules, err := repo.CasbinRepo.GetPolicyList(permission.Domain, permission.Name, permission.ResourceType, permission.ResourceId)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		var policy []string
		if rule.Ptype == "g" {
			policy = []string{rule.V0, rule.V1, rule.V2}
			if _, err := enforcer.AddGroupingPolicy(policy); err != nil {
				log.Errorf("Failed to add policy: %v", err)
				return nil, err
			}
		} else {
			policy = []string{rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5}
			if _, err := enforcer.AddPolicy(policy); err != nil {
				log.Errorf("Failed to add policy: %v", err)
				return nil, err
			}
		}
	}

	return enforcer, nil
}

var PermissionSvc PermissionServiceInterface

func init() {
	PermissionSvc = &PermissionServiceImpl{}
}
