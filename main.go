package main

import (
	"fmt"
	"sync"
	"time"
)

// Определение ролей пользователей
type Role string

const (
	ClientRole Role = "Клиент"
	BarberRole Role = "Барбер"
)

// Структуры данных
type Barber struct {
	ID      int
	Name    string
	Slots   []string // Доступные слоты для записи
	Reviews []string // Отзывы о барбере
}

type Client struct {
	ID   string
	Name string
	Role Role
}

type Appointment struct {
	ClientID string
	BarberID int
	Slot     string
	Date     string // Дата записи
}

var (
	barbers          = make(map[int]*Barber) // Используем указатели на Barber
	clients          = make(map[string]Client)
	appointments     = []Appointment{}
	appointmentsLog  = []Appointment{} // История записей
	barbersLock      sync.Mutex
	clientsLock      sync.Mutex
	appointmentsLock sync.Mutex
	nextBarberID     = 1
	clientCounter    = 1 // Счетчик клиентов за день
)

// Генерация уникального ID для клиента
func generateClientID() string {
	today := time.Now().Format("0201") // ДДММ
	id := fmt.Sprintf("%s%04d", today, clientCounter)
	clientCounter++
	return id
}

// Регистрация пользователя с выбором роли
func registerUser(name string, role Role) string {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	clientID := generateClientID()
	client := Client{
		ID:   clientID,
		Name: name,
		Role: role,
	}
	clients[clientID] = client

	fmt.Printf("✅ Пользователь %s зарегистрирован как %s (ID: %s)\n", name, role, clientID)
	return clientID
}

// Добавление нового барбера
func addBarber(name string, slots []string) {
	barbersLock.Lock()
	defer barbersLock.Unlock()

	barber := &Barber{
		ID:    nextBarberID,
		Name:  name,
		Slots: slots,
	}
	barbers[nextBarberID] = barber
	fmt.Printf("💈 Барбер %s добавлен с ID %d\n", name, nextBarberID)
	nextBarberID++
}

// Бронирование времени у барбера
func bookAppointment(clientID string) {
	barbersLock.Lock()
	appointmentsLock.Lock()
	defer barbersLock.Unlock()
	defer appointmentsLock.Unlock()

	var barberID int
	var slot string

	fmt.Print("Введите ID барбера: ")
	fmt.Scanln(&barberID)

	barber, exists := barbers[barberID]
	if !exists {
		fmt.Println("❌ Барбер не найден!")
		return
	}

	fmt.Print("Введите время (пример: 11:00): ")
	fmt.Scanln(&slot)

	for i, s := range barber.Slots {
		if s == slot {
			barber.Slots = append(barber.Slots[:i], barber.Slots[i+1:]...)
			appt := Appointment{ClientID: clientID, BarberID: barberID, Slot: slot, Date: time.Now().Format("02-01-2006")}
			appointments = append(appointments, appt)
			fmt.Printf("✅ Вы забронировали %s у барбера %s\n", slot, barber.Name)
			return
		}
	}
	fmt.Println("❌ Слот не найден!")
}

// Отмена записи
func cancelAppointment(clientID string) {
	appointmentsLock.Lock()
	defer appointmentsLock.Unlock()

	for i, appt := range appointments {
		if appt.ClientID == clientID {
			appointmentsLog = append(appointmentsLog, appt) // Переносим в историю
			appointments = append(appointments[:i], appointments[i+1:]...)
			if barber, exists := barbers[appt.BarberID]; exists {
				barber.Slots = append(barber.Slots, appt.Slot)
			}
			fmt.Printf("❌ Запись на %s отменена\n", appt.Slot)
			return
		}
	}
	fmt.Println("❌ Запись не найдена!")
}

// Показать текущую запись клиента
func showCurrentAppointment(clientID string) {
	for _, appt := range appointments {
		if appt.ClientID == clientID {
			if barber, exists := barbers[appt.BarberID]; exists {
				fmt.Printf("📌 Ваша текущая запись: %s у барбера %s (%s)\n", appt.Slot, barber.Name, appt.Date)
				return
			}
		}
	}
	fmt.Println("❌ У вас нет активных записей.")
}

// Показать историю записей клиента
func showAppointmentHistory(clientID string) {
	fmt.Printf("📖 История ваших записей:\n")
	found := false
	for _, appt := range appointmentsLog {
		if appt.ClientID == clientID {
			if barber, exists := barbers[appt.BarberID]; exists {
				fmt.Printf("- %s у барбера %s (%s)\n", appt.Slot, barber.Name, appt.Date)
				found = true
			}
		}
	}
	if !found {
		fmt.Println("❌ История записей пуста.")
	}
}

// Главное меню приложения (CLI)
func startCLI() {
	var userRole Role
	var clientID string

	fmt.Println("Выберите роль:")
	fmt.Println("1 - Я ищу барбера (Клиент)")
	fmt.Println("2 - Я барбер")
	var roleChoice int
	fmt.Scanln(&roleChoice)

	switch roleChoice {
	case 1:
		userRole = ClientRole
	case 2:
		userRole = BarberRole
	default:
		fmt.Println("Некорректный выбор.")
		return
	}

	var name string
	fmt.Print("Введите ваше имя: ")
	fmt.Scanln(&name)
	clientID = registerUser(name, userRole)

	for {
		fmt.Println("\nВыберите действие:")
		if userRole == ClientRole {
			fmt.Println("1 - Забронировать время")
			fmt.Println("2 - Посмотреть текущую запись")
			fmt.Println("3 - Посмотреть историю записей")
			fmt.Println("4 - Отменить запись")
		} else {
			fmt.Println("1 - Добавить барбера")
			fmt.Println("2 - Посмотреть доступные слоты")
			fmt.Println("3 - Забронировать время для клиента")
			fmt.Println("4 - Отменить запись клиента")
		}
		fmt.Println("0 - Выход")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 0:
			fmt.Println("Выход из программы.")
			return
		case 1:
			if userRole == ClientRole {
				bookAppointment(clientID)
			} else {
				var barberName string
				fmt.Print("Введите имя барбера: ")
				fmt.Scanln(&barberName)
				addBarber(barberName, []string{"10:00", "11:00", "12:00"})
			}
		case 2:
			if userRole == ClientRole {
				showCurrentAppointment(clientID)
			} else {
				fmt.Println("Функционал для барберов еще не реализован.")
			}
		case 3:
			if userRole == ClientRole {
				showAppointmentHistory(clientID)
			} else {
				fmt.Println("Функционал для барберов еще не реализован.")
			}
		case 4:
			cancelAppointment(clientID)
		default:
			fmt.Println("Некорректный выбор.")
		}
	}
}

// Запуск API-сервера
func main() {
	// Запускаем API сервер в фоне
	go startServer()

	// Запускаем CLI-интерфейс
	startCLI()
}
