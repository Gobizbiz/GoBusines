package controller

import (
	"fmt"

	"github.com/aiteung/athelper"
	"gitlab.com/informatics-research-center/auth-service/model"

	"github.com/aiteung/atmodel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/whatsauth/watoken"
	"github.com/whatsauth/whatsauth"
	"gitlab.com/informatics-research-center/auth-service/config"
	"gitlab.com/informatics-research-center/auth-service/database"
	
)

func WsWhatsAuthQR(c *websocket.Conn) { //simpati unpas lama
	whatsauth.RunSocket(c, config.PublicKey, config.Usertables[:], config.Ulbimariaconn)
}

func PostWhatsAuthRequest(c *fiber.Ctx) error { //receiver whtasapp message token
	if string(c.Request().Host()) == config.Internalhost || string(c.Request().Host()) == "127.0.0.1:7777" {
		var req whatsauth.WhatsauthRequest
		var ntfbtn atmodel.NotifButton
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}
		app := watoken.GetAppSubDomain(req.Uuid)
		getapptried := 0
		for (getapptried < 17) && (app == "") {
			app = watoken.GetAppSubDomain(req.Uuid)
			getapptried += getapptried
		}

		if app == "siapbaak" {
			ntfbtn = whatsauth.RunModuleLegacy(req, config.PrivateKey, config.SiapUserTables[:], config.Ulbimssqlconn)
			fmt.Println(ntfbtn)
		} else if config.CheckIsAkademik(app) {
			ntfbtn = whatsauth.RunWithUsernames(req, config.PrivateKey, config.Usertables[:], config.Ulbimariaconn)
		} else {
			ntfbtn = whatsauth.RunWithUsernames(req, config.PrivateKey, config.AptimasTables[:], config.AptimasConn)
		}
		if app == "" {
			ntfbtn.Message.Message.FooterText = ntfbtn.Message.Message.FooterText + req.Uuid
		}
		return c.JSON(ntfbtn)
	} else {
		var ws whatsauth.WhatsauthStatus
		ws.Status = string(c.Request().Host())
		return c.JSON(ws)
	}
}

func PostWhatsAuthRole(c *fiber.Ctx) error { //receiver whtasapp message token
	if string(c.Request().Host()) != config.Internalhost || string(c.Request().Host()) == "127.0.0.1:7777" {
		var ws whatsauth.WhatsauthStatus
		ws.Status = string(c.Request().Host())
		return c.JSON(ws)
	}
	req := new(whatsauth.WhatsAuthRoles)
	err := c.BodyParser(req)
	if err != nil {
		return err
	}
	var ntfbtn atmodel.NotifButton
	app := watoken.GetAppSubDomain(req.Uuid)
	if app == "siapbaak" {
		ntfbtn := whatsauth.SelectedRoles(*req, config.PrivateKey, config.SiapUserTables[:], config.Ulbimssqlconn)
		fmt.Println(ntfbtn)
	} else if config.CheckIsAkademik(app) {
		ntfbtn = whatsauth.SelectedRoles(*req, config.PrivateKey, config.Usertables[:], config.Ulbimariaconn)
	} else {
		ntfbtn = whatsauth.SelectedRoles(*req, config.PrivateKey, config.AptimasTables[:], config.AptimasConn)
	}
	fmt.Printf("\nreturn button from auth : %+q \n", ntfbtn)
	return c.JSON(ntfbtn)
}