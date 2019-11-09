package SendUtils

import (
	"../log"
	"io"
	"net"
	"time"
)

func CloseSocket() {

}

func main() {

}
func checkerror(err error, errortype string) bool {
	if errortype == "readbyte" {
		if err != nil && err != io.EOF { //io.EOF在网络编程中表示对端把链接关闭了。
			log.Error(err)
			log.Debug("Maybe there is no response!")
			return false
		} else {
			//log.Info(string(buf[:a]))
			return true
		}
	} else {
		if err != nil {
			log.Error("error at ", errortype, "Err occur at :", err)
			return false
		} else {
			return true
		}
	}
}
func OpenSocket(hex []byte) {
	coon, err := net.Dial("tcp", "210.34.130.61:7777")
	log.Info("Connect sucess!")
	if checkerror(err, "TCP") {
		n, err := coon.Write(hex)
		//err = coon.Close()
		if checkerror(err, "Senderror") {

			log.Info("Alerady sent ", n, " bytes")
			log.Debug("Send Buffer was",hex)
			buf := make([]byte, 1024)
			err = coon.SetReadDeadline(time.Now().Add(time.Millisecond * 1000))
			if checkerror(err, "Deanline") {
				_, err := coon.Read(buf)
				//TODO 加入若没有反应则断开
				if checkerror(err, "readbyte") {
					//io.EOF在网络编程中表示对端把链接关闭了。
					log.Info(string(buf))
				}
				err = coon.Close()
				if checkerror(err, "TCP Close") {

				} else {
					panic("Close error!!!")
				}
				log.Info("Connection Closed!")
			}
		}
	}
}
