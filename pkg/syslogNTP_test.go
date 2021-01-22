// MIT License
// 
// (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
// 
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.


package bmc_nwprotocol

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"reflect"
	"time"
	"strings"
	"os"
)

var __global_t *testing.T
var nwProtoGlobal RedfishNWProtocol
var stallG bool
var statusCodeG int

func clearNWP(nwp *NWPData) {
	nwp.SyslogSpec    = ""
	nwp.NTPSpec       = ""
	nwp.SSHKey        = ""
	nwp.SSHConsoleKey = ""
	nwp.BootOrder = []string{}
}

func printnwp(nwp RedfishNWProtocol) string {
	var rstr1 = ""
	var rstr2 = ""

	if (nwp.Oem != nil) {
		rstr1 = fmt.Sprintf("Oem.Syslog.ProtcolEnabled: %t, Oem.Syslog.SyslogServers: %q, Oem.Syslog.Transport: %q, Oem.Syslog.Port: %d",
					nwp.Oem.Syslog.ProtocolEnabled,
					nwp.Oem.Syslog.SyslogServers,
					nwp.Oem.Syslog.Transport,
					nwp.Oem.Syslog.Port)
		if (nwp.Oem.SSHAdmin != nil) {
			rstr1 += fmt.Sprintf(" SSHAdmin: '%s'",nwp.Oem.SSHAdmin.AuthorizedKeys)
		}
		if (nwp.Oem.SSHConsole != nil) {
			rstr1 += fmt.Sprintf(" SSHConsole: '%s'",nwp.Oem.SSHConsole.AuthorizedKeys)
		}
	}
	if (nwp.NTP != nil) {
		rstr2 = fmt.Sprintf("NTP.NTPServers: %q, NTP.ProtocolEnabled: %t, NTP.Port: %d",
					nwp.NTP.NTPServers,
					nwp.NTP.ProtocolEnabled,
					nwp.NTP.Port)
	}
	rstr := fmt.Sprintf("valid: %t, %s %s",nwp.valid,rstr1,rstr2)
	return rstr
}

func verifyNWP(t *testing.T, rfp RedfishNWProtocol) {
	ba,baerr := json.Marshal(&rfp)
	if (baerr != nil) {
		t.Error("ERROR marshalling RF NWP data:",baerr)
		return
	}
	rstr := string(ba)

	if (rfp.NTP == nil) {
		if (strings.Contains(rstr,"NTP")) {
			t.Errorf("ERROR, marshalled JSON with no NTP info contains NTP entry: '%s'",
				rstr)
		}
	} else {
		if ((len(rfp.NTP.NTPServers) == 0) && strings.Contains(rstr,"NTPServers")) {
			t.Errorf("ERROR< marshalled JSON with no NTP servers contains NTPServers entry: '%s'",
				rstr)
		}
	}

	if (rfp.Oem == nil) {
		if (strings.Contains(rstr,"Oem")) {
			t.Errorf("ERROR, marshalled JSON with no Oem info contains Oem entry: '%s'",
				rstr)
		}
	} else {
		if ((rfp.Oem.Syslog == nil) && (strings.Contains(rstr,"Syslog"))) {
			t.Errorf("ERROR, marshalled JSON with no Syslog info contains Syslog entry: '%s'",
				rstr)
		}
		if ((rfp.Oem.Syslog != nil) && (len(rfp.Oem.Syslog.SyslogServers) == 0) &&
			strings.Contains(rstr,"SyslogServer")) {
			t.Errorf("ERROR, marshalled JSON with no SyslogServer info contains SyslogServer entry: '%s'",
				rstr)
		}
		if ((rfp.Oem.SSHAdmin == nil) && strings.Contains(rstr,"SSHAdmin")) {
			t.Errorf("ERROR, marshalled JSON with no SSHAdmin info contains SSHAdmin entry: '%s'",
				rstr)
		}
		if ((rfp.Oem.SSHConsole == nil) && strings.Contains(rstr,"SSHConsole")) {
			t.Errorf("ERROR, marshalled JSON with no SSHConsole info contains SSHConsole entry: '%s'",
				rstr)
		}
	}
}

func Test_makeNPData(t *testing.T) {
	var err error
	var emsgBase,rfSuffix string
	var nwProtoInfo RedfishNWProtocol
	var rfCopy RedfishNWProtocol
	var expNWP1 = RedfishNWProtocol{valid: true,
                                   }

	var expNWP2 = RedfishNWProtocol{valid: true,
	                                Oem: &OemData{Syslog: &SyslogData{ProtocolEnabled: true,
	                                                                 SyslogServers: []string{"10.11.12.13","110.111.112.113",},
	                                                                 Transport: "udp",
	                                                                 Port: 123,
	                                                                },
	                                             },
	                               }
	var expNWP3 = RedfishNWProtocol{valid: true,
	                                NTP: &NTPData{NTPServers: []string{"120.121.122.123","130.131.132.133"},
	                                              ProtocolEnabled: true,
	                                              Port: 789,
	                                             },
	                               }
	var expNWP4 = RedfishNWProtocol{valid: true,
	                                Oem: &OemData{Syslog: &SyslogData{ProtocolEnabled: true,
	                                                                 SyslogServers: []string{"10.11.12.13","110.111.112.113",},
	                                                                 Transport: "udp",
	                                                                 Port: 123,
	                                                                },
	                                             },
	                               NTP: &NTPData{NTPServers: []string{"120.121.122.123","130.131.132.133"},
	                                             ProtocolEnabled: true,
	                                             Port: 789,
	                                            },
	                               }
	var expNWP4a = RedfishNWProtocol{valid: true,
	                                Oem: &OemData{Syslog: &SyslogData{ProtocolEnabled: true,
	                                                                 SyslogServers: []string{"10.11.12.13","110.111.112.113",},
	                                                                 Transport: "udp",
	                                                                 Port: 123,
	                                                                },
	                                              SSHAdmin: &SSHAdminData{"aabbccdd"},
	                                              SSHConsole: &SSHAdminData{"eeffgghh"},
	                                             },
	                               NTP: &NTPData{NTPServers: []string{"120.121.122.123","130.131.132.133"},
	                                             ProtocolEnabled: true,
	                                             Port: 789,
	                                            },
	                               }

	var expNWP5 = RedfishNWProtocol{valid: true,
	                                Oem: &OemData{SSHAdmin: &SSHAdminData{AuthorizedKeys: "aabbccdd11",},},
	                               }

	var expNWP6 = RedfishNWProtocol{valid: true,
	                                Oem: &OemData{SSHConsole: &SSHAdminData{AuthorizedKeys: "wwxxyyzz",},},
	                               }

	///////////// Init() stuff ////////////////////////

	//First test no settings.  Should be OK but result in invalid 
	//nwProtoInfo variable

	var nwp NWPData

	clearNWP(&nwp)
	rfSuffix = ""

	emsgBase = "Init() with no parameters"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Error("No params Init() error:",err)
	}
	if (!reflect.DeepEqual(nwProtoInfo,expNWP1)) {
		t.Errorf("%s: data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP1),printnwp(nwProtoInfo))
	}
	verifyNWP(t,nwProtoInfo)
	rfCopy = CopyRFNetworkProtocol(&nwProtoInfo)
	if (!reflect.DeepEqual(rfCopy,expNWP1)) {
		t.Errorf("%s: rfCopy data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP1),printnwp(rfCopy))
	}

	//Check syslog info only

	emsgBase = "Init() with only syslogTarg"
	clearNWP(&nwp)
	nwp.SyslogSpec = expNWP2.Oem.Syslog.SyslogServers[0] + "," +
	                 expNWP2.Oem.Syslog.SyslogServers[1] + ":" +
	                 fmt.Sprintf("%d",expNWP2.Oem.Syslog.Port)
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Errorf(emsgBase,"error:",err)
	}
	if (!reflect.DeepEqual(nwProtoInfo,expNWP2)) {
		t.Errorf("%s: data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP2),printnwp(nwProtoInfo))
	}
	verifyNWP(t,nwProtoInfo)
	rfCopy = CopyRFNetworkProtocol(&nwProtoInfo)
	if (!reflect.DeepEqual(rfCopy,expNWP2)) {
		t.Errorf("%s: rfCopy data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP2),printnwp(rfCopy))
	}

	//Check NTP info only

	emsgBase = "Init() with only ntpTarg"
	clearNWP(&nwp)
	nwp.NTPSpec = expNWP3.NTP.NTPServers[0] + "," + expNWP3.NTP.NTPServers[1] +
	              ":" + fmt.Sprintf("%d",expNWP3.NTP.Port)
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Error(emsgBase,"error:",err)
	}
	if (!reflect.DeepEqual(nwProtoInfo,expNWP3)) {
		t.Errorf("%s: data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP3),printnwp(nwProtoInfo))
	}
	verifyNWP(t,nwProtoInfo)
	rfCopy = CopyRFNetworkProtocol(&nwProtoInfo)
	if (!reflect.DeepEqual(rfCopy,expNWP3)) {
		t.Errorf("%s: rfCopy data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP3),printnwp(rfCopy))
	}

	//Check both syslog and NTP, and rf suffix

	emsgBase = "Init() with syslogTarg and ntpTarg"
	clearNWP(&nwp)
	nwp.SyslogSpec = expNWP4.Oem.Syslog.SyslogServers[0] + "," +
	                 expNWP4.Oem.Syslog.SyslogServers[1] + ":" +
	                 fmt.Sprintf("%d",expNWP4.Oem.Syslog.Port)
	nwp.NTPSpec = expNWP4.NTP.NTPServers[0] + "," + expNWP4.NTP.NTPServers[1] +
	              ":" + fmt.Sprintf("%d",expNWP4.NTP.Port)
	rfSuffix = "/a/b/c"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Error(emsgBase,"error:",err)
	}
	if (!reflect.DeepEqual(nwProtoInfo,expNWP4)) {
		t.Errorf("%s: data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP4),printnwp(nwProtoInfo))
	}
	if (redfishNPSuffix != rfSuffix) {
		t.Errorf("Redfish suffix not applied to library!\n")
	}
	verifyNWP(t,nwProtoInfo)
	rfCopy = CopyRFNetworkProtocol(&nwProtoInfo)
	if (!reflect.DeepEqual(rfCopy,expNWP4)) {
		t.Errorf("%s: rfCopy data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP4),printnwp(rfCopy))
	}

	//Check SSHKey only

	emsgBase = "Init() with only SSHKey"
	clearNWP(&nwp)
	nwp.SSHKey = "aabbccdd11"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Error(emsgBase,"error:",err)
	}
	if (!reflect.DeepEqual(nwProtoInfo,expNWP5)) {
		t.Errorf("%s: data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP5),printnwp(nwProtoInfo))
	}
	verifyNWP(t,nwProtoInfo)
	rfCopy = CopyRFNetworkProtocol(&nwProtoInfo)
	if (!reflect.DeepEqual(rfCopy,expNWP5)) {
		t.Errorf("%s: rfCopy data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP5),printnwp(rfCopy))
	}

	//Check SSHConsoleKey only

	emsgBase = "Init() with only SSHConsoleKey"
	clearNWP(&nwp)
	nwp.SSHConsoleKey = "wwxxyyzz"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Error(emsgBase,"error:",err)
	}
	if (!reflect.DeepEqual(nwProtoInfo,expNWP6)) {
		t.Errorf("%s: data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP5),printnwp(nwProtoInfo))
	}
	verifyNWP(t,nwProtoInfo)
	rfCopy = CopyRFNetworkProtocol(&nwProtoInfo)
	if (!reflect.DeepEqual(rfCopy,expNWP6)) {
		t.Errorf("%s: rfCopy data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP6),printnwp(rfCopy))
	}

	//Check with all params set, and rf suffix

	emsgBase = "Init() with ALL params present"
	clearNWP(&nwp)
	nwp.SyslogSpec = expNWP4a.Oem.Syslog.SyslogServers[0] + "," +
	                 expNWP4a.Oem.Syslog.SyslogServers[1] + ":" +
	                 fmt.Sprintf("%d",expNWP4a.Oem.Syslog.Port)
	nwp.NTPSpec = expNWP4a.NTP.NTPServers[0] + "," +
	              expNWP4a.NTP.NTPServers[1] + ":" +
	              fmt.Sprintf("%d",expNWP4a.NTP.Port)
	nwp.SSHKey = "aabbccdd"
	nwp.SSHConsoleKey = "eeffgghh"
	rfSuffix = "/a/b/c"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Error(emsgBase,"error:",err)
	}
	if (!reflect.DeepEqual(nwProtoInfo,expNWP4a)) {
		t.Errorf("%s: data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP4a),printnwp(nwProtoInfo))
	}
	if (redfishNPSuffix != rfSuffix) {
		t.Errorf("Redfish suffix not applied to library!\n")
	}
	verifyNWP(t,nwProtoInfo)
	rfCopy = CopyRFNetworkProtocol(&nwProtoInfo)
	if (!reflect.DeepEqual(rfCopy,expNWP4a)) {
		t.Errorf("%s: rfCopy data mismatch, exp:\n%s\ngot:\n%s\n",
			emsgBase,printnwp(expNWP4a),printnwp(rfCopy))
	}

	//Check bad formats

	clearNWP(&nwp)
	nwp.SyslogSpec = "1.1.1.1:xxx"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err == nil) {
		t.Error("syslogTarg with bad format:",nwp.SyslogSpec,"didn't fail!")
	}
	if ((nwProtoInfo.Oem != nil) && (nwProtoInfo.Oem.Syslog != nil) && 
	    (len(nwProtoInfo.Oem.Syslog.SyslogServers) > 0)) {
		t.Error("Bad syslogTarg incorrectly resulted in valid SyslogServer.")
	}
	if ((nwProtoInfo.NTP != nil) && (len(nwProtoInfo.NTP.NTPServers) == 0)) {
		t.Error("Bad syslogTarg resulted in bad NTP server data (shouldn't have).")
	}

	clearNWP(&nwp)
	nwp.SyslogSpec = "1.1.1.1:123"
	nwp.NTPSpec = "1.1.1.1:xxx"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err == nil) {
		t.Error("ntpTarg with bad format:",nwp.NTPSpec,"didn't fail!")
	}
	if (len(nwProtoInfo.Oem.Syslog.SyslogServers) == 0) {
		t.Error("Bad ntpTarg resulted in bad Syslog server data (shouldn't have).")
	}
	if ((nwProtoInfo.NTP != nil) && (len(nwProtoInfo.NTP.NTPServers)) > 0) {
		t.Error("Bad ntpTarg incorrectly resulted in valid SyslogServer.")
	}

	clearNWP(&nwp)
	nwp.SyslogSpec = "1.1.1.1"
	nwp.NTPSpec = "2.2.2.2:123"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err == nil) {
		t.Error("Bad syslogTarg (no port) didn't fail!")
	}

	clearNWP(&nwp)
	nwp.SyslogSpec = "1.1.1.1:123"
	nwp.NTPSpec = "2.2.2.2"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err == nil) {
		t.Error("Bad NTPTarg (no port) didn't fail!")
	}

	clearNWP(&nwp)
	nwp.SyslogSpec = "1.1.1.x:123"
	nwp.NTPSpec = "2.2.2.2:234"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err == nil) {
		t.Error("Bad syslogTarg (invalid IP) didn't fail!")
	}

	nwp.SyslogSpec = "1.1.1.1:123"
	nwp.NTPSpec = "2.2.2.x:234"
	nwProtoInfo,err = Init(nwp,rfSuffix)
	if (err == nil) {
		t.Error("Bad NTPTarg (invalid IP) didn't fail!")
	}
}

var expSvcName = "NWPTest"
var gotUA bool

func hasUserAgentHeader(r *http.Request) bool {
    if (len(r.Header) == 0) {
        return false
    }

    alist,ok := r.Header["User-Agent"]
    if (!ok) {
        return false
    }

	for _,hdr := range(alist) {
		if (hdr == expSvcName) {
			return true
		}
	}
    return true
}

func nwpPatchHandler(w http.ResponseWriter, req *http.Request) {
	var nwp RedfishNWProtocol

	gotUA = hasUserAgentHeader(req)

	if (stallG) {
		time.Sleep(18 * time.Second)
	}

	body,err := ioutil.ReadAll(req.Body)
	if (err != nil) {
		__global_t.Error("ERROR dnsPutHandler() reading request body:",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body,&nwp)
	if (err != nil) {
		__global_t.Error("ERROR nwpPatchHandler() unmarshaling json body:",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Set the 'valid' fields to true in both instances -- it won't unmarshal, 
	//by design.

	nwProtoGlobal.valid = true
	nwp.valid = true

	if (!reflect.DeepEqual(nwProtoGlobal,nwp)) {
		__global_t.Errorf("ERROR, nwpPatchHandler() data miscompare.  exp:\n%s\ngot:\n%s\n",
			printnwp(nwProtoGlobal),printnwp(nwp))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCodeG)
}

func Test_setXnameNWPInfo(t *testing.T) {
	var nwpServer *httptest.Server
	var rfSuffix string
	var nwp NWPData
	var err error

	__global_t = t
	statusCodeG = http.StatusOK

	nwp.SyslogSpec = "150.151.152.153,160.161.162.163:999"
	nwp.NTPSpec = "170.171.172.173,180.181.182.183:888"
	nwp.SSHKey = "aabbccdd"
	nwp.SSHConsoleKey = "wwxxyyzz"
	rfSuffix = "__test__"	//force "test mode"

	nwProtoGlobal,err = InitInstance(nwp,rfSuffix,expSvcName)
	if (err != nil) {
		t.Error("ERROR creating network protocol data:",err)
	}
	nwpServer = httptest.NewServer(http.HandlerFunc(nwpPatchHandler))

	t.Log("Testing SetXNameNWPInfo happy path.")

	//"overload" the 'address' parameter of this func to pass the test URL

	err = SetXNameNWPInfo(nwProtoGlobal,nwpServer.URL,"user","pass")
	if (err != nil) {
		t.Error("ERROR in setXnameNWPInfo():",err)
	}
	if (!gotUA) {
		t.Error("ERROR, didn't get User-Agent header in request.")
	}
	nwpServer.Close()

	//Same test, using default service name

	serviceName = ""
	expSvcName,_ = os.Hostname()
	nwProtoGlobal,err = Init(nwp,rfSuffix)
	if (err != nil) {
		t.Error("ERROR creating network protocol data:",err)
	}
	nwpServer = httptest.NewServer(http.HandlerFunc(nwpPatchHandler))

	t.Log("Testing SetXNameNWPInfo happy path.")

	//"overload" the 'address' parameter of this func to pass the test URL

	err = SetXNameNWPInfo(nwProtoGlobal,nwpServer.URL,"user","pass")
	if (err != nil) {
		t.Error("ERROR in setXnameNWPInfo():",err)
	}
	if (!gotUA) {
		t.Error("ERROR, didn't get User-Agent header in request.")
	}
	nwpServer.Close()


	// Test with no valid bit set

	nwProtoGlobal.valid = false
	t.Log("Testing SetXNameNWPInfo without valid bit set.")
	nwpServer = httptest.NewServer(http.HandlerFunc(nwpPatchHandler))
	err = SetXNameNWPInfo(nwProtoGlobal,nwpServer.URL,"user","pass")
	nwpServer.Close()
	if (err == nil) {
		t.Error("SetXNameNWPInfo without valid should have failed!")
	}
	nwProtoGlobal.valid = true

	// Test with bad status code

	statusCodeG = http.StatusBadRequest
	t.Log("Testing SetXNameNWPInfo forcing error handler return.")
	nwpServer = httptest.NewServer(http.HandlerFunc(nwpPatchHandler))
	err = SetXNameNWPInfo(nwProtoGlobal,nwpServer.URL,"user","pass")
	nwpServer.Close()
	statusCodeG = http.StatusOK
	if (err == nil) {
		t.Error("SetXNameNWPInfo should have returned an HTTP error.")
	}

	//Test timeout by forcing a stall in the handler

	t.Log("Testing SetXNameNWPInfo timeout path (17 seconds).")
	nwpServer = httptest.NewServer(http.HandlerFunc(nwpPatchHandler))

	start := time.Now()
	stallG = true
	SetXNameNWPInfo(nwProtoGlobal,nwpServer.URL,"user","pass")
	nwpServer.Close()
	end := time.Now().Sub(start)
	elapsed := int(float64(end) / float64(time.Second))
	if (elapsed < 17) {
		t.Errorf("ERROR, timeout of RF operation too short, was %d, should be >= 17.",
			elapsed)
	}
	stallG = false
}

