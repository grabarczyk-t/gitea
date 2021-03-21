// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package externalaccount

import (
	"strings"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/structs"

	"github.com/markbates/goth"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

const (
	FacebookGetGroupListUrl = "https://graph.facebook.com/v10.0/%s/groups?access_token=%s"
	ApprovedGroupId = "2654430914769218"
	ApprovedGroupRegion = "Małopolska"
)

type UserGroupsResponse struct {
	Groups []UserGroupsData `json:"data"`
}

type UserGroupsData struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ErrUserNotInApprovedGroup struct {
	Email string
}

func IsErrUserNotInApprovedGroup(err error) bool {
	_, ok := err.(ErrUserNotInApprovedGroup)
	return ok
}

// TODO put organization ids to db, add other regions, handle members and volunteers
func GetUserGroupRegion(gothUser goth.User) (string, error) {
	url := fmt.Sprintf(FacebookGetGroupListUrl, gothUser.UserID, gothUser.AccessToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var respObj UserGroupsResponse
	json.Unmarshal(body, &respObj)

	for i, group := range respObj.Groups {
		if group.Id == ApprovedGroupId {
			return ApprovedGroupRegion, nil
		}
	}
	return "", ErrUserNotInApprovedGroup{gothUser.Email}
}

// LinkAccountToUser link the gothUser to the user
func LinkAccountToUser(user *models.User, gothUser goth.User) error {
	loginSource, err := models.GetActiveOAuth2LoginSourceByName(gothUser.Provider)
	if err != nil {
		return err
	}

	externalLoginUser := &models.ExternalLoginUser{
		ExternalID:        gothUser.UserID,
		UserID:            user.ID,
		LoginSourceID:     loginSource.ID,
		RawData:           gothUser.RawData,
		Provider:          gothUser.Provider,
		Email:             gothUser.Email,
		Name:              gothUser.Name,
		FirstName:         gothUser.FirstName,
		LastName:          gothUser.LastName,
		NickName:          gothUser.NickName,
		Description:       gothUser.Description,
		AvatarURL:         gothUser.AvatarURL,
		Location:          gothUser.Location,
		AccessToken:       gothUser.AccessToken,
		AccessTokenSecret: gothUser.AccessTokenSecret,
		RefreshToken:      gothUser.RefreshToken,
		ExpiresAt:         gothUser.ExpiresAt,
	}

	if err := models.LinkExternalToUser(user, externalLoginUser); err != nil {
		return err
	}

	externalID := externalLoginUser.ExternalID

	var tp structs.GitServiceType
	for _, s := range structs.SupportedFullGitService {
		if strings.EqualFold(s.Name(), gothUser.Provider) {
			tp = s
			break
		}
	}

	if tp.Name() != "" {
		return models.UpdateMigrationsByType(tp, externalID, user.ID)
	}

	return nil
}
