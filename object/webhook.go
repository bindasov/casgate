// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"

	"github.com/casdoor/casdoor/orm"

	"github.com/xorm-io/core"

	"github.com/casdoor/casdoor/util"
)

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Webhook struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Organization string `xorm:"varchar(100) index" json:"organization"`

	Url            string    `xorm:"varchar(100)" json:"url"`
	Method         string    `xorm:"varchar(100)" json:"method"`
	ContentType    string    `xorm:"varchar(100)" json:"contentType"`
	Headers        []*Header `xorm:"mediumtext" json:"headers"`
	Events         []string  `xorm:"varchar(1000)" json:"events"`
	IsUserExtended bool      `json:"isUserExtended"`
	IsEnabled      bool      `json:"isEnabled"`
}

func GetPaginationWebhooks(owner, organization string, offset, limit int, field, value, sortField, sortOrder string) ([]*Webhook, error) {
	webhooks := []*Webhook{}
	session := orm.GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&webhooks, &Webhook{Organization: organization})
	if err != nil {
		return nil, err
	}

	return webhooks, nil
}

func getWebhooksByOrganization(organization string) ([]*Webhook, error) {
	webhooks := []*Webhook{}
	err := orm.AppOrmer.Engine.Desc("created_time").Find(&webhooks, &Webhook{Organization: organization})
	if err != nil {
		return webhooks, err
	}

	return webhooks, nil
}

func getWebhook(owner string, name string) (*Webhook, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	webhook := Webhook{Owner: owner, Name: name}
	existed, err := orm.AppOrmer.Engine.Get(&webhook)
	if err != nil {
		return &webhook, err
	}

	if existed {
		return &webhook, nil
	} else {
		return nil, nil
	}
}

func GetWebhook(id string) (*Webhook, error) {
	owner, name, err := util.GetOwnerAndNameFromUsernameWithOrg(id)
	if err != nil {
		return nil, err
	}
	return getWebhook(owner, name)
}

func UpdateWebhook(id string, webhook *Webhook) (bool, error) {
	owner, name, err := util.GetOwnerAndNameFromUsernameWithOrg(id)
	if err != nil {
		return false, err
	}
	if w, err := getWebhook(owner, name); err != nil {
		return false, err
	} else if w == nil {
		return false, nil
	}

	affected, err := orm.AppOrmer.Engine.ID(core.PK{owner, name}).AllCols().Update(webhook)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func AddWebhook(webhook *Webhook) (bool, error) {
	affected, err := orm.AppOrmer.Engine.Insert(webhook)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteWebhook(webhook *Webhook) (bool, error) {
	affected, err := orm.AppOrmer.Engine.ID(core.PK{webhook.Owner, webhook.Name}).Delete(&Webhook{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (p *Webhook) GetId() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Name)
}
