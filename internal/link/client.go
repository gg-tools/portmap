package link

import (
	"github.com/gg-tools/portmap/util"
	"log"
	"net"
	"time"
)

type Client struct {
	remote string
	local  string
	pwd    string
}

func (c *Client) Connect() {
	proxy, err := net.DialTimeout("tcp", c.remote, 5*time.Second)
	if err != nil {
		log.Println("CAN'T CONNECT:", c.remote, " err:", err)
		return
	}
	defer proxy.Close()
	util.WriteString(proxy, c.pwd+"\n"+util.C2PConnect)

	for {
		proxy.SetReadDeadline(time.Now().Add(2 * time.Second))
		msg, err := util.ReadString(proxy)
		//	proxy.SetReadDeadline(time.Time{})
		if err == nil {
			if msg == util.P2CNewSession {
				go c.session()
			} else {
				log.Println(msg)
			}
		} else {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				//log.Println("Timeout")
				proxy.SetWriteDeadline(time.Now().Add(2 * time.Second))
				_, werr := util.WriteString(proxy, util.C2PKeepAlive) //send KeepAlive msg
				if werr != nil {
					log.Println("CAN'T WRITE, err:", werr)
					return
				}

				continue
			} else {
				log.Println("SERVER CLOSE, err:", err)
				return
			}
		}
		//time.Sleep(2*time.Second)
	}
}

//客户端单次连接处理
func (c *Client) session() {
	log.Println("Create Session")
	rp, err := net.Dial("tcp", c.remote)
	if err != nil {
		log.Println("Can't' connect:", c.remote, " err:", err)
		return
	}
	//defer util.CloseConn(rp)
	util.WriteString(rp, c.pwd+"\n"+util.C2PSession)
	lp, err := net.Dial("tcp", c.local)
	if err != nil {
		log.Println("Can't' connect:", c.local, " err:", err)
		rp.Close()
		return
	}
	go util.CopyFromTo(rp, lp, nil)
	go util.CopyFromTo(lp, rp, nil)
}
