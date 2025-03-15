package main

import (
	"fmt"
	"sync"
	"time"
)

// Структура данных

type Barber struct {
	ID    int
	Name  string
	Slots []string // Доступные слоты для записи
}

type Client struct {
	ID   string
	Name string
}

type Appointment struct {
	ClientID string
	BarberID int
	Slot     string
}

var (
	barbers          = make(map[int]*Barber) // Используем указатели на Barber
	clients          = make(map[string]Client)
	appointments     = []Appointment{}
	barbersLock      sync.Mutex
	clientsLock      sync.Mutex
	appointmentsLock sync.Mutex
	nextBarberID     = 1
	clientCounter    = 1 // Счетчик клиентов за день
)

// Генерация уникального ID для клиента
func generateClientID() string {
	today := time.Now().Format("02.01") // ДД.ММ
	id := fmt.Sprintf("%s.%04d", today, clientCounter)
	clientCounter++
	return id
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

// Добавление нового клиента
func addClient(name string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	clientID := generateClientID()
	client := Client{
		ID:   clientID,
		Name: name,
	}
	clients[clientID] = client
	fmt.Printf("👤 Клиент %s добавлен с ID %s\n", name, clientID)
}

// Получение доступных слотов барбера
func getAvailableSlots(barberID int) {
	barbersLock.Lock()
	defer barbersLock.Unlock()

	barber, exists := barbers[barberID]
	if !exists {
		fmt.Println("❌ Барбер не найден!")
		return
	}

	fmt.Printf("🗓 Доступные слоты у барбера %s:\n", barber.Name)
	for _, slot := range barber.Slots {
		fmt.Println("-", slot)
	}
}

// Бронирование времени у барбера
func bookAppointment(clientID string, barberID int, slot string) {
	barbersLock.Lock()
	appointmentsLock.Lock()
	defer barbersLock.Unlock()
	defer appointmentsLock.Unlock()

	barber, exists := barbers[barberID]
	if !exists {
		fmt.Println("❌ Барбер не найден!")
		return
	}

	// Проверяем, есть ли слот в доступных
	for i, s := range barber.Slots {
		if s == slot {
			// Удаляем слот из доступных
			barber.Slots = append(barber.Slots[:i], barber.Slots[i+1:]...)
			appointments = append(appointments, Appointment{ClientID: clientID, BarberID: barberID, Slot: slot})
			fmt.Printf("✅ Клиент %s забронировал %s у барбера %s\n", clientID, slot, barber.Name)
			return
		}
	}
	fmt.Println("❌ Слот не найден!")
}

// Отмена записи
func cancelAppointment(clientID string, slot string) {
	appointmentsLock.Lock()
	defer appointmentsLock.Unlock()

	for i, appt := range appointments {
		if appt.ClientID == clientID && appt.Slot == slot {
			appointments = append(appointments[:i], appointments[i+1:]...)
			if barber, exists := barbers[appt.BarberID]; exists {
				barber.Slots = append(barber.Slots, slot)
			}
			fmt.Printf("❌ Запись на %s отменена клиентом %s\n", slot, clientID)
			return
		}
	}
	fmt.Println("❌ Запись не найдена!")
}

// Показать записи клиента
func showAppointments(clientID string) {
	fmt.Printf("📋 Записи клиента %s:\n", clientID)
	for _, appt := range appointments {
		if appt.ClientID == clientID {
			if barber, exists := barbers[appt.BarberID]; exists {
				fmt.Printf("- %s с барбером %s\n", appt.Slot, barber.Name)
			}
		}
	}
}

func main() {
	for {
		var choice int
		fmt.Println("\nВыберите действие:")
		fmt.Println("1 - Добавить барбера")
		fmt.Println("2 - Добавить клиента")
		fmt.Println("3 - Посмотреть доступные слоты")
		fmt.Println("4 - Забронировать время")
		fmt.Println("5 - Отменить запись")
		fmt.Println("6 - Показать бронирования")
		fmt.Println("0 - Выход")
		fmt.Print("Ваш выбор: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var name string
			fmt.Print("Введите имя барбера: ")
			fmt.Scanln(&name)
			addBarber(name, []string{"10:00", "11:00", "12:00"})
		case 2:
			var name string
			fmt.Print("Введите имя клиента: ")
			fmt.Scanln(&name)
			addClient(name)
		case 3:
			var barberID int
			fmt.Print("Введите ID барбера: ")
			fmt.Scanln(&barberID)
			getAvailableSlots(barberID)
		case 4:
			var clientID string
			var barberID int
			var slot string
			fmt.Print("Введите ID клиента: ")
			fmt.Scanln(&clientID)
			fmt.Print("Введите ID барбера: ")
			fmt.Scanln(&barberID)
			fmt.Print("Введите время: ")
			fmt.Scanln(&slot)
			bookAppointment(clientID, barberID, slot)
		case 5:
			var clientID string
			var slot string
			fmt.Print("Введите ID клиента: ")
			fmt.Scanln(&clientID)
			fmt.Print("Введите время: ")
			fmt.Scanln(&slot)
			cancelAppointment(clientID, slot)
		case 6:
			var clientID string
			fmt.Print("Введите ID клиента: ")
			fmt.Scanln(&clientID)
			showAppointments(clientID)
		case 0:
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}
}
