package main

import (
    "bufio"
    "os"
    "regexp"
    "encoding/json"
    "strings"
    "flag"
    "log"
)

var Logger *log.Logger
var vendorOuis []vendorMap

type vendorMap struct {
    Oui string `json:"oui"`
    Vendor string `json:"vendor"`
}

var Vendors = func() map[string]string {
    return map[string]string {
        "Dell": "Dell Inc.",
        "Nokia": "Nokia Corporation" }
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

func getOuisByVendor(vendor string, filePath string) {
    re := regexp.MustCompile(vendor)
    re2 := regexp.MustCompile("[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}")

    fileHandle, err := os.Open(filePath)
    check(err)
    defer fileHandle.Close()
    fileScanner := bufio.NewScanner(fileHandle)

    for fileScanner.Scan() {
        text := fileScanner.Text()
        if (re.MatchString(text)) {
            Logger.Println("Matched vendor text: ", text)
            match := re2.FindString(text)
            if(match != "") {
                vendorOuis = append(vendorOuis, vendorMap{strings.Replace(match, "-", ":", 2), vendor})
            }
    	}
    }
}

func generateVendorMap(jsondata string) {
    fileHandle, err := os.Create("vendormap.json")
    check(err)
    defer fileHandle.Close()

    w := bufio.NewWriter(fileHandle)
    _, err = w.WriteString(jsondata)
    check(err)

    w.Flush();
}

func main() {
    filePtr := flag.String("file", "oui.txt", "a path to oui.txt file")
    logPtr := flag.String("log", "create-mapping.log", "a file path for logging")

    flag.Parse()
    Init(*logPtr)

    Logger.Println("Vendor Dict: ", Vendors())
    for _,v := range Vendors() {
        getOuisByVendor(v, *filePtr)
        Logger.Printf("Vendor ouis: %+v", vendorOuis)
    }

    jsondata, err := json.Marshal(vendorOuis)
    check(err)
    Logger.Println(string(jsondata))

    generateVendorMap(string(jsondata))
}

