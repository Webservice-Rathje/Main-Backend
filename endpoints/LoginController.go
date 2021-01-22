package endpoints

import (
	"encoding/json"
	"github.com/Webservice-Rathje/Main-Backend/generalModels"
	"github.com/Webservice-Rathje/Main-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Alert   string `json:"alert"`
	Token   string `json:"token"`
}

func LoginController(c *fiber.Ctx) error {
	obj := LoginRequest{}
	err := json.Unmarshal(c.Body(), &obj)
	if err != nil {
		return c.JSON(generalModels.ErrorResponseModel{
			err.Error(),
			"Error in JSON syntax",
			"Improving your JSON syntax",
			"alert alert-danger",
		})
	}
	conn := utils.GetConn()
	stmt, _ := conn.Prepare("SELECT `KundenID`, `Password` FROM `kunden` WHERE `Email` =?;")
	type cacheStruct struct {
		KundenID string `json:"KundenID"`
		Password string `json:"Password"`
	}
	resp, _ := stmt.Query(obj.Email)
	var user cacheStruct
	for resp.Next() {
		err = resp.Scan(&user.KundenID, &user.Password)
		if err != nil {
			panic(err)
		}
	}
	defer resp.Close()
	status := utils.CheckPasswordsMatch(obj.Password, conn, user.KundenID)
	if status {
		stmt, _ := conn.Prepare("UPDATE `kunden` SET `Token`=? WHERE `KundenID`=?;")
		token := utils.GenerateToken()
		stmt.Exec(token, user.KundenID)
		defer stmt.Close()
		defer conn.Close()
		return c.JSON(LoginResponse{
			"Anmelden erfolgreich",
			"alert alert-success",
			token,
		})
	} else {
		defer stmt.Close()
		defer conn.Close()
		return c.JSON(LoginResponse{
			"Die Anmeldedaten sind leider falsch.",
			"alert alert-danger",
			"None",
		})
	}
}
