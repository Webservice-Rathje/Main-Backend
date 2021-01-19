package kundenmanagement

import (
	"encoding/json"
	"github.com/Webservice-Rathje/Main-Backend/generalModels"
	"github.com/Webservice-Rathje/Main-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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
					err.Error(),
					"Server software",
					"Contact the Webservice Rathje development team",
					"alert alert-danger",
				})
				c.WriteMessage(mt, res)
				break
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
			type response_struct struct {
				Message string `json:"message"`
				Alert   string `json:"alert"`
			}
			if data.TwoFA_Code == "null" {
				kID := utils.Generate_KundenID()
				hash := utils.HashPassword(data.UserData.Password, utils.GenerateSalt())
				conn := utils.GetConn()
				stmt, _ := conn.Prepare("INSERT INTO `kunden` (`ID`, `KundenID`, `Nachname`, `Vorname`, `Password`, `Token`, `AuftragsIDs`, `Telefonnummer`, `Email`, `Geburtsdatum`, `Geschlecht`, `Wohnort`, `Postleitzahl`, `Strasse`, `Hausnummer`, `Mailverified`, `2FA`) VALUES (NULL, ?, ?, ?, ?, 'null', '', ?, ?, ?, ?, ?, ?, ?, ?, 0, 1);")
				stmt.Exec(kID, data.UserData.Nachname, data.UserData.Vorname, hash, data.UserData.Telefonnummer, data.UserData.Mail, data.UserData.Geburtsdatum, data.UserData.Geschlecht, data.UserData.Wohnort, data.UserData.Postleitzahl, data.UserData.Strasse, data.UserData.Hausnummer)
				two_fa_code := utils.Generate_2FA_Code()
				stmt, _ = conn.Prepare("INSERT INTO `2FA-Codes` (`ID`, `Code`, `KundenID`, `timestamp`) VALUES (NULL, ?, ?, CURRENT_TIMESTAMP());")
				stmt.Exec(two_fa_code, kID)
				// Sending Mail via registration service
				resp, _ := json.Marshal(response_struct{
					"Account wurde erfolgreich erstellt. Geben sie nun ihren 2FA-Code ein, den wir ihnen per Mail geschickt haben.",
					"alert alert-success",
				})
				stmt.Close()
				conn.Close()
				c.WriteMessage(mt, resp)
			} else {
				conn := utils.GetConn()
				stmt, _ := conn.Prepare("SELECT `KundenID` FROM `2FA-Codes` WHERE `Code`=?")
				resp, _ := stmt.Query(data.TwoFA_Code)
				type cacheStruct struct {
					KundenID string `json:"KundenID"`
				}
				var kid_struct cacheStruct
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
				}
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
	})
}
