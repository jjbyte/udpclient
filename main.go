package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"udpclient/config"
	"udpclient/controllers"
)

func udpProcess(conn *net.UDPConn,clientAddr *net.UDPAddr,ownerName string) {
	// 最大读取数据大小
	data := make([]byte, 1024)
	for {
		//读取UDP报文
		len, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("failed read udp msg, error: " + err.Error())
			return
		}

		//解析读取到的报文
		var rsp controllers.RspJson
		err = json.Unmarshal(data[:len],&rsp)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Response:",rsp)

		if rsp.Result {
			switch v := (rsp.Msg).(type) {
			case map[string]interface{}:
				msg, ok := (rsp.Msg).(map[string]interface{})
				if ok {
					conn.Close()
					controllers.BidirectionalHole(clientAddr,msg["addr"].(string),ownerName)
				}
			default:
				log.Println("msg:",v)
			}
		}
	}
}

func main() {

	addr := config.Conf.Get("common.Addr").(string)
	ownerName := os.Args[1]
	var otherName string
	if len(os.Args) > 2 {
		otherName = os.Args[2]
	}

	clientAddr := &net.UDPAddr{
		IP: net.IPv4zero,
		Port: 6688,
	}
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	conn,err := net.DialUDP("udp",clientAddr,serverAddr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	req := &controllers.ReqJson{
		Command:"1",
		OwnerName:ownerName,
		EndName:otherName,
	}
	rspBuff, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = conn.Write(rspBuff)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go udpProcess(conn,clientAddr,ownerName)

	//如果有对端用户，则发送获取对端用户信息
	if len(otherName) > 0 {
		req = &controllers.ReqJson{
			Command:"3",
			OwnerName:ownerName,
			EndName:otherName,
		}
		rspBuff, err = json.Marshal(req)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		_, err = conn.Write(rspBuff)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	select {

	}

}
