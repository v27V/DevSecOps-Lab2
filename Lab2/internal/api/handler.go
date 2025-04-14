package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"project/internal/models"
	"project/internal/storage"
)

// APIHandler обрабатывает API запросы
type APIHandler struct {
	dbManager *storage.DBManager
}

// NewAPIHandler создает новый обработчик API
func NewAPIHandler(dbManager *storage.DBManager) *APIHandler {
	return &APIHandler{dbManager: dbManager}
}

// ServeHTTP обрабатывает HTTP запросы
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Парсим путь запроса
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "Неверный путь запроса", http.StatusBadRequest)
		return
	}

	dbName := pathParts[0]
	resource := pathParts[1]

	// Получаем базу данных
	db, err := h.dbManager.GetDB(dbName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Обрабатываем запрос в зависимости от метода HTTP
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r, db, resource, pathParts)
	case http.MethodPost:
		h.handlePost(w, r, db, resource)
	case http.MethodPut:
		h.handlePut(w, r, db, resource, pathParts)
	case http.MethodDelete:
		h.handleDelete(w, r, db, resource, pathParts)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// handleGet обрабатывает GET запросы
func (h *APIHandler) handleGet(w http.ResponseWriter, r *http.Request, db *storage.Database, resource string, pathParts []string) {
	if resource != "products" {
		http.Error(w, "Ресурс не найден", http.StatusNotFound)
		return
	}

	// Если в пути есть идентификатор - возвращаем один товар
	if len(pathParts) > 2 {
		id, err := strconv.Atoi(pathParts[2])
		if err != nil {
			http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
			return
		}

		product, exists := db.GetProduct(id)
		if !exists {
			http.Error(w, "Товар не найден", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(product)
		return
	}

	// Возвращаем все товары
	products := db.GetAllProducts()
	json.NewEncoder(w).Encode(products)
}

// handlePost обрабатывает POST запросы (создание нового товара)
func (h *APIHandler) handlePost(w http.ResponseWriter, r *http.Request, db *storage.Database, resource string) {
	if resource != "products" {
		http.Error(w, "Ресурс не найден", http.StatusNotFound)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Ошибка чтения данных: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляем товар в "базу данных"
	if err := db.AddProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Товар с ID %d успешно создан в базе %s", product.ID, db.Name),
	})
}

// handlePut обрабатывает PUT запросы (обновление товара)
func (h *APIHandler) handlePut(w http.ResponseWriter, r *http.Request, db *storage.Database, resource string, pathParts []string) {
	if resource != "products" || len(pathParts) <= 2 {
		http.Error(w, "Неверный путь запроса", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Ошибка чтения данных: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Удостоверимся, что ID в пути и в теле запроса совпадают
	if product.ID != id {
		http.Error(w, "ID в пути и в теле запроса не совпадают", http.StatusBadRequest)
		return
	}

	// Обновляем товар
	if err := db.UpdateProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Товар с ID %d успешно обновлен в базе %s", id, db.Name),
	})
}

// handleDelete обрабатывает DELETE запросы (удаление товара)
func (h *APIHandler) handleDelete(w http.ResponseWriter, r *http.Request, db *storage.Database, resource string, pathParts []string) {
	if resource != "products" || len(pathParts) <= 2 {
		http.Error(w, "Неверный путь запроса", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	// Удаляем товар
	if err := db.DeleteProduct(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Товар с ID %d успешно удален из базы %s", id, db.Name),
	})
}
