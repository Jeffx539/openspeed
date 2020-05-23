package main


import (
	"github.com/gorilla/websocket"
	"github.com/oschwald/geoip2-golang"
	"net"
	"net/http"
	"math/rand"
	"log"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

var testMemory []byte
var geoIPDB *geoip2.Reader

const(
	memoryMax = 1024*1024*512


)

type InfoResponse struct {
	RemoteAddress    string `json:"remoteAddress"`
	ASN uint `json:"autonomousSystemNumber"`
	ASOrg string `json:"autonomousSystemOrganisation"`
}
  


var upgrader = websocket.Upgrader{} // use default options

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

		if (string(message) == "ping"){
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





func handleInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err,port)
		http.Error(w, "Failed to Parse IP", http.StatusInternalServerError)
	}

	record, err := geoIPDB.ASN(net.ParseIP(ip))
	if err != nil {
		log.Println(err)
		return
	}

	json, err := json.Marshal(InfoResponse{RemoteAddress:ip, ASN:record.AutonomousSystemNumber, ASOrg: record.AutonomousSystemOrganization})
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


	if(len(chunkSizeStr) != 1 ){
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request - size not set"))
		return
	} 
  
	chunkSize, err := strconv.ParseInt(chunkSizeStr[0],10, 64)

	if(err != nil){
		http.Error(w, "Failed to Chunk Size", http.StatusInternalServerError)
	}
	

	if(chunkSize < 0 || chunkSize > memoryMax){
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid chunk size must be between 0 and memoryMax"))
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(chunkSize, 10))
	w.Write(testMemory[0:chunkSize])
}



func handleUpload(w http.ResponseWriter, r *http.Request){
	enableCors(&w)

	if r.Method == "GET"{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Upload route, you should probably post to this instead"))
	}


	buf, err := ioutil.ReadAll(r.Body)
    if err != nil {
		log.Printf("upload request failed ",err)
	}

	
	log.Println(len(buf))
	w.Write([]byte("ok"))

}




func allocTestMemory(){
	log.Println("Allocating Memory for tests")
	testMemory = make([]byte, memoryMax)
	rand.Read(testMemory)
	log.Println("Memory Allocated")
}



func main() {

	
	geo, err := geoip2.Open("GeoLite2-ASN.mmdb")
	

	if err != nil {
		log.Fatal("Failed to Open GeoIP Database")
	}
	
	geoIPDB = geo
	allocTestMemory()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/chunk", handleChunk)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/info", handleInfo)

	
	log.Println("We're Ready to go. Now Listening")
	http.ListenAndServe("0.0.0.0:4000",nil)
}
