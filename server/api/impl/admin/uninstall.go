// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package admin

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func (adm *Admin) UninstallApp(appID apps.AppID, sessionToken apps.SessionToken) error {
	app, err := adm.store.App().Get(appID)
	if err != nil {
		return errors.Wrapf(err, "failed to get app. appID: %s", appID)
	}

	if app.OnUninstall != nil {
		creq := &apps.CallRequest{
			Call: *app.OnUninstall,
		}
		resp := adm.proxy.Call(sessionToken, creq)
		if resp.Type == apps.CallResponseTypeError {
			return errors.Wrapf(resp, "call %s failed", creq.Path)
		}
	}

	// delete oauth app
	conf := adm.conf.GetConfig()
	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	client.SetToken(string(sessionToken))

	if app.OAuth2ClientID != "" {
		success, response := client.DeleteOAuthApp(app.OAuth2ClientID)
		if !success {
			if response.Error != nil {
				return errors.Wrapf(response.Error, "failed to delete OAuth2 App - %s", app.AppID)
			}
			return errors.Errorf("failed to delete OAuth2 App - returned with status code %d", response.StatusCode)
		}
	}

	// delete the bot account
	if err := adm.mm.Bot.DeletePermanently(app.BotUserID); err != nil {
		return errors.Wrapf(err, "can't delete bot account for App - %s", app.AppID)
	}

	// delete app from proxy plugin, not removing the data
	if err := adm.store.App().Delete(app.AppID); err != nil {
		return errors.Wrapf(err, "can't delete app - %s", app.AppID)
	}

	adm.mm.Log.Info("Uninstalled the app", "app_id", app.AppID)

	return nil
}
