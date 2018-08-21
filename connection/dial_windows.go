package connection

import (
	"net"
	"time"

	winio "github.com/Microsoft/go-winio"
)

func dial(sockname string, timeout time.Duration) (net.Conn, error) {
	return winio.DialPipe(sockname, &timeout)
}
