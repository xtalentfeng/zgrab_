/*
 * ZGrab Copyright 2015 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 */

package zlib

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/xtalentfeng/zgrab_/ztools/ssh"
	"github.com/xtalentfeng/zgrab_/ztools/x509"
	"github.com/xtalentfeng/zgrab_/ztools/zlog"
)

type HTTPConfig struct {
	Method      string
	Endpoint    string
	UserAgent   string
	ProxyDomain string
	MaxSize     int
}

type SSHScanConfig struct {
	SSH               bool
	Client            string
	KexAlgorithms     string
	HostKeyAlgorithms string
}

func (sc *SSHScanConfig) GetClientImplementation() (*ssh.ClientImplementation, bool) {
	if sc.Client == "" {
		return &ssh.OpenSSH_6_6p1, true
	}
	return ssh.ClientImplementationByName(sc.Client)
}

func (sc *SSHScanConfig) readNameList(reader io.Reader) (ssh.NameList, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, errors.New("TAHW")
	}
	nameList := ssh.NameList(records[0])
	return nameList, nil
}

func (sc *SSHScanConfig) MakeKexNameList() (ssh.NameList, error) {
	if sc.KexAlgorithms == "" {
		c, _ := sc.GetClientImplementation()
		return c.KexAlgorithms(), nil
	}
	kexReader := strings.NewReader(sc.KexAlgorithms)
	return sc.readNameList(kexReader)
}

func (sc *SSHScanConfig) MakeHostKeyNameList() (ssh.NameList, error) {
	if sc.HostKeyAlgorithms == "" {
		c, _ := sc.GetClientImplementation()
		return c.HostKeyAlgorithms(), nil
	}
	hostKeyReader := strings.NewReader(sc.HostKeyAlgorithms)
	return sc.readNameList(hostKeyReader)
}

func (sc *SSHScanConfig) MakeConfig() *ssh.Config {
	config := new(ssh.Config)
	config.KexAlgorithms, _ = sc.MakeKexNameList()
	config.HostKeyAlgorithms, _ = sc.MakeHostKeyNameList()
	return config
}

type Config struct {
	// Connection
	Port               uint16
	Timeout            time.Duration
	Senders            uint
	ConnectionsPerHost uint

	// Encoding
	Encoding string

	// TLS
	TLS               bool
	TLSVersion        uint16
	Heartbleed        bool
	RootCAPool        *x509.CertPool
	DHEOnly           bool
	ExportsOnly       bool
	ExportsDHOnly     bool
	FirefoxOnly       bool
	FirefoxNoDHE      bool
	ChromeOnly        bool
	ChromeNoDHE       bool
	SafariOnly        bool
	SafariNoDHE       bool
	NoSNI             bool
	TLSExtendedRandom bool

	// SSH
	SSH SSHScanConfig

	// Banners and Data
	Banners  bool
	SendData bool
	Data     []byte
	Raw      bool

	// Mail
	SMTP       bool
	IMAP       bool
	POP3       bool
	SMTPHelp   bool
	EHLODomain string
	EHLO       bool
	StartTLS   bool

	// FTP
	FTP bool

	// Modbus
	Modbus bool

	// s7
	S7 bool

	// HTTP
	HTTP HTTPConfig

	// Error handling
	ErrorLog *zlog.Logger

	// Go Runtime Config
	GOMAXPROCS int
}
