package main

import (
	"fmt"
	"sync"
)

// Структура данных

type Barber struct {
	ID    int
	Name  string
	Slots []string // Доступные слоты для записи
}

type Client struct {
	ID   int
	Name string
}

type Appointment struct {
	ClientID int
	BarberID int
	Slot     string
}

var (
	barbers          = make(map[int]*Barber) // Изменено: теперь указатель на Barber
	clients          = make(map[int]Client)
	appointments     = []Appointment{}
	barbersLock      sync.Mutex
	clientsLock      sync.Mutex
	appointmentsLock sync.Mutex
	nextBarberID     = 1
	nextClientID     = 1
)

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

	client := Client{
		ID:   nextClientID,
		Name: name,
	}
	clients[nextClientID] = client
	fmt.Printf("👤 Клиент %s добавлен с ID %d\n", name, nextClientID)
	nextClientID++
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
func bookAppointment(clientID, barberID int, slot string) {
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
			fmt.Printf("✅ Клиент %d забронировал %s у барбера %s\n", clientID, slot, barber.Name)
			return
		}
	}
	fmt.Println("❌ Слот не найден!")
}

// Отмена записи
func cancelAppointment(clientID int, slot string) {
	appointmentsLock.Lock()
	defer appointmentsLock.Unlock()

	for i, appt := range appointments {
		if appt.ClientID == clientID && appt.Slot == slot {
			appointments = append(appointments[:i], appointments[i+1:]...)
			if barber, exists := barbers[appt.BarberID]; exists {
				barber.Slots = append(barber.Slots, slot)
			}
			fmt.Printf("❌ Запись на %s отменена клиентом %d\n", slot, clientID)
			return
		}
	}
	fmt.Println("❌ Запись не найдена!")
}

// Показать записи клиента
func showAppointments(clientID int) {
	fmt.Printf("📋 Записи клиента %d:\n", clientID)
	for _, appt := range appointments {
		if appt.ClientID == clientID {
			if barber, exists := barbers[appt.BarberID]; exists {
				fmt.Printf("- %s с барбером %s\n", appt.Slot, barber.Name)
			}
		}
	}
}

func main() {
	// Добавляем барберов и клиентов
	addBarber("Алексей", []string{"10:00", "11:00", "12:00"})
	addBarber("Иван", []string{"14:00", "15:00"})
	addClient("Петр")
	addClient("Сергей")

	// Демонстрация работы
	getAvailableSlots(1)
	bookAppointment(1, 1, "10:00")
	bookAppointment(2, 1, "12:00")
	showAppointments(1)
	showAppointments(2)
	cancelAppointment(1, "10:00")
	showAppointments(1)
	getAvailableSlots(1)
}
