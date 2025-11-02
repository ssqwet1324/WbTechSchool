const API_BASE = 'http://localhost:8081';
let currentUserId = null;
let currentReservation = null;
let reservationTimer = null;

// Показать сообщение
function showMessage(message, type = 'info') {
    const messageArea = document.getElementById('messageArea');
    const alertClass = type === 'error' ? 'alert-error' : type === 'success' ? 'alert-success' : 'alert-info';
    messageArea.innerHTML = `<div class="alert ${alertClass}">${message}</div>`;
    setTimeout(() => {
        messageArea.innerHTML = '';
    }, 5000);
}

// Установка пользователя
function setUser() {
    const userIdInput = document.getElementById('userId');
    const userId = userIdInput.value.trim();

    if (!userId) {
        showMessage('Введите ID пользователя', 'error');
        return;
    }

    currentUserId = userId;
    document.getElementById('currentUser').textContent = userId;
    document.getElementById('userInfo').style.display = 'block';
    document.getElementById('eventsSection').style.display = 'block';
    
    // Проверяем сохраненную бронь при входе
    checkSavedReservation();
    
    loadEvents();
}

// Проверка сохраненной брони при загрузке страницы
function checkSavedReservation() {
    if (currentUserId) {
        const savedReservation = localStorage.getItem(`reservation_${currentUserId}`);
        if (savedReservation) {
            try {
                const reservation = JSON.parse(savedReservation);
                const reservationTime = new Date(reservation.timestamp);
                const now = new Date();
                const elapsedMinutes = (now - reservationTime) / 1000 / 60;

                if (elapsedMinutes < 10) {
                    currentReservation = { eventId: reservation.eventId, seatNumber: reservation.seatNumber };
                    // Автоматически открываем детали мероприятия с бронированием
                    setTimeout(() => {
                        viewEventDetails(reservation.eventId);
                    }, 500);
                } else {
                    localStorage.removeItem(`reservation_${currentUserId}`);
                }
            } catch (e) {
                console.error('Ошибка при чтении сохраненной брони:', e);
            }
        }
    }
}

// Форматирование даты
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU');
}

// Загрузка списка мероприятий
async function loadEvents() {
    if (!currentUserId) {
        showMessage('Сначала войдите в систему', 'error');
        return;
    }

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

        events.forEach(event => {
            const eventCard = document.createElement('div');
            eventCard.className = 'event-card';
            eventCard.innerHTML = `
                <h3>${event.title}</h3>
                <div class="event-info">
                    <div class="event-info-item"><strong>Дата:</strong> ${formatDate(event.date)}</div>
                    <div class="event-info-item"><strong>Всего мест:</strong> ${event.total_seats}</div>
                </div>
                <button onclick="viewEventDetails('${event.id}')" class="btn btn-primary" style="margin-top: 10px;">
                    Посмотреть места и забронировать
                </button>
            `;
            eventsList.appendChild(eventCard);
        });
    } catch (error) {
        eventsList.innerHTML = `<div class="alert alert-error">Ошибка загрузки мероприятий: ${error.message}</div>`;
    }
}

// Просмотр деталей мероприятия
async function viewEventDetails(eventId) {
    if (!currentUserId) {
        showMessage('Сначала войдите в систему', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/events/${eventId}`);
        const eventInfo = await response.json();

        if (!response.ok) {
            throw new Error(eventInfo.error || 'Ошибка загрузки мероприятия');
        }

        // Показываем секцию с деталями
        document.getElementById('eventDetailsSection').style.display = 'block';
        document.getElementById('eventDetailsTitle').textContent = eventInfo.title;
        
        // Отображаем информацию о мероприятии
        document.getElementById('eventDetails').innerHTML = `
            <div class="event-info">
                <div class="event-info-item"><strong>ID:</strong> ${eventInfo.id}</div>
                <div class="event-info-item"><strong>Дата:</strong> ${formatDate(eventInfo.date)}</div>
                <div class="event-info-item"><strong>Всего мест:</strong> ${eventInfo.total_seats}</div>
            </div>
        `;

        // Отображаем места
        renderSeats(eventInfo.seats, eventId);

        // Проверяем, есть ли активная бронь у этого пользователя
        checkActiveReservation(eventInfo.seats, eventId);

        // Прокручиваем к секции с местами
        document.getElementById('eventDetailsSection').scrollIntoView({ behavior: 'smooth' });
    } catch (error) {
        showMessage(`Ошибка загрузки мероприятия: ${error.message}`, 'error');
    }
}

// Отображение мест
function renderSeats(seats, eventId) {
    const seatsContainer = document.getElementById('seatsContainer');
    seatsContainer.innerHTML = '<h3>Выбор места:</h3><div class="seats-grid"></div>';
    const grid = seatsContainer.querySelector('.seats-grid');

    seats.forEach(seat => {
        const seatElement = document.createElement('div');
        seatElement.className = `seat ${seat.status}`;
        seatElement.textContent = seat.seat_number;
        seatElement.title = `Место ${seat.seat_number} - ${getStatusText(seat.status)}`;

        if (seat.status === 'free') {
            seatElement.onclick = () => bookSeat(eventId, seat.seat_number);
        } else {
            seatElement.classList.add('disabled');
        }

        grid.appendChild(seatElement);
    });
}

// Получение текста статуса
function getStatusText(status) {
    const statusMap = {
        'free': 'Свободно',
        'reserving': 'В резерве',
        'booked': 'Забронировано'
    };
    return statusMap[status] || status;
}

// Бронирование места
async function bookSeat(eventId, seatNumber) {
    if (!currentUserId) {
        showMessage('Сначала войдите в систему', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/events/${eventId}/book`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                seat_number: seatNumber,
                user_id: currentUserId
            })
        });

        const result = await response.json();

        if (response.ok) {
            showMessage(`Место ${seatNumber} успешно забронировано! У вас есть 10 минут для подтверждения оплаты.`, 'success');
            currentReservation = { eventId, seatNumber };
            
            // Сохраняем бронь в localStorage
            const reservationData = {
                eventId,
                seatNumber,
                userId: currentUserId,
                timestamp: new Date().toISOString()
            };
            localStorage.setItem(`reservation_${currentUserId}`, JSON.stringify(reservationData));
            
            startReservationTimer(reservationData.timestamp);
            viewEventDetails(eventId); // Перезагружаем детали
        } else {
            showMessage(`Ошибка: ${result.error}`, 'error');
        }
    } catch (error) {
        showMessage(`Ошибка при бронировании: ${error.message}`, 'error');
    }
}

// Проверка активной брони
function checkActiveReservation(seats, eventId) {
    if (!currentUserId) return;

    // Проверяем сохраненную бронь в localStorage
    const savedReservation = localStorage.getItem(`reservation_${currentUserId}`);
    
    if (savedReservation) {
        try {
            const reservation = JSON.parse(savedReservation);
            const reservationTime = new Date(reservation.timestamp);
            const now = new Date();
            const elapsedMinutes = (now - reservationTime) / 1000 / 60;

            // Если прошло меньше 10 минут и это тот же eventId
            if (elapsedMinutes < 10 && reservation.eventId === eventId) {
                const seat = seats.find(s => s.seat_number === reservation.seatNumber);
                if (seat && seat.status === 'reserving') {
                    currentReservation = { eventId, seatNumber: reservation.seatNumber };
                    showReservationInfo();
                    startReservationTimer(reservation.timestamp);
                    return;
                }
            } else if (elapsedMinutes >= 10) {
                // Бронь истекла, удаляем из localStorage
                localStorage.removeItem(`reservation_${currentUserId}`);
            }
        } catch (e) {
            console.error('Ошибка при чтении сохраненной брони:', e);
        }
    }

    hideReservationInfo();
}

// Показать информацию о брони
function showReservationInfo() {
    if (currentReservation) {
        document.getElementById('reservedSeat').textContent = currentReservation.seatNumber;
        document.getElementById('reservationStatus').textContent = 'Ожидает подтверждения';
        document.getElementById('reservationInfo').style.display = 'block';
    }
}

// Скрыть информацию о брони
function hideReservationInfo() {
    document.getElementById('reservationInfo').style.display = 'none';
    if (reservationTimer) {
        clearInterval(reservationTimer);
        reservationTimer = null;
    }
}

// Таймер бронирования
function startReservationTimer(reservationTimestamp) {
    if (reservationTimer) {
        clearInterval(reservationTimer);
    }

    const timerElement = document.getElementById('reservationTimer');
    const reservationTime = reservationTimestamp ? new Date(reservationTimestamp) : new Date();
    const expirationTime = new Date(reservationTime.getTime() + 10 * 60 * 1000); // +10 минут
    
    reservationTimer = setInterval(() => {
        const now = new Date();
        const timeLeft = Math.max(0, Math.floor((expirationTime - now) / 1000));
        
        if (timeLeft > 0) {
            const minutes = Math.floor(timeLeft / 60);
            const seconds = timeLeft % 60;
            timerElement.textContent = `Осталось времени для подтверждения: ${minutes}:${seconds.toString().padStart(2, '0')}`;
        } else {
            timerElement.textContent = 'Время истекло. Бронирование будет освобождено автоматически.';
            timerElement.style.background = '#f8d7da';
            clearInterval(reservationTimer);
            reservationTimer = null;
            const expiredEventId = currentReservation ? currentReservation.eventId : null;
            localStorage.removeItem(`reservation_${currentUserId}`);
            currentReservation = null;
            
            setTimeout(() => {
                if (expiredEventId) {
                    viewEventDetails(expiredEventId);
                }
                hideReservationInfo();
            }, 2000);
        }
    }, 1000);
}

// Подтверждение брони
async function confirmReservation() {
    if (!currentReservation || !currentUserId) {
        showMessage('Нет активной брони', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/events/${currentReservation.eventId}/confirm`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                user_id: currentUserId,
                seat_number: currentReservation.seatNumber
            })
        });

        const result = await response.json();

        if (response.ok) {
            showMessage(`Бронирование места ${currentReservation.seatNumber} успешно подтверждено!`, 'success');
            
            // Удаляем бронь из localStorage
            localStorage.removeItem(`reservation_${currentUserId}`);
            
            hideReservationInfo();
            if (reservationTimer) {
                clearInterval(reservationTimer);
                reservationTimer = null;
            }
            
            const eventId = currentReservation.eventId;
            currentReservation = null;
            
            // Обновляем детали мероприятия
            setTimeout(() => {
                viewEventDetails(eventId);
            }, 1000);
        } else {
            showMessage(`Ошибка: ${result.error}`, 'error');
        }
    } catch (error) {
        showMessage(`Ошибка при подтверждении: ${error.message}`, 'error');
    }
}

