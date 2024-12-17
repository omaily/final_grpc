package delivery

import (
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler POST REGISTER"))
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler POST LOGIN"))
}
func Balance(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler GET BALANCE"))
}
func Deposit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler POST DEPOSiT"))
}
func Withdraw(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler POST WITHDRAW"))
}
func Rates(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler GET RATES"))
}
func Exchange(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler POST EXCHANGE"))
}
