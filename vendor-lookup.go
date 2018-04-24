package main

import (
    "bufio"
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
    "strings"
    "flag"
    "log"
)

var Logger *log.Logger
var Ouis map[string]string

type vendorMap struct {
    Oui string `json:"oui"`
    Vendor string `json:"vendor"`
}

func Init(logfile string) {
    file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file", err)
    }

    Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func check(e error) {
    if e != nil {
        Logger.Panic("Error: ", e)
    }
}

func parseVendorMap(confpath string) {
    dat, err := ioutil.ReadFile(confpath)
    check(err)
    Logger.Println("File Content: ", string(dat))

    var results []vendorMap
    json.Unmarshal(dat, &results)
    Logger.Println("Parsed JSON Results: ", results)

    Ouis = make(map[string]string)
    for _, v := range results {
        Ouis[v.Oui] = v.Vendor
    }
}

func getVendor(macAddr string) string {
    vendoroui := (macAddr[0:8])
    return Ouis[strings.ToUpper(vendoroui)]
}

func getVendors(filePath string) {
    fileHandle, err := os.Open(filePath)
    check(err)
    defer fileHandle.Close()

    fileScanner := bufio.NewScanner(fileHandle)
    for fileScanner.Scan() {
        text := fileScanner.Text()
        mac := (strings.Split(text, "\""))[1]
        Logger.Println("MacAddress from file: ", mac)

        vendor := getVendor(mac)
        if(vendor == "") {
            Logger.Panic("Unknown MAC Address. No vendor found")
        }

        fmt.Printf("MACAddress:%v     Vendor: %v", mac, vendor)
        fmt.Println();
    }
}

func main() {
    macPtr := flag.String("mac", "", "a MAC address")
    filePtr := flag.String("file", "", "a file path to the MAC addresses")
    confPtr := flag.String("conf", "vendormap.json", "a config file with vendor oui mapping")
    logPtr := flag.String("log", "get-vendor.log", "a file path for logging")

    flag.Parse()
    Init(*logPtr)

    macAddr := *macPtr
    Logger.Println("MacAddress passed in: ", macAddr)

    filePath := *filePtr
    Logger.Println("File path passed in: ", filePath)

    if(filePath == "" && macAddr == "") {
        reader := bufio.NewReader(os.Stdin)
        fmt.Println("Please enter a MAC address: ")
        macAddr, _ = reader.ReadString('\n')
        Logger.Println("MacAddress inputted: ", macAddr)
    }

    parseVendorMap(*confPtr)

    if(macAddr != "") {
        vendor := getVendor(macAddr)
        if(vendor == "") {
            Logger.Panic("Unknown MAC Address. No vendor found")
        }

        fmt.Printf("Vendor: %v", vendor)
        fmt.Println();
    } else {
        getVendors(filePath)
    }
}