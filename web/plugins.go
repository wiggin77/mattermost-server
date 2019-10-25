// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package web

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/mattermost/mattermost-server/mlog"
)

func (w *Web) InitPlugins() {
	mlog.Debug("Initializing plugin routes")

	pluginsRoute := w.MainRouter.PathPrefix("/plugins/{plugin_id:[A-Za-z0-9\\_\\-\\.]+}").Subrouter()
	pluginsRoute.Handle("", w.pluginHandler(servePluginRequest))
	pluginsRoute.Handle("/public/{public_file:.*}", w.pluginHandler(servePluginPublicRequest))
	pluginsRoute.Handle("/{anything:.*}", w.pluginHandler(servePluginRequest))
}

func servePluginRequest(c *Context, w http.ResponseWriter, r *http.Request) {
	c.App.ServePluginRequest(w, r)
}

func servePluginPublicRequest(c *Context, w http.ResponseWriter, r *http.Request) {
	c.App.ServePluginPublicRequest(w, r)
}

// pluginHandler provides a handler which minimizes additional processing and instead largely
// passes through the request unmodified for handling by a plugin.
func (w *Web) pluginHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &Handler{
		GetGlobalAppOptions: w.GetGlobalAppOptions,
		HandleFunc:          h,
		HandlerName:         GetHandlerName(h),
		RequireSession:      false,
		TrustRequester:      false,
		RequireMfa:          false,
		IsStatic:            false,
		Passthrough:         true,
		AllowFormCSRFTokens: true,
	}
	if *w.ConfigService.Config().ServiceSettings.WebserverMode == "gzip" {
		return gziphandler.GzipHandler(handler)
	}

	return handler
}
