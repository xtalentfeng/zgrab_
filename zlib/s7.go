/*
 * ZGrab Copyright 2015 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 */

package zlib

import (
	"bytes"
	"encoding/binary"
//	"encoding/json"
	"fmt"
//	"strconv"
)

//Request
var TKPTHeaderBytes = []byte{
	0x03,0x00,
}

type TKPTRequest struct {
	version byte
	reserved byte
	size int
	data []byte
}

type TKPTResponse struct {
	version byte
	reserved byte
	size int
	data []byte
}
type COTPConnectionRequest struct {
	size byte
	pdu_type byte
	dst_ref int
	src_ref int
	flag byte
	parameter_code0 byte
	parameter_length0 byte
	src_tsap int
	parameter_code1 byte
	parameter_length1 byte
	dst_tsap int
	parameter_code2 byte
	parameter_length2 byte
	tpdu_size byte
}

func (r *TKPTRequest) MarshalBinary() (data []byte, err error){
	data = make([]byte,4+len(r.data))
	data[0] = r.version
	data[1] = r.reserved
	binary.BigEndian.PutUint16(data[2:4], uint16(len(r.data)+4))
	copy(data[4:],r.data)
	return
}

func (r *COTPConnectionRequest) MarshalBinary() (data []byte, err error){
	data = make([]byte,18)
	data[0] = r.size
	data[1] = r.pdu_type
	binary.BigEndian.PutUint16(data[2:4], uint16(r.dst_ref))
	binary.BigEndian.PutUint16(data[4:6], uint16(r.src_ref))
	data[6] = r.flag
	data[7] = r.parameter_code0
	data[8] = r.parameter_length0
	binary.BigEndian.PutUint16(data[9:11], uint16(r.src_tsap))
	data[11] = r.parameter_code1
	data[12] = r.parameter_length1
	binary.BigEndian.PutUint16(data[13:15], uint16(r.dst_tsap))
	data[15] = r.parameter_code2
	data[16] = r.parameter_length2
	data[17] = r.tpdu_size
	return
}

//response


type TKPTEvent struct {
	version byte `json:"version"`
	reserved byte `json:"reserved"`
	size int `json:"length"`
	pdu_type byte `json:"pdu_type"`
	dst_tsap int `json:"dst_tsap"`
	src_tsap int `json:"src_tsap"`
	Response []byte `json:"raw_response,omitempty"`
}

func (m *TKPTEvent) parseSelf(){
	if len(m.Response) < 18 {
		return
	}
	if m.Response[1] != 0x0e {
		return
	} 
	m.pdu_type = m.Response[1]
	m.dst_tsap = int (binary.BigEndian.Uint16(m.Response[13:15]))
	m.src_tsap = int (binary.BigEndian.Uint16(m.Response[9:11]))
}

func (c *Conn) ReadMin1(res []byte, bytes int) (cnt int, err error) {
	for cnt < bytes {
		var n int
		n, err = c.getUnderlyingConn().Read(res[cnt:])
		cnt += n

		if err != nil && cnt >= len(res) {
			err = fmt.Errorf("S7: response buffer too small")
		}

		if err != nil {
			return
		}
	}

	return
}



func (c *Conn) GetS7Response() (res TKPTResponse, err error){
	var cnt int
	buf := make([]byte,1024)
	header := buf[0:4]
	buf = buf[4:]

	cnt, err = c.ReadMin1(header,4)
	if err != nil{
		err = fmt.Errorf("modbus:could not get response: %s", err.Error())
		return
	}

	//first 2 bytes should be known, verify them
	if !bytes.Equal(header[0:2],TKPTHeaderBytes) {
		err = fmt.Errorf("s7: not a s7 response")
		return
	}
	msglen := int (binary.BigEndian.Uint16(header[2:4]))

	cnt = 0
	if msglen >len(buf){
		msglen = len(buf)
	}//

	//one of bytes in length counts as part of header
	for cnt < msglen-1 {
		var n int
		n, err = c.getUnderlyingConn().Read(buf[cnt:])
		cnt += n

		if err != nil && cnt >= len(buf){
			err = fmt.Errorf("s7: response buffer too small")

		}

		if err != nil{
			break
		}
	}

	if cnt > len(buf) {
		cnt = len(buf)
	}

	var d []byte
	if cnt > 1 {
		d = buf[0:cnt]
	}

	res = TKPTResponse{
		version: 0x03,
		reserved: 0x00,
		size: msglen,
		data: d,
	}
	return
}



//应该放到conn.go的SendS7Echo()
/*
func (c *Conn) SendS7Echo() (int,error){

	COTP_req := COTPConnectionRequest{
		size:17,
		pdu_type:0x0e,
		dst_ref:0,
		src_ref:0x04, //
		flag:0,
		parameter_code0:0xc1,
		parameter_length0:2,
		src_tsap:0x100,
		parameter_code1:0xc2
		parameter_length1:2,
		dst_tsap:0x102,
		parameter_code2:0xc0,
		parameter_length2:1,
		tpdu_size:0x0a,
	}
	req := TKPTRequest{
		version:3,
		reserved:0,
		size:len(COTP_req.MarshalBinary()+4),
		data:COTP_req.MarshalBinary(),
	}

	event := new(TKPTEvent)
	data,err := req.MarshalBinary()
	w :=0
	for w < len(data) {
		wiritten, err = c.getUnderlyingConn().Write(data[w:])
		w += written
		if err != nil {
			c.grabData.S7 = event
			return w, errors.New("Could not write s7TKTPT request")
		}
	}
	res, err := c.GetS7Response()
	event.version = res.version
	event.reserved = res.reserved
	event.size = res.size
	event.Response = res.Data
	event.parseSelf()
	//  
	c.grabData.S7 = event
	return w, err

}
*/
