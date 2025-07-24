package auth

import (
	"crocus-calculator/internal/api/middleware"
	"crocus-calculator/internal/dto"
	ht "crocus-calculator/internal/lib/http"
	"net/http"
	"strconv"
)

func (c AuthController) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Info("Handling user registration")

		req := r.Context().Value(middleware.DataKey).(dto.RegisterRequest)

		err := c.authService.Register(req)
		if err != nil {
			ht.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		ht.SendMessage(w, r, http.StatusOK, "User registered successfully")
		c.logger.Info("User registered successfully", "login", req.Login)
	}
}

func (c AuthController) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infof("Handling user login")

		req := r.Context().Value(middleware.DataKey).(dto.LoginRequest)

		tokenPair, err := c.authService.Login(req)
		if err != nil {
			ht.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		access, refresh := tokenPair.Access, tokenPair.Refresh

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refresh.Token,
			HttpOnly: true,
			Secure:   true, // Только для HTTPS
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   int(refresh.TTL.Seconds()),
		})

		response := map[string]string{
			"access_token": access.Token,
			"token_type":   "Bearer",
			"expires_in":   strconv.Itoa(int(access.TTL.Seconds())),
		}

		ht.SendJSON(w, r, http.StatusOK, response)
	}
}

func (c AuthController) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infof("Handling refresh")

		userID := r.Context().Value(middleware.UserIDKey).(uint)

		tokenPair, err := c.authService.Refresh(userID)
		if err != nil {
			ht.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		access, refresh := tokenPair.Access, tokenPair.Refresh

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refresh.Token,
			HttpOnly: true,
			Secure:   true, // Только для HTTPS
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   int(refresh.TTL.Seconds()),
		})

		response := map[string]string{
			"access_token": access.Token,
			"token_type":   "Bearer",
			"expires_in":   strconv.Itoa(int(access.TTL.Seconds())),
		}

		ht.SendJSON(w, r, http.StatusOK, response)
	}
}
