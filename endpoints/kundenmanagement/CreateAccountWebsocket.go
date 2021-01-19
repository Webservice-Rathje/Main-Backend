package kundenmanagement

import (
	"encoding/json"
	"github.com/Webservice-Rathje/Main-Backend/generalModels"
	"github.com/Webservice-Rathje/Main-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"net/http"
)

type CreateAccountRequestModel struct {
	UserData struct {
		Nachname      string `json:"nachname"`
		Vorname       string `json:"vorname"`
		Password      string `json:"password"`
		Telefonnummer string `json:"telefonnummer"`
		Mail          string `json:"mail"`
		Geschlecht    string `json:"geschlecht"`
		Geburtsdatum  string `json:"geburtsdatum"`
		Wohnort       string `json:"wohnort"`
		Postleitzahl  string `json:"postleitzahl"`
		Strasse       string `json:"strasse"`
		Hausnummer    string `json:"hausnummer"`
	} `json:"user_data"`
	TwoFA_Code string `json:"two_fa_code"`
}

func CreateAccountWebsocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				res, _ := json.Marshal(generalModels.ErrorResponseModel{
					Error:          err.Error(),
					CausedBy:       "Server software",
					CouldBeFixedBy: "Contact the Webservice Rathje development team",
					Alert:          "alert alert-danger",
				})
				c.WriteMessage(mt, res)
			}
			var data CreateAccountRequestModel
			err = json.Unmarshal(msg, &data)
			if err != nil {
				resp, _ := json.Marshal(generalModels.ErrorResponseModel{
					Error:          err.Error(),
					CausedBy:       "Your invalid JSON string",
					CouldBeFixedBy: "Fixing problems with your JSON string",
					Alert:          "alert alert-danger",
				})
				c.WriteMessage(mt, resp)
			}
			if !checkRequestData(data) {
				resp, _ := json.Marshal(generalModels.ErrorResponseModel{
					Error:          "Invalid request data",
					CausedBy:       "Your invalid JSON string",
					CouldBeFixedBy: "Fixing problems with your JSON string",
					Alert:          "alert alert-danger",
				})
				c.WriteMessage(mt, resp)
			}
			type response_struct struct {
				Message string `json:"message"`
				Alert   string `json:"alert"`
			}

			if data.TwoFA_Code == "null" {
				kID := utils.Generate_KundenID()
				hash := utils.HashPassword(data.UserData.Password, utils.GenerateSalt())
				conn := utils.GetConn()
				stmt, _ := conn.Prepare("SELECT * FROM `kunden` WHERE `Email`=?")
				resp, _ := stmt.Query(data.UserData.Mail)
				mail_already_exists := false
				for resp.Next() {
					mail_already_exists = true
				}
				if mail_already_exists {
					res, _ := json.Marshal(generalModels.ErrorResponseModel{
						"Für diese Email-Addresse existiert bereits ein Konto.",
						"Your Email",
						"Entering another mail",
						"alert alert-warning",
					})
					c.WriteMessage(mt, res)
				}
				stmt, _ = conn.Prepare("INSERT INTO `kunden` (`ID`, `KundenID`, `Nachname`, `Vorname`, `Password`, `Token`, `AuftragsIDs`, `Telefonnummer`, `Email`, `Geburtsdatum`, `Geschlecht`, `Wohnort`, `Postleitzahl`, `Strasse`, `Hausnummer`, `Mailverified`, `2FA`) VALUES (NULL, ?, ?, ?, ?, 'null', '', ?, ?, ?, ?, ?, ?, ?, ?, 0, 1);")
				stmt.Exec(kID, data.UserData.Nachname, data.UserData.Vorname, hash, data.UserData.Telefonnummer, data.UserData.Mail, data.UserData.Geburtsdatum, data.UserData.Geschlecht, data.UserData.Wohnort, data.UserData.Postleitzahl, data.UserData.Strasse, data.UserData.Hausnummer)
				two_fa_code := utils.Generate_2FA_Code()
				stmt, _ = conn.Prepare("INSERT INTO `2FA-Codes` (`ID`, `Code`, `KundenID`, `timestamp`) VALUES (NULL, ?, ?, CURRENT_TIMESTAMP());")
				stmt.Exec(two_fa_code, kID)
				url := "http://10.11.0.5:8081/sendRegistrationCodeMail?mail=" + data.UserData.Mail + "&name=" + data.UserData.Vorname + "%20" + data.UserData.Nachname + "&2FA-Code=" + two_fa_code
				_, err := http.Get(url)
				if err != nil {
					res, _ := json.Marshal(generalModels.ErrorResponseModel{
						err.Error(),
						"Unknown",
						"Unknown",
						"alert alert-danger",
					})
					c.WriteMessage(mt, res)
					break
				}
				res, _ := json.Marshal(response_struct{
					"Account wurde erfolgreich erstellt. Geben sie nun ihren 2FA-Code ein, den wir ihnen per Mail geschickt haben.",
					"alert alert-success",
				})
				stmt.Close()
				conn.Close()
				c.WriteMessage(mt, res)
			} else {
				conn := utils.GetConn()
				stmt, _ := conn.Prepare("SELECT `KundenID` FROM `2FA-Codes` WHERE `Code`=?")
				resp, _ := stmt.Query(data.TwoFA_Code)
				type cacheStruct struct {
					KundenID string `json:"KundenID"`
				}
				var kid_struct cacheStruct
				code_exists := false
				for resp.Next() {
					err := resp.Scan(&kid_struct.KundenID)
					if err != nil {
						res, _ := json.Marshal(generalModels.ErrorResponseModel{
							Error:          err.Error(),
							CausedBy:       "MySQL Query",
							CouldBeFixedBy: "Contact the Webservice Rathje Developer team",
							Alert:          "alert alert-danger",
						})
						c.WriteMessage(mt, res)
						break
					}
					code_exists = true
				}
				if !code_exists {
					res, _ := json.Marshal(generalModels.ErrorResponseModel{
						"Ihr 2FA-Code ist falsch. Bitte überprüfen sie, ob sie den richtigen Code eingegeben haben.",
						"wrong code",
						"check your code",
						"alert alert-waring",
					})
					c.WriteMessage(mt, res)
				} else {
					stmt, _ = conn.Prepare("UPDATE `kunden` SET `Mailverified`=1 WHERE `KundenID`=? AND `Password`=?")
					stmt.Exec(kid_struct.KundenID, utils.CheckPasswordsMatch(data.UserData.Password, conn, kid_struct.KundenID))
					stmt, _ = conn.Prepare("DELETE FROM `2FA-Codes` WHERE `Code`=?;")
					stmt.Exec(data.TwoFA_Code)
					res, _ := json.Marshal(response_struct{
						"Registrierung erfolgreich abgeschlossen. Sie können nun fortfahren.",
						"alert alert-success",
					})
					c.WriteMessage(mt, res)
					break
				}
			}
		}
	})
}

func checkRequestData(data CreateAccountRequestModel) bool {
	if data.TwoFA_Code != "" && data.UserData.Nachname != "" && data.UserData.Vorname != "" &&
		data.UserData.Password != "" && data.UserData.Telefonnummer != "" && data.UserData.Mail != "" &&
		data.UserData.Geschlecht != "" && data.UserData.Geburtsdatum != "" && data.UserData.Wohnort != "" &&
		data.UserData.Postleitzahl != "" && data.UserData.Strasse != "" && data.UserData.Hausnummer != "" {
		return true
	} else {
		return false
	}
}
