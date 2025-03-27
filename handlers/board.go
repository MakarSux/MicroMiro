package handlers

import (
	"database/sql"
	"micromiro/database"
	"micromiro/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateBoard создает новую доску
func CreateBoard(c *gin.Context) {
	var req models.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	query := `INSERT INTO boards (title, description, creator_id, is_public, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var boardID int
	err = db.QueryRow(query, req.Title, req.Description, userID, req.IsPublic, time.Now(), time.Now()).Scan(&boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания доски"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Доска успешно создана", "board_id": boardID})
}

// GetBoards получает список досок пользователя
func GetBoards(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Получаем доски, созданные пользователем
	query := `SELECT id, title, description, creator_id, is_public, created_at, updated_at 
              FROM boards 
              WHERE creator_id = $1`

	rows, err := db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения досок"})
		return
	}
	defer rows.Close()

	boards := []models.Board{}
	for rows.Next() {
		var board models.Board
		if err := rows.Scan(&board.ID, &board.Title, &board.Description, &board.CreatorID, &board.IsPublic, &board.CreatedAt, &board.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
			return
		}
		boards = append(boards, board)
	}

	// Также получаем доски, к которым у пользователя есть доступ через разрешения
	query = `SELECT b.id, b.title, b.description, b.creator_id, b.is_public, b.created_at, b.updated_at 
             FROM boards b
             JOIN board_permissions bp ON b.id = bp.board_id
             WHERE bp.user_id = $1 AND b.creator_id != $1`

	rows, err = db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения досок с разрешениями"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var board models.Board
		if err := rows.Scan(&board.ID, &board.Title, &board.Description, &board.CreatorID, &board.IsPublic, &board.CreatedAt, &board.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
			return
		}
		boards = append(boards, board)
	}

	c.JSON(http.StatusOK, boards)
}

// GetBoard получает информацию о конкретной доске
func GetBoard(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID доски"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Проверяем, имеет ли пользователь доступ к доске
	var board models.Board
	query := `SELECT id, title, description, creator_id, is_public, created_at, updated_at 
              FROM boards 
              WHERE id = $1 AND (creator_id = $2 OR is_public = true OR EXISTS (
                  SELECT 1 FROM board_permissions WHERE board_id = $1 AND user_id = $2
              ))`

	err = db.QueryRow(query, boardID, userID).Scan(&board.ID, &board.Title, &board.Description, &board.CreatorID, &board.IsPublic, &board.CreatedAt, &board.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Доска не найдена или у вас нет доступа"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доски"})
		return
	}

	// Получаем элементы доски
	query = `SELECT id, board_id, type, content, position_x, position_y, width, height, created_at, updated_at 
             FROM board_elements 
             WHERE board_id = $1`

	rows, err := db.Query(query, boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения элементов доски"})
		return
	}
	defer rows.Close()

	elements := []models.BoardElement{}
	for rows.Next() {
		var element models.BoardElement
		if err := rows.Scan(&element.ID, &element.BoardID, &element.Type, &element.Content, &element.PositionX, &element.PositionY, &element.Width, &element.Height, &element.CreatedAt, &element.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных элементов"})
			return
		}
		elements = append(elements, element)
	}

	c.JSON(http.StatusOK, gin.H{"board": board, "elements": elements})
}

// UpdateBoard обновляет информацию о доске
func UpdateBoard(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID доски"})
		return
	}

	var req models.UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Проверяем, имеет ли пользователь права на редактирование доски
	var creatorID int
	query := `SELECT creator_id FROM boards WHERE id = $1`
	err = db.QueryRow(query, boardID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Доска не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доски"})
		return
	}

	// Проверяем, является ли пользователь создателем доски или имеет права на редактирование
	if creatorID != userID.(int) {
		var canEdit bool
		query = `SELECT can_edit FROM board_permissions WHERE board_id = $1 AND user_id = $2`
		err = db.QueryRow(query, boardID, userID).Scan(&canEdit)
		if err == sql.ErrNoRows || !canEdit {
			c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет прав на редактирование этой доски"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки прав доступа"})
			return
		}
	}

	// Обновляем доску
	query = `UPDATE boards SET title = $1, description = $2, is_public = $3, updated_at = $4 WHERE id = $5`
	_, err = db.Exec(query, req.Title, req.Description, req.IsPublic, time.Now(), boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления доски"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Доска успешно обновлена"})
}

// DeleteBoard удаляет доску
func DeleteBoard(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID доски"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Проверяем, является ли пользователь создателем доски
	var creatorID int
	query := `SELECT creator_id FROM boards WHERE id = $1`
	err = db.QueryRow(query, boardID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Доска не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доски"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только создатель доски может удалить её"})
		return
	}

	// Начинаем транзакцию для удаления доски и всех связанных данных
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка начала транзакции"})
		return
	}

	// Удаляем элементы доски
	_, err = tx.Exec(`DELETE FROM board_elements WHERE board_id = $1`, boardID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления элементов доски"})
		return
	}

	// Удаляем разрешения доски
	_, err = tx.Exec(`DELETE FROM board_permissions WHERE board_id = $1`, boardID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления разрешений доски"})
		return
	}

	// Удаляем саму доску
	_, err = tx.Exec(`DELETE FROM boards WHERE id = $1`, boardID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления доски"})
		return
	}

	// Завершаем транзакцию
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка завершения транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Доска успешно удалена"})
}

// CreateBoardElement создает новый элемент на доске
func CreateBoardElement(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID доски"})
		return
	}

	var req models.CreateBoardElementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Проверяем, имеет ли пользователь права на редактирование доски
	var creatorID int
	query := `SELECT creator_id FROM boards WHERE id = $1`
	err = db.QueryRow(query, boardID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Доска не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доски"})
		return
	}

	// Проверяем, является ли пользователь создателем доски или имеет права на редактирование
	if creatorID != userID.(int) {
		var canEdit bool
		query = `SELECT can_edit FROM board_permissions WHERE board_id = $1 AND user_id = $2`
		err = db.QueryRow(query, boardID, userID).Scan(&canEdit)
		if err == sql.ErrNoRows || !canEdit {
			c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет прав на редактирование этой доски"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки прав доступа"})
			return
		}
	}

	// Добавляем новый элемент
	query = `INSERT INTO board_elements (board_id, type, content, position_x, position_y, width, height, created_at, updated_at) 
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var elementID int
	err = db.QueryRow(query, boardID, req.Type, req.Content, req.PositionX, req.PositionY, req.Width, req.Height, time.Now(), time.Now()).Scan(&elementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания элемента доски"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Элемент успешно создан", "element_id": elementID})
}

// UpdateBoardElement обновляет элемент на доске
func UpdateBoardElement(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID доски"})
		return
	}

	elementID, err := strconv.Atoi(c.Param("element_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID элемента"})
		return
	}

	var req models.UpdateBoardElementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Проверяем, имеет ли пользователь права на редактирование доски
	var creatorID int
	query := `SELECT creator_id FROM boards WHERE id = $1`
	err = db.QueryRow(query, boardID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Доска не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доски"})
		return
	}

	// Проверяем, является ли пользователь создателем доски или имеет права на редактирование
	if creatorID != userID.(int) {
		var canEdit bool
		query = `SELECT can_edit FROM board_permissions WHERE board_id = $1 AND user_id = $2`
		err = db.QueryRow(query, boardID, userID).Scan(&canEdit)
		if err == sql.ErrNoRows || !canEdit {
			c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет прав на редактирование этой доски"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки прав доступа"})
			return
		}
	}

	// Проверяем, существует ли элемент и принадлежит ли он указанной доске
	var count int
	query = `SELECT COUNT(*) FROM board_elements WHERE id = $1 AND board_id = $2`
	err = db.QueryRow(query, elementID, boardID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки элемента"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден или не принадлежит указанной доске"})
		return
	}

	// Обновляем элемент
	query = `UPDATE board_elements SET type = $1, content = $2, position_x = $3, position_y = $4, width = $5, height = $6, updated_at = $7 
             WHERE id = $8 AND board_id = $9`
	_, err = db.Exec(query, req.Type, req.Content, req.PositionX, req.PositionY, req.Width, req.Height, time.Now(), elementID, boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления элемента"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Элемент успешно обновлен"})
}

// DeleteBoardElement удаляет элемент с доски
func DeleteBoardElement(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID доски"})
		return
	}

	elementID, err := strconv.Atoi(c.Param("element_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID элемента"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Проверяем, имеет ли пользователь права на редактирование доски
	var creatorID int
	query := `SELECT creator_id FROM boards WHERE id = $1`
	err = db.QueryRow(query, boardID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Доска не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доски"})
		return
	}

	// Проверяем, является ли пользователь создателем доски или имеет права на редактирование
	if creatorID != userID.(int) {
		var canEdit bool
		query = `SELECT can_edit FROM board_permissions WHERE board_id = $1 AND user_id = $2`
		err = db.QueryRow(query, boardID, userID).Scan(&canEdit)
		if err == sql.ErrNoRows || !canEdit {
			c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет прав на редактирование этой доски"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки прав доступа"})
			return
		}
	}

	// Проверяем, существует ли элемент и принадлежит ли он указанной доске
	var count int
	query = `SELECT COUNT(*) FROM board_elements WHERE id = $1 AND board_id = $2`
	err = db.QueryRow(query, elementID, boardID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки элемента"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден или не принадлежит указанной доске"})
		return
	}

	// Удаляем элемент
	query = `DELETE FROM board_elements WHERE id = $1 AND board_id = $2`
	_, err = db.Exec(query, elementID, boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления элемента"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Элемент успешно удален"})
}
