package postgres

import (
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCardStorage_CreateCard(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewCardStorage(db)

	// Создаём корректную карту
	cardDTO := dto.CreateCardDTO{
		UserID:         1,
		BankName:       "Test Bank",
		Num:            "1234567812345678", // Правильный номер карты
		CVV:            "123",
		ExpDate:        "12/25",
		CardHolderName: "Test User",
	}

	// Сохраняем карту
	err := storage.CreateCard(cardDTO)
	assert.NoError(t, err, "CreateCard should not return an error")

	// Проверяем, что данные карты были сохранены в базе
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM cards WHERE user_id = $1 AND bank_name = $2", cardDTO.UserID, cardDTO.BankName).Scan(&count)
	assert.NoError(t, err, "Failed to query cards table")
	assert.Equal(t, 1, count, "Expected one card to be inserted")
}

func TestCardStorage_GetAllCardsByUserId(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewCardStorage(db)

	// Создаём карту для пользователя
	cardDTO := dto.CreateCardDTO{
		UserID:         1,
		BankName:       "Test Bank",
		Num:            "1234567812345678", // Правильный номер карты
		CVV:            "123",
		ExpDate:        "12/25",
		CardHolderName: "Test User",
	}

	// Сохраняем карту
	err := storage.CreateCard(cardDTO)
	assert.NoError(t, err, "CreateCard should not return an error")

	// Получаем все карты для пользователя
	cards, err := storage.GetAllCardsByUserId(int64(cardDTO.UserID))
	assert.NoError(t, err, "GetAllCardsByUserId should not return an error")
	assert.Len(t, cards, 1, "Expected one card for the user")

	// Проверяем, что сохранённые данные совпадают с теми, что были отправлены
	assert.Equal(t, cardDTO.Num, cards[0].Number, "Card number should match")
	assert.Equal(t, cardDTO.CVV, cards[0].CVV, "CVV should match")
	assert.Equal(t, cardDTO.ExpDate, cards[0].ExpDate, "Expiration date should match")
	assert.Equal(t, cardDTO.CardHolderName, cards[0].CardHolderName, "Card holder name should match")
}

func TestCardStorage_UpdateCard(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewCardStorage(db)

	// Создаём карту для обновления
	cardDTO := dto.CreateCardDTO{
		UserID:         1,
		BankName:       "Test Bank",
		Num:            "1234567812345678", // Правильный номер карты
		CVV:            "123",
		ExpDate:        "12/25",
		CardHolderName: "Test User",
	}

	// Сохраняем карту
	err := storage.CreateCard(cardDTO)
	assert.NoError(t, err, "CreateCard should not return an error")

	// Обновляем данные карты
	updatedCardDTO := dto.UpdateCardDTO{
		Num:            "8765432187654321", // Новый номер карты
		CVV:            "321",
		ExpDate:        "12/30",
		CardHolderName: "Updated User",
	}

	// Обновляем карту
	err = storage.UpdateCard(1, updatedCardDTO) // Предположим, что ID карты - 1
	assert.NoError(t, err, "UpdateCard should not return an error")

	// Проверяем, что данные были обновлены
	var updatedCard entities.Card
	err = db.QueryRow("SELECT num, cvv, exp_date, card_holder_name FROM cards WHERE id = $1", 1).
		Scan(&updatedCard.Number, &updatedCard.CVV, &updatedCard.ExpDate, &updatedCard.CardHolderName)
	assert.NoError(t, err, "Failed to query updated card data")

	// Проверяем обновлённые данные
	assert.Equal(t, updatedCardDTO.Num, updatedCard.Number, "Card number should be updated")
	assert.Equal(t, updatedCardDTO.CVV, updatedCard.CVV, "CVV should be updated")
	assert.Equal(t, updatedCardDTO.ExpDate, updatedCard.ExpDate, "Expiration date should be updated")
	assert.Equal(t, updatedCardDTO.CardHolderName, updatedCard.CardHolderName, "Card holder name should be updated")
}
