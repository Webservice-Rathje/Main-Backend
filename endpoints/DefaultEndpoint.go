package endpoints

import "github.com/gofiber/fiber/v2"

type DefaultEndpointResponse struct {
	ApplicationName                string `json:"application-name"`
	ApplicationVersion             string `json:"application-version"`
	ApplicationLanguage            string `json:"application-language"`
	ApplicationMode                string `json:"application-mode"`
	ApplicationProgrammingLanguage string `json:"application-programming-language"`
	ApplicationHost                string `json:"application-host"`
}

func DefaultEndpoint(c *fiber.Ctx) error {
	return c.JSON(DefaultEndpointResponse{
		"Webservice Rathje backend",
		"v0.0.1-dev",
		"de-DE",
		"development",
		"golang",
		"Webservice Rathje",
	})
}
