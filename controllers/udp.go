package controllers

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type ClientInfo struct {
	Addr		string `json:"addr"`
	ClientName	string `json:"name"`
}

type ReqJson struct {
	Command 	string `json:"command"`
	OwnerName	string `json:"owner"`
	EndName		string `json:"end"`
}

type RspJson struct {
	Result 		bool `json:"result"`
	Msg			interface{} `json:"msg"`
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func BidirectionalHole(srcAddr *net.UDPAddr, dstAddr string,ownerName string) {

	anotherAddr := parseAddr(dstAddr)
	conn, err := net.DialUDP("udp", srcAddr, &anotherAddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),
	//用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息

	if _,err = conn.Write([]byte("handshake"));err != nil {
		log.Println("send handshake:", err)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err := conn.Write([]byte("from [" + ownerName + "]"));err != nil {
				log.Println("send msg fail", err)
				continue
			}
			log.Println("send a msg ok.")
		}
	}()

	for {
		data := make([]byte, 1024)
		len, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("recv a msg:%s\n", data[:len])
		}
	}
}
