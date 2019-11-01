package main_test

import (
	"bytes"
	"encoding/json"
	. "github.com/benjamintf1/unmarshalledmatchers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"instatasks/config"
	"instatasks/database"
	"instatasks/database/migrations"
	"instatasks/models"
	"instatasks/redis_storage"
	"instatasks/router"
	. "instatasks/test_helpers"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
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
	// autocleaner = DatabaseAutocleaner(db)
	migrations.Migrate()
	models.Init()
	userAgent := &models.UserAgent{Name: "user_agent_with_default_settings"}
	userAgent.Create()
})

var _ = AfterSuite(func() {
	// autocleaner()
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
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusOK))
			})

			It("Should return pong JSON", func() {
				req, err := http.NewRequest("GET", "/ping", nil)
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)

				Ω(err).ShouldNot(HaveOccurred())
				Ω(w.Body).Should(MatchUnorderedJSON(`{"message":"pong"}`))
			})
		})
	})
	/////////////////////////////////////////////////////////////
	Describe("Accaunt (POST /accaunt) route", func() {
		reqBody := []byte(`{
			"instagramid": 666,
			"deviceid":    "device1"
		}`)

		Context("None banned User", func() {
			wrongReqBody := []byte(`{
				"instagramid": 666
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
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)
				///*******
				LogHttpData(req, w)

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
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body).Should(MatchUnorderedJSON(expected_response))
			})

			It("must be present deviceid", func() {

				req, _ := http.NewRequest("POST", "/accaunt", bytes.NewBuffer(wrongReqBody))
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
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
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)
				///*******
				LogHttpData(req, w)

				Ω(w.Code).Should(Equal(http.StatusForbidden))
			})

			It("bunned Device must force ForbiddenError response status", func() {
				bannedDevice := models.BannedDevice{Deviceid: "device2"}
				db.Save(&bannedDevice)

				bannedDeviceReqBody := []byte(`{
					"instagramid": 667,
					"deviceid":    "device2"
				}`)

				req, _ := http.NewRequest("POST", "/accaunt", bytes.NewBuffer(bannedDeviceReqBody))
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)
				///*******
				LogHttpData(req, w)

				Ω(w.Code).Should(Equal(http.StatusForbidden))
			})
		})
	})

	/////////////////////////////////////////////////////////////
	Describe("User Agent settings (POST /setting) route", func() {

		It("must return User Agent settings", func() {

			req, _ := http.NewRequest("POST", "/setting", nil)
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`)
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			Ω(w.Code).Should(Equal(http.StatusOK))
			Ω(w.Body).Should(ContainUnorderedJSON(`{"activitylimit": 0, "like": true }`))
		})
	})

	/////////////////////////////////////////////////////////////
	Describe("Create Task (POST /newwork) route", func() {

		It("User must create new Task if User have enough Coins", func() {
			reqBody := []byte(`{ 
				"instagramid": 777,
				"count": 10,
				"type": "like",
				"mediaid":"mediaid1",
				"photourl":"url/blabla", 
				"instagramusername":"url/blabla" 
			}`)
			user := models.User{Instagramid: 777, Coins: 20}
			user.Save()

			req, _ := http.NewRequest("POST", "/newwork", bytes.NewBuffer(reqBody))
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`) // like default price 1
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			Ω(db.First(&models.Task{Instagramid: 777}).RecordNotFound()).Should(BeFalse())
			Ω(w.Code).Should(Equal(http.StatusOK))
			Ω(w.Body).Should(MatchUnorderedJSON(`{"coins": 10}`))
		})

		It("User must get error code 406 NotAcceptable if User does not have enough coins", func() {
			reqBody := []byte(`{
				"instagramid": 888,
				"count": 20,
				"type": "like",
				"mediaid":"mediaid1",
				"photourl":"url/blabla", 
				"instagramusername":"url/blabla" 
			}`)
			user := models.User{Instagramid: 888, Coins: 10}
			user.Save()

			req, _ := http.NewRequest("POST", "/newwork", bytes.NewBuffer(reqBody))
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`) // like default price 1
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			Ω(db.Where(&models.Task{Instagramid: 888}).First(&models.Task{}).RecordNotFound()).Should(BeTrue())
			Ω(w.Code).Should(Equal(http.StatusNotAcceptable))
			Ω(w.Body).Should(MatchUnorderedJSON(`{"error": "Not Enough Coins"}`))
		})
	})

	/////////////////////////////////////////////////////////////
	Describe("Get Tasks history (POST /history) route", func() {

		It("must return last 10 User tasks (soft deleted too)", func() {
			task := models.Task{
				Instagramid:       777,
				Type:              "like",
				Count:             20,
				Photourl:          "url/blabla",
				Instagramusername: "url/blabla",
				Mediaid:           "mediaid2",
			}

			for i := 0; i < 4; i++ {
				temp := task
				db.Create(&temp)
			}
			for i := 0; i < 4; i++ {
				temp := task
				db.Create(&temp)
				db.Delete(&temp)
			}
			for i := 0; i < 6; i++ {
				temp := task
				db.Create(&temp)
			}

			req, _ := http.NewRequest("POST", "/history", strings.NewReader(`{"instagramid": 777}`))
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`)
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			resp := []struct {
				Taskid string `json:"taskid"`
			}{}
			json.Unmarshal([]byte(w.Body.String()), &resp)

			Ω(len(resp)).Should(Equal(10))
			first_id, _ := strconv.Atoi(resp[0].Taskid)
			last_id, _ := strconv.Atoi(resp[len(resp)-1].Taskid)

			Ω(first_id > last_id).Should(BeTrue())
			Ω(w.Code).Should(Equal(http.StatusOK))

		})
	})

	/////////////////////////////////////////////////////////////
	Describe("Get Tasks (POST /gettasks) route", func() {
		Context("Return last 5 tasks by type and Mediad/Instagramusername User done/cencel verification ", func() {

			It("must show 'like' Tasks", func() {

				task := models.Task{
					Instagramid:       888,
					Type:              "like",
					Count:             20,
					Photourl:          "url/blabla",
					Instagramusername: "url/blabla",
					Mediaid:           "mediaid2",
				}

				for i := 0; i < 2; i++ {
					temp := task
					db.Create(&temp)
				}
				for i := 0; i < 2; i++ {
					temp := task
					temp.Type = "follow"
					db.Create(&temp)
				}
				for i := 0; i < 2; i++ {
					temp := task
					db.Create(&temp)
					db.Delete(&temp)
				}

				tasks := []models.Task{}
				db.Order("id desc").Limit(6).Find(&tasks)

				// doned tasks
				tasks[0].Done()
				user_mediaid := models.UserMediaid{Instagramid: 888, Mediaid: tasks[0].Mediaid}
				db.Create(&user_mediaid)

				tasks[1].Done()
				user_mediaid = models.UserMediaid{Instagramid: 889, Mediaid: tasks[1].Mediaid}
				db.Create(&user_mediaid)

				req, _ := http.NewRequest("POST", "/gettasks", strings.NewReader(`{"instagramid": 777, "type":"like"}`))
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)
				///*******
				LogHttpData(req, w)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body).Should(ContainUnorderedJSON(`[{"type": "like"}]`))
			})

			It("must show 'all' Tasks ", func() {

				req, _ := http.NewRequest("POST", "/gettasks", strings.NewReader(`{"instagramid": 888, "type":"all"}`))
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("User-Agent", `user_agent_with_default_settings`)
				r.ServeHTTP(w, req)
				///*******
				LogHttpData(req, w)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body).ShouldNot(ContainUnorderedJSON(`[{"mediaid": "mediaid1"}]`))
			})
		})
	})

	/////////////////////////////////////////////////////////////
	Describe("Done Task (POST /done) route", func() {

		It("must return ResetContent (205) status if success CCANCEL status DONE", func() {
			cancel_req := `{
				"instagramid": 888,
				"taskid": "19",
				"status": "cancel"
			}`

			req, _ := http.NewRequest("POST", "/done", strings.NewReader(cancel_req))
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`)
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			Ω(w.Code).Should(Equal(http.StatusResetContent))
		})

		It("must deposit and return coins if success OK status DONE", func() {
			ok_req := `{
				"instagramid": 888,
				"taskid": "19", 
				"status": "ok"
			}`

			req, _ := http.NewRequest("POST", "/done", strings.NewReader(ok_req))
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`)
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			Ω(w.Code).Should(Equal(http.StatusOK))
			/// task 19 is follow, price 5, must return 15
			Ω(w.Body).Should(ContainUnorderedJSON(`{"coins": 15}`))
		})
	})

	/////////////////////////////////////////////////////////////
	Describe("Rateus (POST /rateus) route", func() {

		It("must return ballance if task is not already done", func() {

			req, _ := http.NewRequest("POST", "/rateus", strings.NewReader(`{"instagramid": 888}`))
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`)
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			Ω(w.Code).Should(Equal(http.StatusOK))
			/// default price 20, previous balance 15, must return 35
			Ω(w.Body).Should(ContainUnorderedJSON(`{"coins": 35}`))
		})

		It("must return NotAcceptable (406) if task already done", func() {

			req, _ := http.NewRequest("POST", "/rateus", strings.NewReader(`{"instagramid": 888}`))
			req.Header.Add("Content-Type", `application/json`)
			req.Header.Add("User-Agent", `user_agent_with_default_settings`)
			r.ServeHTTP(w, req)
			///*******
			LogHttpData(req, w)

			Ω(w.Code).Should(Equal(http.StatusNotAcceptable))
		})
	})
})

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
				reqBody := []byte(`{
					"name": "useragent1",
					"activitylimit": 1
				}`)

				req, _ := http.NewRequest("POST", "/admin/useragent", bytes.NewBuffer(reqBody))
				req.Header.Add("Content-Type", `application/json`)
				req.Header.Add("Authorization", AuthorizationHeader(validSuperadminUsername, validSuperadminPassword))
				r.ServeHTTP(w, req)
				///*******
				LogHttpData(req, w)

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
				///*******
				LogHttpData(req, w)

				Ω(w.Code).Should(Equal(http.StatusOK))
				Ω(w.Body.String()).Should(ContainSubstring("BEGIN RSA PUBLIC KEY"))
			})
		})
	})
})
