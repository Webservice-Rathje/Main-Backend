package endpoints

import (
	"encoding/json"
	"github.com/Webservice-Rathje/Main-Backend/generalModels"
	"github.com/Webservice-Rathje/Main-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type CheckTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func CheckTokenController(c *fiber.Ctx) error {
	obj := CheckTokenRequest{}
	err := json.Unmarshal(c.Body(), &obj)
	if err != nil {
		return c.JSON(generalModels.ErrorResponseModel{
			Error:          err.Error(),
			CausedBy:       "Error in JSON syntax",
			CouldBeFixedBy: "Improving your JSON syntax",
			Alert:          "alert alert-danger",
		})
	}
	if !checkCheckTokenRequest(obj) {
		return c.JSON(generalModels.ErrorResponseModel{
			Error:          "Wrong JSON syntax",
			CausedBy:       "Error in JSON syntax",
			CouldBeFixedBy: "Improving your JSON syntax",
			Alert:          "alert alert-danger",
		})
	}
	conn := utils.GetConn()
	stmt, _ := conn.Prepare("SELECT `KundenID`, `Password` FROM `kunden` WHERE `Email` =? AND `Token`=?;")
	type cacheStruct struct {
		KundenID string `json:"KundenID"`
		Password string `json:"Password"`
	}
	resp, _ := stmt.Query(obj.Email, obj.Token)
	var user cacheStruct
	for resp.Next() {
		err = resp.Scan(&user.KundenID, &user.Password)
		if err != nil {
			panic(err)
		}
	}
	defer resp.Close()
	status := utils.CheckPasswordsMatch(obj.Password, conn, user.KundenID)
	defer stmt.Close()
	defer conn.Close()
	if status {
		return c.JSON(LoginResponse{
			"Anmelden erfolgreich",
			"alert alert-success",
			"null",
		})
	} else {
		return c.JSON(LoginResponse{
			"Die Anmeldedaten sind leider falsch.",
			"alert alert-danger",
			"null",
		})
	}
}

func checkCheckTokenRequest(obj CheckTokenRequest) bool {
	return obj.Email != "" && obj.Password != "" && obj.Token != ""
}
