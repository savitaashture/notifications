package strategies_test

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserStrategy", func() {
	var strategy strategies.UserStrategy
	var options postal.Options
	var tokenLoader *fakes.TokenLoader
	var userLoader *fakes.UserLoader
	var templatesLoader *fakes.TemplatesLoader
	var mailer *fakes.Mailer
	var clientID string
	var receiptsRepo *fakes.ReceiptsRepo
	var conn *fakes.DBConn
	var users map[string]uaa.User

	BeforeEach(func() {
		clientID = "mister-client"

		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}

		tokenClaims := map[string]interface{}{
			"client_id": "mister-client",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		tokenLoader = fakes.NewTokenLoader()
		tokenLoader.Token = fakes.BuildToken(tokenHeader, tokenClaims)

		users = map[string]uaa.User{
			"user-123": uaa.User{
				ID:     "user-123",
				Emails: []string{"user-123@example.com"},
			},
		}

		receiptsRepo = fakes.NewReceiptsRepo()
		mailer = fakes.NewMailer()
		userLoader = fakes.NewUserLoader()
		userLoader.Users = users
		templatesLoader = fakes.NewTemplatesLoader()
		strategy = strategies.NewUserStrategy(tokenLoader, userLoader, templatesLoader, mailer, receiptsRepo)
	})

	Describe("Dispatch", func() {
		BeforeEach(func() {
			options = postal.Options{
				KindID:            "forgot_password",
				KindDescription:   "Password reminder",
				SourceDescription: "Login system",
				Text:              "Please reset your password by clicking on this link...",
				HTML: postal.HTML{
					BodyContent: "<p>Please reset your password by clicking on this link...</p>",
				},
			}
		})

		It("records a receipt for the user", func() {
			_, err := strategy.Dispatch(clientID, "user-123", options, conn)
			if err != nil {
				panic(err)
			}

			Expect(receiptsRepo.CreateUserGUIDs).To(Equal([]string{"user-123"}))
			Expect(receiptsRepo.ClientID).To(Equal(clientID))
			Expect(receiptsRepo.KindID).To(Equal(options.KindID))
		})

		It("calls mailer.Deliver with the correct arguments for a user", func() {
			templates := postal.Templates{
				Subject: "default-missing-subject",
				Text:    "default-space-text",
				HTML:    "default-space-html",
			}

			templatesLoader.Templates = templates

			_, err := strategy.Dispatch(clientID, "user-123", options, conn)
			if err != nil {
				panic(err)
			}

			Expect(templatesLoader.ContentSuffix).To(Equal(models.UserBodyTemplateName))
			Expect(mailer.DeliverArguments).To(ContainElement(conn))
			Expect(mailer.DeliverArguments).To(ContainElement(templates))
			Expect(mailer.DeliverArguments).To(ContainElement(users))
			Expect(mailer.DeliverArguments).To(ContainElement(options))
			Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerOrganization{}))
			Expect(mailer.DeliverArguments).To(ContainElement(cf.CloudControllerSpace{}))
			Expect(mailer.DeliverArguments).To(ContainElement(clientID))

			Expect(userLoader.LoadedGUIDs).To(Equal([]string{"user-123"}))
		})

		Context("failure cases", func() {
			Context("when a token cannot be loaded", func() {
				It("returns the error", func() {
					loadError := errors.New("BOOM!")
					tokenLoader.LoadError = loadError
					_, err := strategy.Dispatch(clientID, "user-123", options, conn)

					Expect(err).To(Equal(loadError))
				})
			})

			Context("when a user cannot be loaded", func() {
				It("returns the error", func() {
					loadError := errors.New("BOOM!")
					userLoader.LoadError = loadError
					_, err := strategy.Dispatch(clientID, "user-123", options, conn)

					Expect(err).To(Equal(loadError))
				})
			})

			Context("when a template cannot be loaded", func() {
				It("returns a TemplateLoadError", func() {
					templatesLoader.LoadError = errors.New("BOOM!")

					_, err := strategy.Dispatch(clientID, "user-123", options, conn)

					Expect(err).To(BeAssignableToTypeOf(postal.TemplateLoadError("")))
				})
			})

			Context("when create receipts call returns an err", func() {
				It("returns an error", func() {
					receiptsRepo.CreateReceiptsError = true

					_, err := strategy.Dispatch(clientID, "user-123", options, conn)
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})

	Describe("Trim", func() {
		Describe("TrimFields", func() {
			It("trims the specified fields from the response object", func() {
				responses, err := json.Marshal([]strategies.Response{
					{
						Status:         "delivered",
						Recipient:      "user-123",
						NotificationID: "123-456",
					},
				})

				trimmedResponses := strategy.Trim(responses)

				var result []map[string]string
				err = json.Unmarshal(trimmedResponses, &result)
				if err != nil {
					panic(err)
				}

				Expect(result).To(ContainElement(map[string]string{"status": "delivered",
					"recipient":       "user-123",
					"notification_id": "123-456",
				}))
			})
		})
	})
})