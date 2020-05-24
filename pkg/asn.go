package pkg

import(
	"net"
	"compress/gzip"
	"os"
	"bufio"
	"strings"
	"sort"
	"strconv"
)


// ASNInfo docs 
type ASNInfo struct {
	
	startAddress net.IP
	endAddress net.IP
	ASNumber uint64
	CountryCode string
	Description string

}


type ASNDatabse struct {
	ipRanges []ASNInfo
}



func LoadASNInfoFile(path string) (*ASNDatabse, error){

	file, err := os.Open(path)
	if err != nil {
		return  nil, err
	}

	zr, err := gzip.NewReader(file)
	if err != nil {
		return  nil, err
	}
	
	scanner := bufio.NewScanner(zr)
	asnDB := ASNDatabse{}

	for scanner.Scan(){
		splitString := strings.Split(scanner.Text(), "\t")

		asn, err := strconv.ParseUint(splitString[2], 10, 64)

		if err != nil {
			return  nil, err
		}


		ds := ASNInfo{
			startAddress: net.ParseIP(splitString[0]),
			endAddress: net.ParseIP(splitString[1]),
			ASNumber: asn,
			CountryCode: splitString[3],
			Description: splitString[4],
		}
	
		asnDB.ipRanges = append(asnDB.ipRanges, ds)

	}

	return &asnDB,nil

}


func (db *ASNDatabse) IPToASN(ipAddress net.IP) (*ASNInfo){

	searchResult := sort.Search(len(db.ipRanges), func(i int) bool {
		return string(ipAddress.To16()) < string(db.ipRanges[i].startAddress.To16())
	})

	if searchResult-1 < 0 {
		return nil
	}

	if string(db.ipRanges[searchResult-1].endAddress.To16()) < string(ipAddress.To16()){
		return nil
	}

	return &db.ipRanges[searchResult-1]

}
