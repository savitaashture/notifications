package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Updater", func() {
	Describe("#Update", func() {
		var templatesRepo *fakes.TemplatesRepo
		var template models.Template
		var updater services.TemplateUpdater

		BeforeEach(func() {
			templatesRepo = fakes.NewTemplatesRepo()
			template = models.Template{
				Name: "gobble template",
				Text: "gobble",
				HTML: "<p>gobble</p>",
			}

			updater = services.NewTemplateUpdater(templatesRepo, fakes.NewDatabase())
		})

		It("Inserts templates into the templates repo", func() {
			Expect(templatesRepo.Templates).ToNot(ContainElement(template))
			err := updater.Update("my-awesome-id", template)
			Expect(err).ToNot(HaveOccurred())
			Expect(templatesRepo.Templates).To(ContainElement(template))
		})

		It("propagates errors from repo", func() {
			expectedErr := errors.New("Boom!")

			templatesRepo.UpdateError = expectedErr
			err := updater.Update("unimportant", template)

			Expect(err).To(Equal(expectedErr))
		})
	})
})
