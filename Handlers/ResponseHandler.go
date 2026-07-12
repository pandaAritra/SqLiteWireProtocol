package handlers

import (
	"encoding/json"
	"net"

	utils "github.com/pandaAritra/sqliteWireProtocol/Utils"
	"github.com/pandaAritra/sqliteWireProtocol/models"
)

// SendResponse marshals the models.Response to JSON and sends it over the connection
// using the wire protocol format (message type 0x03 for Response).
func SendResponse(client net.Conn, response *models.Response) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return err
	}

	packet := utils.MakePKT(0x03, string(jsonBytes))
	_, err = client.Write(*packet)
	return err
}
