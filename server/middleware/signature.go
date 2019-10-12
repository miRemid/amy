package middleware

import (
	"log"
	"fmt"
	"io"
	"io/ioutil"
	"crypto/sha1"
	"crypto/hmac"	
	"net/http"	
)

// SignatureMiddleware CQHTTP消息验证中间
func SignatureMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		log.Println("Signature")
		sig := r.Header.Get("X-Signature")
		if sig == "" {
			return
		}
		sig = sig[len("sha1="):]
		mac := hmac.New(sha1.New, []byte("amy"))
		byteData, _ := ioutil.ReadAll(r.Body)

		io.WriteString(mac, string(byteData))
		res := fmt.Sprintf("%x", mac.Sum(nil))
		log.Printf("CQ HMAC: %s\nAmy HMAC: %s\n", sig, res)
		if res == sig {
			log.Println("It's CoolQ")
		}else{
			log.Println("Data From CoolQ Wrong")
			return
		}		
		h.ServeHTTP(w, r)
	})
}