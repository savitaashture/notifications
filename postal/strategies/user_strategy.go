package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

const UserEndorsement = "This message was sent directly to you."

type UserStrategy struct {
	mailer MailerInterface
}

func NewUserStrategy(mailer MailerInterface) UserStrategy {
	return UserStrategy{
		mailer: mailer,
	}
}

func (strategy UserStrategy) Dispatch(clientID, guid string, options postal.Options, conn models.ConnectionInterface) ([]Response, error) {
	options.Endorsement = UserEndorsement
	responses := strategy.mailer.Deliver(conn, []User{{GUID: guid}}, options, cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, clientID, "")

	return responses, nil
}
