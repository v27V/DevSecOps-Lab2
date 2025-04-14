package storage

import (
	"fmt"
	"sync"

	"project/internal/models"
)

// Database имитирует базу данных в памяти
type Database struct {
	Name     string
	Products map[int]models.Product
	mu       sync.RWMutex
}

// DBManager управляет несколькими базами данных
type DBManager struct {
	Databases map[string]*Database
}

// NewDBManager создает и инициализирует менеджер баз данных
func NewDBManager() *DBManager {
	manager := &DBManager{
		Databases: make(map[string]*Database),
	}

	// Создаем три "базы данных"
	manager.Databases["products_db"] = &Database{
		Name:     "products_db",
		Products: make(map[int]models.Product),
	}
	manager.Databases["suppliers_db"] = &Database{
		Name:     "suppliers_db",
		Products: make(map[int]models.Product),
	}
	manager.Databases["inventory_db"] = &Database{
		Name:     "inventory_db",
		Products: make(map[int]models.Product),
	}

	// Заполняем тестовыми данными
	manager.initializeTestData()

	return manager
}

// Инициализация тестовых данных
func (m *DBManager) initializeTestData() {
	// Данные для products_db
	m.Databases["products_db"].Products[1] = models.Product{
		ID:          1,
		Name:        "Кирпич облицовочный",
		Category:    "Стеновые материалы",
		Price:       15.50,
		Description: "Кирпич керамический облицовочный",
		InStock:     true,
		Supplier:    "ООО Кирпичный завод",
	}
	m.Databases["products_db"].Products[2] = models.Product{
		ID:          2,
		Name:        "Цемент М500",
		Category:    "Вяжущие материалы",
		Price:       350.00,
		Description: "Цемент М500 Д0, мешок 50 кг",
		InStock:     true,
		Supplier:    "Евроцемент",
	}

	// Данные для suppliers_db
	m.Databases["suppliers_db"].Products[1] = models.Product{
		ID:          1,
		Name:        "Клей для плитки",
		Category:    "Клеевые составы",
		Price:       280.00,
		Description: "Клей для керамической плитки, 25 кг",
		InStock:     true,
		Supplier:    "Цемикс",
	}

	// Данные для inventory_db
	m.Databases["inventory_db"].Products[1] = models.Product{
		ID:          1,
		Name:        "Гипсокартон",
		Category:    "Листовые материалы",
		Price:       450.00,
		Description: "Гипсокартон влагостойкий, 12.5 мм, 1.2x2.5 м",
		InStock:     false,
		Supplier:    "Кнауф",
	}
}

// GetDB возвращает базу данных по имени
func (m *DBManager) GetDB(name string) (*Database, error) {
	db, exists := m.Databases[name]
	if !exists {
		return nil, fmt.Errorf("база данных '%s' не найдена", name)
	}
	return db, nil
}

// GetProduct возвращает продукт по ID
func (db *Database) GetProduct(id int) (models.Product, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	product, exists := db.Products[id]
	return product, exists
}

// GetAllProducts возвращает все продукты из базы данных
func (db *Database) GetAllProducts() []models.Product {
	db.mu.RLock()
	defer db.mu.RUnlock()

	products := make([]models.Product, 0, len(db.Products))
	for _, product := range db.Products {
		products = append(products, product)
	}
	return products
}

// AddProduct добавляет новый продукт в базу данных
func (db *Database) AddProduct(product models.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Проверяем, существует ли продукт с таким ID
	if _, exists := db.Products[product.ID]; exists {
		return fmt.Errorf("продукт с ID %d уже существует", product.ID)
	}

	db.Products[product.ID] = product
	return nil
}

// UpdateProduct обновляет существующий продукт
func (db *Database) UpdateProduct(product models.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Проверяем, существует ли продукт с таким ID
	if _, exists := db.Products[product.ID]; !exists {
		return fmt.Errorf("продукт с ID %d не найден", product.ID)
	}

	db.Products[product.ID] = product
	return nil
}

// DeleteProduct удаляет продукт по ID
func (db *Database) DeleteProduct(id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Проверяем, существует ли продукт с таким ID
	if _, exists := db.Products[id]; !exists {
		return fmt.Errorf("продукт с ID %d не найден", id)
	}

	delete(db.Products, id)
	return nil
}
