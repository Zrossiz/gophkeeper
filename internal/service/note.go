package service

import (
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

type NoteService struct {
	noteDB       NoteStorage
	cryptoModule CryptoModule
	log          *zap.Logger
}

type NoteStorage interface {
	Create(body dto.CreateNoteDTO) error
	Update(noteID int, body dto.UpdateNoteDTO) error
	GetAllByUser(userID int) ([]entities.Note, error)
}

func NewNoteService(
	db NoteStorage,
	cryptoModule CryptoModule,
	log *zap.Logger,
) *NoteService {
	return &NoteService{
		noteDB:       db,
		cryptoModule: cryptoModule,
		log:          log,
	}
}

func (n *NoteService) Create(body dto.CreateNoteDTO) error {
	encryptedTitle, err := n.cryptoModule.Encrypt(body.Title, body.Key)
	if err != nil {
		return err
	}

	encryptedTextData, err := n.cryptoModule.Encrypt(body.TextData, body.Key)
	if err != nil {
		return err
	}

	body.Title = encryptedTitle
	body.TextData = encryptedTextData

	err = n.noteDB.Create(body)
	if err != nil {
		return err
	}

	return nil
}

func (n *NoteService) Update(noteID int, body dto.UpdateNoteDTO) error {
	encryptedTitle, err := n.cryptoModule.Encrypt(body.Title, body.Key)
	if err != nil {
		return err
	}

	encryptedTextData, err := n.cryptoModule.Encrypt(body.TextData, body.Key)
	if err != nil {
		return err
	}

	body.Title = encryptedTitle
	body.TextData = encryptedTextData

	err = n.noteDB.Update(noteID, body)
	if err != nil {
		return err
	}

	return nil
}

func (n *NoteService) GetAll(userID int, key string) ([]entities.Note, error) {
	encryptedData, err := n.noteDB.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	decryptedData := n.decryptNotesArray(encryptedData, key)

	return decryptedData, nil
}

func (n *NoteService) decryptNotesArray(
	encryptedData []entities.Note,
	key string,
) []entities.Note {
	decryptedData := make([]entities.Note, 0, len(encryptedData))

	for i := 0; i < len(encryptedData); i++ {
		decryptedNote, err := n.decryptNote(encryptedData[i], key)
		if err != nil {
			continue
		}

		decryptedData = append(decryptedData, *decryptedNote)
	}

	return decryptedData
}

func (n *NoteService) decryptNote(
	encryptedNote entities.Note,
	key string,
) (*entities.Note, error) {
	decryptedTitle, err := n.cryptoModule.Decrypt(encryptedNote.Title, key)
	if err != nil {
		return nil, err
	}

	decryptedTextData, err := n.cryptoModule.Decrypt(encryptedNote.TextData, key)
	if err != nil {
		return nil, err
	}

	encryptedNote.Title = decryptedTitle
	encryptedNote.TextData = decryptedTextData

	return &encryptedNote, nil
}
