package constants

import "os"

const ApiVersion int = 1
const AppName string = "soapstone"

var JwtKey string = os.Getenv("JWT_AUTH_KEY")
