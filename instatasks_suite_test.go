package main_test

import (
	"bytes"
	// "encoding/json"
	. "github.com/benjamintf1/unmarshalledmatchers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"instatasks/config"
	"instatasks/database"
	"instatasks/models"
	"instatasks/redis_storage"
	"instatasks/router"
	. "instatasks/test_helpers"
	"net/http"
	"net/http/httptest"
	"os"
	// "strconv"
	// "fmt"
	"testing"
)

var (
	w           *httptest.ResponseRecorder
	r           *gin.Engine
	db          *gorm.DB
	autocleaner func()
)

func TestInstatasks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Instatasks Suite")
}

var _ = BeforeSuite(func() {
	os.Setenv("APP_ENV", "test")
	config.InitConfig()
	r = router.SetupRouter()
	db = database.InitDB()
	autocleaner = DatabaseAutocleaner(db)
	// redis_storage.InitCache()
	models.Init()
})

var _ = AfterSuite(func() {
	autocleaner()
	db.Close()
})

var _ = Describe("Instatasks API", func() {

	BeforeEach(func() {
		w = httptest.NewRecorder()
	})
	//////////////////////////////////////////////////////////////
	Describe("Ping (GET /ping) route", func() {
		Context("When ping succesfully", func() {
			It("Should return Ok code", func() {
				req, _ := http.NewRequest("GET", "/ping", nil)
				req.Header.Add("Content-Type", `application/json`)
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusOK))
			})

			It("Should return pong JSON", func() {
				req, err := http.NewRequest("GET", "/ping", nil)
				req.Header.Add("Content-Type", `application/json`)
				r.ServeHTTP(w, req)

				Ω(err).ShouldNot(HaveOccurred())
				Ω(w.Body).Should(MatchUnorderedJSON(`{"message":"pong"}`))
			})
		})
	})
	/////////////////////////////////////////////////////////////
	Describe("Accaunt (POST /accaunt) route", func() {
		reqBody := []byte(`{ "data": {
			"instagramid": 666,
			"deviceid":    "device1"
		}
		}`)

		Context("None banned User", func() {
			wrongReqBody := []byte(`{ "data": {
				"instagramid": 666
			}
			}`)

			expected_response := []byte(`{
				"instagramid": 666,
				"coins":       0,
				"rateus":      true
			}`)

			It("must create User if not exist  and return default User data, must be cached", func() {
				var cachedUser models.CachedUser

				Ω(db.First(&models.User{Instagramid: 666}).RecordNotFound()).Should(BeTrue())

				req, _ := http.NewRequest("POST", "/accaunt", bytes.NewBuffer(reqBody))
				req.Header.Add("Content-Type", `application/json`)
				r.ServeHTTP(w, req)

				Ω(db.First(&models.User{Instagramid: 666}).RecordNotFound()).Should(BeFalse())
				cache := redis_storage.GetCacheCodec("User")
				cache.Get("666", &cachedUser)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body).Should(MatchUnorderedJSON(expected_response))
				Ω(cachedUser).Should(Equal(models.CachedUser{false, 0, true}))
			})

			It("if User exist: return User data", func() {

				Ω(db.First(&models.User{Instagramid: 666}).RecordNotFound()).Should(BeFalse())

				req, _ := http.NewRequest("POST", "/accaunt", bytes.NewBuffer(reqBody))
				req.Header.Add("Content-Type", `application/json`)
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body).Should(MatchUnorderedJSON(expected_response))
			})

			It("must be present deviceid", func() {

				req, _ := http.NewRequest("POST", "/accaunt", bytes.NewBuffer(wrongReqBody))
				req.Header.Add("Content-Type", `application/json`)
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusBadRequest))
			})
		})

		Context("Banned User/Device", func() {

			It("banned User must force ForbiddenError response status", func() {
				var bannedUser models.User
				db.First(&bannedUser)
				bannedUser.Banned = true
				bannedUser.Save()

				req, _ := http.NewRequest("POST", "/accaunt", bytes.NewBuffer(reqBody))
				req.Header.Add("Content-Type", `application/json`)
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusForbidden))
			})

			It("bunned Device must force ForbiddenError response status", func() {
				bannedDevice := models.BannedDevice{Deviceid: "device2"}
				db.Save(&bannedDevice)

				bannedDeviceReqBody := []byte(`{ "data": {
					"instagramid": 667,
					"deviceid":    "device2"
				}
				}`)

				req, _ := http.NewRequest("POST", "/accaunt", bytes.NewBuffer(bannedDeviceReqBody))
				req.Header.Add("Content-Type", `application/json`)
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusForbidden))
			})
		})
	})

	/////////////////////////////////////////////////////////////
	Describe("User Agent settings (POST /setting) route", func() {
		It("must return User Agent settings", func() {
			userAgent := &models.UserAgent{Name: "user_agent_with_default_settings"}
			userAgent.Save()

			req, _ := http.NewRequest("POST", "/setting", nil)
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`)
			r.ServeHTTP(w, req)

			Ω(w.Code).Should(Equal(http.StatusOK))
			Ω(w.Body).Should(ContainUnorderedJSON(`{"activitylimit": 0, "like": true }`))
		})

	})
})

/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////

var _ = Describe("Instatasks Admin API. Admin (/admin) protected route group", func() {

	BeforeEach(func() {
		w = httptest.NewRecorder()
	})

	/////////////////////////////////////////////////////////////
	Describe("UserAgent (/admin/useragent) protected route", func() {

		Context("Unauthorized Admin", func() {

			It("must 401 Unauthorized", func() {
				req, _ := http.NewRequest("POST", "/admin/useragent", nil)
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("Authorization", AuthorizationHeader("wrongname", "wrongpass"))
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusUnauthorized))
			})
		})

		Context("Authorized Admin", func() {

			It("must generate rsa keys, AES encript rsa keys, create UserAgent (with default value) in DB", func() {
				validSuperadminUsername := config.GetConfig().Server.Superadmin.Username
				validSuperadminPassword := config.GetConfig().Server.Superadmin.Password
				reqBody := []byte(`{ "data": {
					"name": "useragent1",
					"activitylimit": 1
				}
				}`)

				req, _ := http.NewRequest("POST", "/admin/useragent", bytes.NewBuffer(reqBody))
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("Authorization", AuthorizationHeader(validSuperadminUsername, validSuperadminPassword))
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body).Should(ContainUnorderedJSON(`{"name": "useragent1", "activitylimit": 1, "like": true }`))
			})

			It("must get RSA Public Key (GET /admin/useragent/pkey)", func() {
				validSuperadminUsername := config.GetConfig().Server.Superadmin.Username
				validSuperadminPassword := config.GetConfig().Server.Superadmin.Password
				reqBody := []byte(`{"name": "useragent1"}`)

				req, _ := http.NewRequest("GET", "/admin/useragent/pkey", bytes.NewBuffer(reqBody))
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("Authorization", AuthorizationHeader(validSuperadminUsername, validSuperadminPassword))
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body.String()).Should(ContainSubstring("BEGIN RSA PUBLIC KEY"))
			})
		})
	})
})
