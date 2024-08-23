package handler

import "net/http"

// TODO: Implement a function that checks if the token is valid
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//r.ParseForm()
			//username := r.Form.Get("username")
			//token := r.Form.Get("token")
			//
			//if len(username) < 3 || !IsTokenValid(token) {
			//	http.Redirect(w, r, "/user/login", http.StatusOK)
			//}
			h(w, r)

		})

}
