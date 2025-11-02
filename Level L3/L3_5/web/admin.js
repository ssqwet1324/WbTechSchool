const API_BASE = 'http://localhost:8081';

// Показать сообщение
function showMessage(message, type = 'info') {
    const messageArea = document.getElementById('messageArea');
    const alertClass = type === 'error' ? 'alert-error' : type === 'success' ? 'alert-success' : 'alert-info';
    messageArea.innerHTML = `<div class="alert ${alertClass}">${message}</div>`;
    setTimeout(() => {
        messageArea.innerHTML = '';
    }, 5000);
}

// Форматирование даты
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU');
}

// Подсчет статистики мест
function countSeatsByStatus(seats) {
    const stats = { free: 0, reserving: 0, booked: 0 };
    seats.forEach(seat => {
        if (stats[seat.status] !== undefined) {
            stats[seat.status]++;
        }
    });
    return stats;
}

// Создание мероприятия
document.getElementById('createEventForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const eventId = document.getElementById('eventId').value;
    const title = document.getElementById('eventTitle').value;
    const date = document.getElementById('eventDate').value;
    const totalSeats = parseInt(document.getElementById('totalSeats').value);
    const rows = parseInt(document.getElementById('rows').value);
    const seatsPerRow = parseInt(document.getElementById('seatsPerRow').value);
    const startNumber = parseInt(document.getElementById('startNumber').value);

    const eventData = {
        event: {
            id: eventId,
            title: title,
            date: new Date(date).toISOString(),
            total_seats: totalSeats
        },
        layout: {
            rows: rows,
            seats_per_row: seatsPerRow,
            start_number: startNumber
        }
    };

    try {
        const response = await fetch(`${API_BASE}/events`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(eventData)
        });

        const result = await response.json();

        if (response.ok) {
            showMessage(`Мероприятие "${title}" успешно создано! ID: ${result.event_id}`, 'success');
            document.getElementById('createEventForm').reset();
            loadEvents();
        } else {
            showMessage(`Ошибка: ${result.error}`, 'error');
        }
    } catch (error) {
        showMessage(`Ошибка при создании мероприятия: ${error.message}`, 'error');
    }
});

// Загрузка списка мероприятий
async function loadEvents() {
    const eventsList = document.getElementById('eventsList');
    eventsList.innerHTML = '<div class="loading">Загрузка...</div>';

    try {
        const response = await fetch(`${API_BASE}/events/all`);
        const events = await response.json();

        if (!response.ok) {
            throw new Error(events.error || 'Ошибка загрузки мероприятий');
        }

        if (events.length === 0) {
            eventsList.innerHTML = '<p>Мероприятия не найдены</p>';
            return;
        }

        eventsList.innerHTML = '';

        for (const event of events) {
            // Загружаем детали каждого мероприятия
            try {
                const eventDetailsResponse = await fetch(`${API_BASE}/events/${event.id}`);
                const eventDetails = await eventDetailsResponse.json();

                if (eventDetailsResponse.ok && eventDetails.seats) {
                    const stats = countSeatsByStatus(eventDetails.seats);

                    const eventCard = document.createElement('div');
                    eventCard.className = 'event-card';
                    eventCard.innerHTML = `
                        <h3>${event.title}</h3>
                        <div class="event-info">
                            <div class="event-info-item">
                                <strong>ID:</strong> ${event.id}
                            </div>
                            <div class="event-info-item">
                                <strong>Дата:</strong> ${formatDate(event.date)}
                            </div>
                            <div class="event-info-item">
                                <strong>Всего мест:</strong> ${event.total_seats}
                            </div>
                            <div class="event-info-item">
                                <strong>Свободно:</strong> <span class="status-badge status-free">${stats.free}</span>
                            </div>
                            <div class="event-info-item">
                                <strong>В резерве:</strong> <span class="status-badge status-reserving">${stats.reserving}</span>
                            </div>
                            <div class="event-info-item">
                                <strong>Забронировано:</strong> <span class="status-badge status-booked">${stats.booked}</span>
                            </div>
                        </div>
                    `;
                    eventsList.appendChild(eventCard);
                } else {
                    // Если детали не загрузились, показываем базовую информацию
                    const eventCard = document.createElement('div');
                    eventCard.className = 'event-card';
                    eventCard.innerHTML = `
                        <h3>${event.title}</h3>
                        <div class="event-info">
                            <div class="event-info-item"><strong>ID:</strong> ${event.id}</div>
                            <div class="event-info-item"><strong>Дата:</strong> ${formatDate(event.date)}</div>
                            <div class="event-info-item"><strong>Всего мест:</strong> ${event.total_seats}</div>
                        </div>
                    `;
                    eventsList.appendChild(eventCard);
                }
            } catch (error) {
                console.error(`Ошибка загрузки деталей мероприятия ${event.id}:`, error);
                // Показываем базовую информацию даже при ошибке
                const eventCard = document.createElement('div');
                eventCard.className = 'event-card';
                eventCard.innerHTML = `
                    <h3>${event.title}</h3>
                    <div class="event-info">
                        <div class="event-info-item"><strong>ID:</strong> ${event.id}</div>
                        <div class="event-info-item"><strong>Дата:</strong> ${formatDate(event.date)}</div>
                        <div class="event-info-item"><strong>Всего мест:</strong> ${event.total_seats}</div>
                    </div>
                    <p style="color: #e74c3c;">Не удалось загрузить детальную информацию</p>
                `;
                eventsList.appendChild(eventCard);
            }
        }
    } catch (error) {
        eventsList.innerHTML = `<div class="alert alert-error">Ошибка загрузки мероприятий: ${error.message}</div>`;
    }
}

// Загружаем мероприятия при загрузке страницы
loadEvents();

