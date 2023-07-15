package Middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func handleAuth(w http.ResponseWriter, r *http.Request) {
	// Get the SSO provider from the URL path
	// provider := chi.URLParam(r, "provider")

	// Initiate the authentication flow for the specified provider
	gothic.BeginAuthHandler(w, r)
}

func (c *Middlewares) handleCallback(store *sessions.CookieStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session, err := store.Get(r, "session-name")
		if err != nil {
			// Handle error
			fmt.Println(err)

			fmt.Println("getting session name")
			http.Redirect(w, r, "/err", http.StatusFound)
			return
		}

		// Get the user object from gothic.CompleteUserAuth
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			// Handle error

			fmt.Println("completing auth")
			http.Redirect(w, r, "http://localhost:8080/err", http.StatusFound)
			return
		}

		_, err = c.DB.FindUserByEmail(context.Background(), user.Email)
		if err != nil {

			http.Redirect(w, r, "http://localhost:8080/err", http.StatusFound)
			return
		}
		// Store user information in the session
		session.Values["userID"] = user.UserID
		// Save the session
		err = session.Save(r, w)
		if err != nil {
			// Handle error

			fmt.Println("saving session")
			http.Redirect(w, r, "/err", http.StatusFound)
			return
		}
		print("success Redirect")
		http.Redirect(w, r, "http://localhost:8080/health", http.StatusFound)
	}

}

func SSOAuthMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session-name")
			if err != nil {
				// Handle error
			}
			// Check if the user is authenticated using SSO
			if !isAuthenticatedSSO(r, session) {
				http.Redirect(w, r, "/v1/auth/provider?provider=google", http.StatusFound)
				return
			}

			// User is authenticated, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func isAuthenticatedSSO(r *http.Request, session *sessions.Session) bool {
	// Check if the user is authenticated using SSO
	_, ok := session.Values["userID"].(string)
	if !ok {
		return false
		// User is not authenticated, handle accordingly
	}

	return true
	// Check if the user object contains necessary authentication data
	// ...

}
func initializeGoogleProvider() {

	GOOGLE_CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET := os.Getenv("GOOGLE_CLIENT_SECRET")
	REDIRECT_URL := os.Getenv("REDIRECT_URL")
	if GOOGLE_CLIENT_ID == "" || GOOGLE_CLIENT_SECRET == "" || REDIRECT_URL == "" {
		log.Fatal("Invalid requirement for googleauth")
	}
	goth.UseProviders(
		google.New(GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, REDIRECT_URL+"google", "profile", "email"),
	)
}
func (c *Middlewares) Init() (r *chi.Mux, store *sessions.CookieStore) {
	initializeGoogleProvider()
	// Configure Google SSO provider

	// // Configure Facebook SSO provider
	// goth.UseProviders(
	// 	facebook.New("YOUR_FACEBOOK_CLIENT_ID", "YOUR_FACEBOOK_CLIENT_SECRET", "YOUR_REDIRECT_URL", "email"),
	// )

	// goth.UseProviders(
	// 	linkedin.New("YOUR_FACEBOOK_CLIENT_ID", "YOUR_FACEBOOK_CLIENT_SECRET", "YOUR_REDIRECT_URL", "email"),
	// )
	r = chi.NewRouter()
	SECRET_KEY := os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		log.Fatal("SECRET_KEY is missing from env")
	}

	store = sessions.NewCookieStore([]byte(SECRET_KEY))
	r.Get("/{provider}", handleAuth)

	r.Get("/callback/{provider}", c.handleCallback(store))
	return r, store

	// Add more SSO providers as needed
}
