package tcpserver

import "log"

// ------------------------------------------------------------------------------
// ManageRx
// ------------------------------------------------------------------------------
func (a *SStServer) ManageRx() bool {
	log.Printf("ManageRx[%d]", a.Id)

	rxbuf := make([]byte, 1024)

	n, err := a.Tcp.Read(rxbuf)
	if err != nil {
		if err.Error() != "EOF" {
			log.Printf("Rx Error (%d) [%d bytes]: %v", a.Id, n, err)
		} else {
			log.Printf("Client %d closed connection", a.Id)
		}
		a.close()
		return false
	}

	log.Printf("Received [%d bytes] from %d: %s", n, a.Id, string(rxbuf[:n]))
	return true
}
