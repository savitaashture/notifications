package models

import (
	"database/sql"
	"time"
)

type MessagesRepo struct {
}

func NewMessagesRepo() MessagesRepo {
	return MessagesRepo{}
}

func (repo MessagesRepo) Create(conn ConnectionInterface, message Message) (Message, error) {
	message.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()
	err := conn.Insert(&message)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}

func (repo MessagesRepo) FindByID(conn ConnectionInterface, messageID string) (Message, error) {
	message := Message{}
	err := conn.SelectOne(&message, "SELECT * FROM `messages` WHERE `id`=?", messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Message{}, NewRecordNotFoundError("Message with ID %q could not be found", messageID)
		}
		return Message{}, err
	}
	return message, nil
}

func (repo MessagesRepo) Update(conn ConnectionInterface, message Message) (Message, error) {
	message.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()
	_, err := conn.Update(&message)
	if err != nil {
		return message, err
	}

	return repo.FindByID(conn, message.ID)
}

func (repo MessagesRepo) Upsert(conn ConnectionInterface, message Message) (Message, error) {
	_, err := repo.FindByID(conn, message.ID)

	switch err.(type) {
	case RecordNotFoundError:
		return repo.Create(conn, message)
	case nil:
		return repo.Update(conn, message)
	default:
		return message, err
	}
}

func (repo MessagesRepo) DeleteBefore(conn ConnectionInterface, threshold time.Time) (int, error) {
	result, err := conn.Exec("DELETE FROM `messages` WHERE `updated_at` < ?", threshold.UTC())
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
