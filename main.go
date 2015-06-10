package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ryokdy/go-kintone"
	"github.com/howeyc/gopass"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
)

type Configure struct {
	login string
	password string
	basicAuthUser string
	basicAuthPassword string
	apiToken string
	domain string
	basic string
	format string
	query string
	appId uint64
	fields []string
	filePath string
	deleteAll bool
	encoding string
}

var config Configure

const ROW_LIMIT = 100

type Column struct {
	Code        string
	Type        string
	IsSubField  bool
	Table       string
}

func getEncoding() encoding.Encoding {
	switch config.encoding {
	case "utf-16":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case "utf-16be-with-signature":
		return unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	case "utf-16le-with-signature":
		return unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	case "euc-jp":
		return japanese.EUCJP
	case "sjis":
		return japanese.ShiftJIS
	default:
		return nil
	}
}

func main() {
	var colNames string

	flag.StringVar(&config.login, "u", "", "Login name")
	flag.StringVar(&config.password, "p", "", "Password")
	flag.StringVar(&config.basicAuthUser, "U", "", "Basic authentication user name")
	flag.StringVar(&config.basicAuthPassword, "P", "", "Basic authentication password")
	flag.StringVar(&config.domain, "d", "", "Domain name")
	flag.StringVar(&config.apiToken, "t", "", "API token")
	flag.Uint64Var(&config.appId, "a", 0, "App ID")
	flag.StringVar(&config.format, "o", "csv", "Output format: 'json' or 'csv'(default)")
	flag.StringVar(&config.query, "q", "", "Query string")
	flag.StringVar(&colNames, "c", "", "Field names (comma separated)")
	flag.StringVar(&config.filePath, "f", "", "Input file path")
	flag.BoolVar(&config.deleteAll, "D", false, "Delete all records before insert")
	flag.StringVar(&config.encoding, "e", "utf-8", "Character encoding: 'utf-8'(default), 'utf-16', 'utf-16be-with-signature', 'utf-16le-with-signature, 'sjis' or 'euc-jp'")

	flag.Parse()

	if config.appId == 0 || (config.apiToken == "" && (config.domain == "" || config.login == "")) {
		flag.PrintDefaults()
		return
	}

	if !strings.Contains(config.domain, ".") {
		config.domain += ".cybozu.com"
	}

	if colNames != "" {
		config.fields = strings.Split(colNames, ",")
	}


	var app *kintone.App

	if config.basicAuthUser != "" && config.basicAuthPassword == "" {
		fmt.Printf("Basic authentication password: ")
		config.basicAuthPassword = string(gopass.GetPasswd())
	}

	if config.apiToken == "" {
		if config.password == "" {
			fmt.Printf("Password: ")
			config.password = string(gopass.GetPasswd())
		}

		app = &kintone.App{
			Domain:	  config.domain,
			User:	  config.login,
			Password: config.password,
			AppId:	  config.appId,
		}
	} else {
		app = &kintone.App{
			Domain:	  config.domain,
			ApiToken: config.apiToken,
			AppId:	  config.appId,
		}
	}

	if config.basicAuthUser != "" {
		app.SetBasicAuth(config.basicAuthUser, config.basicAuthPassword)
	}

	var err error
	if config.filePath == "" {
		if config.format == "json" {
			err = writeJson(app)
		} else {
			err = writeCsv(app)
		}
	} else {
		err = readCsv(app, config.filePath)
	}
	if err != nil {
		log.Fatal(err)
	}
}
