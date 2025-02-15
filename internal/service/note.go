// Package service provides business logic for managing encrypted notes.
package service

import (
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// NoteService handles operations related to encrypted note storage.
type NoteService struct {
	noteDB       NoteStorage
	cryptoModule CryptoModule
	log          *zap.Logger
}

// NoteStorage defines an interface for storing, retrieving, and updating encrypted notes.
type NoteStorage interface {
	// Create stores an encrypted note entry.
	Create(ctx context.Context, body dto.CreateNoteDTO) error
	// Update modifies an existing encrypted note.
	Update(ctx context.Context, noteID int, body dto.UpdateNoteDTO) error
	// GetAllByUser retrieves all encrypted notes for a given user ID.
	GetAllByUser(ctx context.Context, userID int) ([]entities.Note, error)
}

// NewNoteService creates a new instance of NoteService with the provided dependencies.
//
// Parameters:
//   - db: An implementation of the NoteStorage interface for data persistence.
//   - cryptoModule: An implementation of CryptoModule for encryption and decryption.
//   - log: A structured logger (zap.Logger) for logging events.
//
// Returns:
//   - A pointer to a NoteService instance.
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

// Create encrypts and stores a note securely.
//
// Parameters:
//   - body: A dto.CreateNoteDTO containing note title, text, and an encryption key.
//
// Returns:
//   - An error if encryption or storage fails.
func (n *NoteService) Create(ctx context.Context, body dto.CreateNoteDTO) error {
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

	return n.noteDB.Create(ctx, body)
}

// Update encrypts and updates an existing note entry.
//
// Parameters:
//   - noteID: The ID of the note being updated.
//   - body: A dto.UpdateNoteDTO containing updated title, text, and an encryption key.
//
// Returns:
//   - An error if encryption or update fails.
func (n *NoteService) Update(ctx context.Context, noteID int, body dto.UpdateNoteDTO) error {
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

	return n.noteDB.Update(ctx, noteID, body)
}

// GetAll retrieves and decrypts all notes for a given user.
//
// Parameters:
//   - userID: The ID of the user whose notes are being retrieved.
//   - key: The encryption key required for decryption.
//
// Returns:
//   - A slice of decrypted entities.Note or an error if retrieval or decryption fails.
func (n *NoteService) GetAll(ctx context.Context, userID int, key string) ([]entities.Note, error) {
	encryptedData, err := n.noteDB.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	decryptedData := n.decryptNotesArray(encryptedData, key)

	return decryptedData, nil
}

// decryptNotesArray decrypts an array of encrypted notes.
//
// Parameters:
//   - encryptedData: A slice of encrypted entities.Note.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A slice of decrypted entities.Note.
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

// decryptNote decrypts a single encrypted note.
//
// Parameters:
//   - encryptedNote: An encrypted entities.Note instance.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A pointer to a decrypted entities.Note or an error if decryption fails.
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
