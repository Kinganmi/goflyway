package config

import (
	"../lru"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

var (
	G_Key      = flag.String("k", "0123456789", "key, max length = 10")
	G_KeyBytes = []byte(*G_Key)

	G_Username = flag.String("u", "username", "proxy username")
	G_Password = flag.String("p", "password", "proxy password")
	G_Upstream = flag.String("up", "", "upstream server address (e.g. 127.0.0.1:8100)")
	G_Local    = flag.String("l", ":8100", "local listening")

	// G_WriteThreshold = flag.Int("wt", 0, "speed threshold in kb/s, 0 for unlimited")

	G_Debug      = flag.Bool("debug", false, "debug mode")
	G_UnsafeHttp = flag.Bool("disable-sh", false, "do not encrypt http content")
	G_NoPA       = flag.Bool("disable-pa", false, "disable proxy authentication")
	G_ProxyAll   = flag.Bool("all", false, "proxy Chinese websites")

	G_SuppressSocketReadWriteError = flag.Bool("ssrwe", false, "suppress socket read/write error")
	G_DNSCacheEntries              = flag.Int("dns-cache", 1024, "DNS cache size")

	G_Cache *lru.Cache
)

func setString(k *string, v interface{}) {
	switch v.(type) {
	case string:
		if v.(string) != "" {
			*k = v.(string)
		}
	}
}

func setBool(k *bool, v interface{}) {
	switch v.(type) {
	case bool:
		*k = v.(bool)
	}
}

func setInt(k *int, v interface{}) {
	switch v.(type) {
	case int:
		*k = v.(int)
	}
}

func LoadConfig(path string) {
	if path == "" {
		return
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(buf, &m)
	if err != nil {
		log.Println(err)
		return
	}

	setString(G_Key, m["key"])
	setString(G_Username, m["username"])
	setString(G_Password, m["password"])
	setString(G_Local, m["listen"])
	setString(G_Upstream, m["upstream"])

	setBool(G_UnsafeHttp, m["unsafehttp"])
	setBool(G_SuppressSocketReadWriteError, m["ssrwe"])
	setBool(G_NoPA, m["disablepa"])
	setBool(G_ProxyAll, m["proxyall"])

	setInt(G_DNSCacheEntries, m["dnscache"])

	if *G_Key == "." {
		// use username/password combination as the key
		*G_Key = *G_Username + *G_Password
	}

	G_KeyBytes = []byte(*G_Key)
	if len(G_KeyBytes) < 10 {
		G_KeyBytes = append(G_KeyBytes, make([]byte, 10-len(G_KeyBytes))...)
	}
}