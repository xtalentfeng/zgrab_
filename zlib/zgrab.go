/*
 * ZGrab Copyright 2015 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 */

package zlib

import (
	"encoding/json"
	"net"
	"time"

	"github.com/zmap/zgrab/ztools/ssh"
	"github.com/zmap/zgrab/ztools/ztls"
)

type Grab struct {
	IP             net.IP
	Domain         string
	Time           time.Time
	Data           GrabData
	Error          error
	ErrorComponent string
}

type encodedGrab struct {
	IP             string    `json:"ip"`
	Domain         string    `json:"domain,omitempty"`
	Time           string    `json:"timestamp"`
	Data           *GrabData `json:"data,omitempty"`
	Error          *string   `json:"error,omitempty"`
	ErrorComponent string    `json:"error_component,omitempty"`
}
//modbus
type GrabData struct {
	Banner       string                `json:"banner,omitempty"`
	Read         string                `json:"read,omitempty"`
	Write        string                `json:"write,omitempty"`
	EHLO         string                `json:"ehlo,omitempty"`
	SMTPHelp     *SMTPHelpEvent        `json:"smtp_help,omitempty"`
	StartTLS     string                `json:"starttls,omitempty"`
	TLSHandshake *ztls.ServerHandshake `json:"tls,omitempty"`
	HTTP         *HTTPRequestResponse  `json:"http,omitempty"`
	Heartbleed   *ztls.Heartbleed      `json:"heartbleed,omitempty"`
	Modbus       *ModbusEvent          `json:"modbus,omitempty"`
	//add s7
	S7			 *TKPTEvent			   `json:"s7,omitempty"`
	SSH          *ssh.HandshakeLog     `json:"ssh,omitempty"`
}

func (g *Grab) MarshalJSON() ([]byte, error) {
	time := g.Time.Format(time.RFC3339)
	var errString *string
	if g.Error != nil {
		s := g.Error.Error()
		errString = &s
	}
	obj := encodedGrab{
		IP:             g.IP.String(),
		Domain:         g.Domain,
		Time:           time,
		Data:           &g.Data,
		Error:          errString,
		ErrorComponent: g.ErrorComponent,
	}
	return json.Marshal(obj)
}

func (g *Grab) UnmarshalJSON(b []byte) error {
	eg := new(encodedGrab)
	err := json.Unmarshal(b, eg)
	if err != nil {
		return err
	}
	g.IP = net.ParseIP(eg.IP)
	g.Domain = eg.Domain
	if g.Time, err = time.Parse(time.RFC3339, eg.Time); err != nil {
		return err
	}
	panic("unimplemented")
}

func (g *Grab) status() status {
	if g.Error != nil {
		return status_failure
	}
	return status_success
}
