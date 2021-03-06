package utilities_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AllUserGUIDs", func() {
	var allUsers utilities.AllUsers
	var uaaClient *fakes.UAAClient
	var users []uaa.User

	BeforeEach(func() {
		uaaClient = fakes.NewUAAClient()
		allUsers = utilities.NewAllUsers(uaaClient)
	})

	Context("when the request succeeds", func() {
		BeforeEach(func() {
			users = []uaa.User{
				{
					Emails: []string{"user-123@example.com"},
					ID:     "user-123",
				},
				{
					Emails: []string{"user-456@example.com"},
					ID:     "user-456",
				},
				{
					Emails: []string{"user-999@example.com"},
					ID:     "user-999",
				},
			}

			uaaClient.AllUsersData = users
		})

		It("returns the UAAUsers, UserGUIDs, and an error", func() {
			guids, err := allUsers.AllUserGUIDs()
			Expect(err).NotTo(HaveOccurred())
			Expect(guids).To(ConsistOf("user-456", "user-999", "user-123"))
		})
	})

	Context("when the request to UAA fails", func() {
		It("bubbles up the error", func() {
			uaaClient.AllUsersError = errors.New("BOOM!")
			_, err := allUsers.AllUserGUIDs()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(errors.New("BOOM!")))
		})
	})
})
