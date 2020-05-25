package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/jeffx539/openspeed/pkg"
)

var testMemory []byte
var asnDB *pkg.ASNDatabse

const (
	memoryMax = 1024 * 1024 * 128
)

// InfoResponse struct contains the response to the /info route
type InfoResponse struct {
	RemoteAddress string `json:"remoteAddress"`
	ASN           uint64 `json:"autonomousSystemNumber"`
	ASOrg         string `json:"autonomousSystemOrganisation"`
	Country       string `json:"autonomousSystemCountry"`
}

var upgrader = websocket.Upgrader{}

func handlePing(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if string(message) == "ping" {
			err = c.WriteMessage(mt, []byte("pong"))
			if err != nil {
				log.Println("write:", err)
				break
			}
		} else {
			break
		}

	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getProxyAddress(r *http.Request) string {

	if os.Getenv("REVERSE_PROXY") == "" {
		return r.RemoteAddr
	}
	if r.Header.Get("X-Real-Ip") != "" {
		return r.Header.Get("X-Real-Ip") + ":1234"
	}

	if r.Header.Get("X-Forwarded-For") != "" {
		return r.Header.Get("X-Forwarded-For") + ":1234"
	}

	return r.RemoteAddr

}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	ip, port, err := net.SplitHostPort(getProxyAddress(r))

	if err != nil {
		log.Println(err, port)
		http.Error(w, "Failed to Parse IP", http.StatusInternalServerError)
	}

	asnInfo := asnDB.IPToASN(net.ParseIP(ip))

	if asnInfo == nil {
		asnInfo = &pkg.ASNInfo{ASNumber: 0, Description: "None", CountryCode: "XX"}
	}

	json, err := json.Marshal(InfoResponse{RemoteAddress: ip, ASN: asnInfo.ASNumber, ASOrg: asnInfo.Description, Country: asnInfo.CountryCode})
	if err != nil {
		http.Error(w, "Failed To Marshal JSON", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func handleChunk(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	chunkSizeStr := r.URL.Query()["size"]

	if len(chunkSizeStr) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request - size not set"))
		return
	}

	chunkSize, err := strconv.ParseInt(chunkSizeStr[0], 10, 64)

	if err != nil {
		http.Error(w, "Failed to Chunk Size", http.StatusInternalServerError)
	}

	if chunkSize < 0 || chunkSize > memoryMax {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid chunk size must be between 0 and memoryMax"))
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(chunkSize, 10))
	w.Write(testMemory[0:chunkSize])
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Upload route, you should probably post to this instead"))
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("upload request failed %s", err.Error())
	}

	log.Println(len(buf))
	w.Write([]byte("ok"))

}

func allocTestMemory() {
	log.Println("Allocating Memory for tests")
	testMemory = make([]byte, memoryMax)
	rand.Read(testMemory)
	log.Println("Memory Allocated")
}

func main() {

	asn, err := pkg.LoadASNInfoFile("./ip2asn-combined.tsv.gz")
	asnDB = asn

	if err != nil {
		log.Fatal(err)
	}

	allocTestMemory()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/chunk", handleChunk)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/info", handleInfo)

	log.Println("We're Ready to go. Now Listening")
	http.ListenAndServe("0.0.0.0:4000", nil)
}
