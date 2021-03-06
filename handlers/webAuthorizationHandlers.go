//Package handlers ...
package handlers

import (
	"fmt"
	m "github.com/Ulbora/GoAuth2/managers"
	"net/http"
	"strconv"
)

/*
 Copyright (C) 2019 Ulbora Labs LLC. (www.ulboralabs.com)
 All rights reserved.

 Copyright (C) 2019 Ken Williamson
 All rights reserved.

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

//PageParams PageParams
type PageParams struct {
	Title      string
	ClientName string
	WebSite    string
	Scope      string
	Error      string
}

//Authorize Authorize
func (h *OauthWebHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	//h.Session.InitSessionStore()
	//fmt.Println("in authorize--------------------------------------")
	s, suc := h.getSession(r)
	// fmt.Println("s", s)
	// fmt.Println("store s", s.Store())
	// fmt.Println("name in getSession s", s.Name())
	// fmt.Println("id getSession s", s.ID)
	// fmt.Println("Options in getSession s", s.Options)
	// fmt.Println("SessionKey in getSession", h.Session.SessionKey)

	if suc {
		loggedInAuth := s.Values["loggedIn"]
		userAuth := s.Values["user"]
		fmt.Println("loggedIn: ", loggedInAuth)
		fmt.Println("user: ", userAuth)

		//session, _ := store.Get(r, "temp-name")

		// fmt.Println("store session", session.Store())
		// fmt.Println("name in getSession session", session.Name())
		// fmt.Println("id getSession session", session.ID)
		// fmt.Println("Options in getSession session", session.Options)

		//loggedIn2 := session.Values["loggedIn"]
		//fmt.Println("loggedIn2", loggedIn2)

		larii := s.Values["authReqInfo"]

		fmt.Println("arii-----auth", larii)

		respTypeAuth := r.URL.Query().Get("response_type")
		fmt.Println("respType: ", respTypeAuth)

		clientIDStrAuth := r.URL.Query().Get("client_id")
		fmt.Println("clientIDStr: ", clientIDStrAuth)

		clientIDAuth, idErr := strconv.ParseInt(clientIDStrAuth, 10, 64)
		fmt.Println("clientIDAuth: ", clientIDAuth)
		fmt.Println("idErr: ", idErr)

		redirectURLAuth := r.URL.Query().Get("redirect_uri")
		fmt.Println("redirURLAuth: ", redirectURLAuth)

		scopeAuth := r.URL.Query().Get("scope")
		fmt.Println("scopeAuth: ", scopeAuth)

		stateAuth := r.URL.Query().Get("state")
		fmt.Println("stateAuth: ", stateAuth)

		var ari AuthorizeRequestInfo
		ari.ResponseType = respTypeAuth
		ari.ClientID = clientIDAuth
		ari.RedirectURI = redirectURLAuth
		ari.Scope = scopeAuth
		ari.State = stateAuth

		if loggedInAuth == true && userAuth != nil {
			fmt.Println("loggedIn: ", loggedInAuth)
			fmt.Println("user: ", userAuth)
			if respTypeAuth == codeRespType {
				var au m.AuthCode
				au.ClientID = clientIDAuth
				au.UserID = userAuth.(string)
				au.Scope = scopeAuth
				au.RedirectURI = redirectURLAuth
				authed := h.Manager.CheckAuthCodeApplicationAuthorization(&au)
				fmt.Println("authed: ", authed)
				if authed {
					ausuc, acode, acodeStr := h.Manager.AuthorizeAuthCode(&au)
					fmt.Println("ausuc: ", ausuc)
					fmt.Println("acode: ", acode)
					fmt.Println("acodeStr: ", acodeStr)
					if ausuc && acode != 0 && acodeStr != "" {
						http.Redirect(w, r, redirectURLAuth+"?code="+acodeStr+"&state="+stateAuth, http.StatusFound)
					} else {
						http.Redirect(w, r, accessDeniedErrorURL, http.StatusFound)
					}
				} else {
					s.Values["authReqInfo"] = ari
					s.Save(r, w)
					larii := s.Values["authReqInfo"]
					fmt.Println("arii-----", larii)
					http.Redirect(w, r, authorizeAppURL, http.StatusFound)
				}
			} else if respTypeAuth == tokenRespType {
				var aut m.Implicit
				aut.ClientID = clientIDAuth
				aut.UserID = userAuth.(string)
				aut.Scope = scopeAuth
				aut.RedirectURI = redirectURLAuth
				iauthed := h.Manager.CheckImplicitApplicationAuthorization(&aut)
				fmt.Println("iauthed: ", iauthed)
				if iauthed {
					isuc, im := h.Manager.AuthorizeImplicit(&aut)
					if isuc && im.Token != "" {
						if h.TokenCompressed {
							im.Token = h.JwtCompress.CompressJwt(im.Token)
						}
						http.Redirect(w, r, redirectURLAuth+"?token="+im.Token+"&token_type=bearer&state="+stateAuth, http.StatusFound)
					} else {
						http.Redirect(w, r, accessDeniedErrorURL, http.StatusFound)
					}
				} else {
					s.Values["authReqInfo"] = ari
					s.Save(r, w)
					http.Redirect(w, r, authorizeAppURL, http.StatusFound)
				}
			} else {
				http.Redirect(w, r, invalidGrantErrorURL, http.StatusFound)
			}
		} else {
			//session, _ := store.Get(r, "temp-name")
			//session.Values["authReqInfo"] = "ari"
			//err := session.Save(r, w)
			//fmt.Println("sesErr", err)
			s.Values["authReqInfo"] = ari
			//s.Values["testval"] = "someTest"
			sesErr := s.Save(r, w)
			fmt.Println("sesErr", sesErr)
			//larii := s.Values["authReqInfo"]
			//fmt.Println("arii-----", larii)
			http.Redirect(w, r, loginURL, http.StatusFound)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//AuthorizeApp AuthorizeApp
func (h *OauthWebHandler) AuthorizeApp(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	if suc {
		loggedInAuth := s.Values["loggedIn"]
		userAuth := s.Values["user"]
		fmt.Println("loggedIn: ", loggedInAuth)
		fmt.Println("user: ", userAuth)

		arii := s.Values["authReqInfo"]
		fmt.Println("arii", arii)
		if arii != nil {
			ari := arii.(*AuthorizeRequestInfo)
			fmt.Println("ari", ari)
			if ari.ResponseType == codeRespType {
				var au m.AuthCode
				au.ClientID = ari.ClientID
				au.RedirectURI = ari.RedirectURI
				fmt.Println("au", au)
				authRes := h.Manager.ValidateAuthCodeClientAndCallback(&au)
				fmt.Println("authRes", authRes)
				if authRes.Valid {
					fmt.Println("authRes", authRes)
					var pg PageParams
					pg.Title = authAppPageTitle
					pg.ClientName = authRes.ClientName
					pg.WebSite = authRes.WebSite
					pg.Scope = ari.Scope
					h.Templates.ExecuteTemplate(w, authorizeHTML, &pg)
				} else {
					var epg PageParams
					epg.Error = invalidRedirectError
					h.Templates.ExecuteTemplate(w, oauthErrorHTML, &epg)
				}
			} else if ari.ResponseType == tokenRespType {
				var auti m.Implicit
				auti.ClientID = ari.ClientID
				auti.RedirectURI = ari.RedirectURI
				iauthr := h.Manager.ValidateImplicitClientAndCallback(&auti)
				if iauthr.Valid {
					var ipg PageParams
					ipg.Title = authAppPageTitle
					ipg.ClientName = iauthr.ClientName
					ipg.WebSite = iauthr.WebSite
					ipg.Scope = ari.Scope
					h.Templates.ExecuteTemplate(w, authorizeHTML, &ipg)
				} else {
					var iepg PageParams
					iepg.Error = invalidRedirectError
					h.Templates.ExecuteTemplate(w, oauthErrorHTML, &iepg)
				}
			} else {
				var ertepg PageParams
				ertepg.Error = invalidRedirectError
				h.Templates.ExecuteTemplate(w, oauthErrorHTML, &ertepg)
			}
		} else {
			var pg PageParams
			pg.Error = invalidReqestError
			h.Templates.ExecuteTemplate(w, oauthErrorHTML, &pg)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//ApplicationAuthorizationByUser ApplicationAuthorizationByUser
func (h *OauthWebHandler) ApplicationAuthorizationByUser(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	if suc {
		aaut := r.URL.Query().Get("authorize")
		fmt.Println("authorize", aaut)

		aaarii := s.Values["authReqInfo"]
		user := s.Values["user"]
		fmt.Println("aaarii", aaarii)
		if aaarii != nil && user != nil {

			aaari := aaarii.(*AuthorizeRequestInfo)
			fmt.Println("ari", aaari)
			if aaut == "true" && aaari.ResponseType == codeRespType {
				var aac m.AuthCode
				aac.ClientID = aaari.ClientID
				aac.UserID = user.(string)
				aac.Scope = aaari.Scope
				aac.RedirectURI = aaari.RedirectURI
				authSuc, authCode, authCodeStr := h.Manager.AuthorizeAuthCode(&aac)
				if authSuc && authCode != 0 && authCodeStr != "" {
					fmt.Println("authSuc", authSuc)
					fmt.Println("authCode", authCode)
					fmt.Println("authCodeStr", authCodeStr)
					http.Redirect(w, r, aaari.RedirectURI+"?code="+authCodeStr+"&state="+aaari.State, http.StatusFound)
				} else {
					var aaacpg PageParams
					aaacpg.Error = accessDenidError
					h.Templates.ExecuteTemplate(w, oauthErrorHTML, &aaacpg)
				}
			} else if aaut == "true" && aaari.ResponseType == tokenRespType {
				var aai m.Implicit
				aai.ClientID = aaari.ClientID
				aai.UserID = user.(string)
				aai.Scope = aaari.Scope
				aai.RedirectURI = aaari.RedirectURI
				aaiSuc, irtn := h.Manager.AuthorizeImplicit(&aai)
				fmt.Println("aaiSuc", aaiSuc)
				fmt.Println("irtn", irtn)
				if aaiSuc && irtn.Token != "" {
					if h.TokenCompressed {
						irtn.Token = h.JwtCompress.CompressJwt(irtn.Token)
					}
					http.Redirect(w, r, aaari.RedirectURI+"?token="+irtn.Token+"&token_type=bearer&state="+aaari.State, http.StatusFound)
				} else {
					var aaipg PageParams
					aaipg.Error = accessDenidError
					h.Templates.ExecuteTemplate(w, oauthErrorHTML, &aaipg)
				}
			} else {
				http.Redirect(w, r, aaari.RedirectURI+"?error=access_denied&state="+aaari.State, http.StatusFound)
			}
		} else {
			var pg PageParams
			pg.Error = invalidReqestError
			h.Templates.ExecuteTemplate(w, oauthErrorHTML, &pg)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//OauthError OauthError
func (h *OauthWebHandler) OauthError(w http.ResponseWriter, r *http.Request) {
	authError := r.URL.Query().Get("error")
	var epg PageParams
	epg.Error = authError
	h.Templates.ExecuteTemplate(w, oauthErrorHTML, &epg)
}
