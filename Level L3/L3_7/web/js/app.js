const API_BASE_URL = 'http://localhost:8081';
let currentToken = null;
let currentUser = null;
let currentRole = null;
let editingProductId = null;
let currentProductIdForLogs = null; // Для сохранения CSV

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    // Проверяем, есть ли сохраненный токен
    const savedToken = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    const savedRole = localStorage.getItem('role');
    
    if (savedToken && savedUser && savedRole) {
        currentToken = savedToken;
        currentUser = savedUser;
        currentRole = savedRole;
        showMainPanel();
        loadProducts();
    }
});

// Переключение между вкладками входа и регистрации
function showLoginTab() {
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('registerForm').style.display = 'none';
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    event.target.classList.add('active');
}

function showRegisterTab() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'block';
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    event.target.classList.add('active');
}

// Регистрация нового пользователя
async function register() {
    const username = document.getElementById('regUsername').value.trim();
    const role = document.getElementById('regRole').value;
    const errorDiv = document.getElementById('registerError');
    const successDiv = document.getElementById('registerSuccess');

    // Скрываем предыдущие сообщения
    errorDiv.classList.remove('show');
    successDiv.classList.remove('show');

    if (!username) {
        showError(errorDiv, 'Введите имя пользователя');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/user/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                role: role
            })
        });

        const data = await response.json();

        if (!response.ok) {
            showError(errorDiv, data.error || 'Ошибка регистрации');
            return;
        }

        // Показываем успешное сообщение
        successDiv.textContent = `Пользователь "${username}" успешно зарегистрирован! Теперь вы можете войти.`;
        successDiv.classList.add('show');

        // Очищаем форму
        document.getElementById('regUsername').value = '';

        // Автоматически переключаемся на вкладку входа через 2 секунды
        setTimeout(() => {
            showLoginTab();
            document.getElementById('username').value = username;
            document.getElementById('userRole').value = role;
        }, 2000);
    } catch (error) {
        showError(errorDiv, 'Ошибка подключения к серверу: ' + error.message);
    }
}

// Вход в систему
async function login() {
    const username = document.getElementById('username').value.trim();
    const role = document.getElementById('userRole').value;
    const errorDiv = document.getElementById('loginError');

    if (!username) {
        showError(errorDiv, 'Введите имя пользователя');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/user/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                role: role
            })
        });

        const data = await response.json();

        if (!response.ok) {
            showError(errorDiv, data.error || 'Ошибка входа');
            return;
        }

        currentToken = data.jwt_token;
        currentUser = username;
        currentRole = role;

        // Сохраняем в localStorage
        localStorage.setItem('token', currentToken);
        localStorage.setItem('user', currentUser);
        localStorage.setItem('role', currentRole);

        showMainPanel();
        loadProducts();
    } catch (error) {
        showError(errorDiv, 'Ошибка подключения к серверу: ' + error.message);
    }
}

// Выход из системы
function logout() {
    currentToken = null;
    currentUser = null;
    currentRole = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('role');
    
    document.getElementById('authPanel').style.display = 'block';
    document.getElementById('mainPanel').style.display = 'none';
    document.getElementById('username').value = '';
    showLoginTab();
}

// Показать главную панель
function showMainPanel() {
    document.getElementById('authPanel').style.display = 'none';
    document.getElementById('mainPanel').style.display = 'block';
    document.getElementById('currentUser').textContent = currentUser;
    document.getElementById('currentRole').textContent = currentRole;

    // Показываем кнопку добавления только для admin
    const addSection = document.getElementById('addProductSection');
    if (currentRole === 'admin') {
        addSection.style.display = 'block';
    } else {
        addSection.style.display = 'none';
    }
}

// Поиск товара по имени
async function searchProduct() {
    const productName = document.getElementById('searchProductName').value.trim();
    const searchResult = document.getElementById('searchResult');

    if (!productName) {
        alert('Введите название товара для поиска');
        return;
    }

    searchResult.style.display = 'block';
    searchResult.innerHTML = '<div class="loading">Поиск...</div>';

    try {
        const encodedName = encodeURIComponent(productName);
        const response = await fetch(`${API_BASE_URL}/api/v1/items/${encodedName}`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                alert('Сессия истекла. Пожалуйста, войдите снова.');
                logout();
                return;
            }
            if (response.status === 404) {
                searchResult.innerHTML = '<div class="error-message show">Товар не найден</div>';
                return;
            }
            throw new Error('Ошибка поиска товара');
        }

        const data = await response.json();
        const product = data.product;

        searchResult.innerHTML = `
            <div class="search-result-card">
                <h3>${escapeHtml(product.name)}</h3>
                <p><strong>ID:</strong> ${product.id}</p>
                <p><strong>Описание:</strong> ${escapeHtml(product.description || 'Нет описания')}</p>
                <p><strong>Количество:</strong> ${product.quantity}</p>
                <p><strong>Обновлено:</strong> ${formatDate(product.updated_at)}</p>
                <div style="margin-top: 10px;">
                    <button onclick="viewLogs('${product.id}')" class="btn btn-warning btn-sm">История</button>
                    ${(currentRole === 'admin' || currentRole === 'manager') ? 
                        `<button onclick="editProduct('${product.id}', '${escapeHtml(product.name)}', '${escapeHtml(product.description || '')}', ${product.quantity})" class="btn btn-primary btn-sm">Редактировать</button>` : ''}
                    ${currentRole === 'admin' ? 
                        `<button onclick="deleteProduct('${escapeHtml(product.name)}')" class="btn btn-danger btn-sm">Удалить</button>` : ''}
                </div>
            </div>
        `;
    } catch (error) {
        searchResult.innerHTML = `<div class="error-message show">Ошибка: ${error.message}</div>`;
    }
}

// Загрузить список товаров
async function loadProducts() {
    const tbody = document.getElementById('productsTableBody');
    tbody.innerHTML = '<tr><td colspan="6" class="loading">Загрузка...</td></tr>';

    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/items`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                alert('Сессия истекла. Пожалуйста, войдите снова.');
                logout();
                return;
            }
            throw new Error('Ошибка загрузки товаров');
        }

        const data = await response.json();
        const products = data.products || [];

        if (products.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="empty-state">Товары не найдены</td></tr>';
            return;
        }

        tbody.innerHTML = products.map(product => {
            const canEdit = currentRole === 'admin' || currentRole === 'manager';
            const canDelete = currentRole === 'admin';
            
            return `
                <tr>
                    <td>${product.id}</td>
                    <td>${escapeHtml(product.name)}</td>
                    <td>${escapeHtml(product.description || '')}</td>
                    <td>${product.quantity}</td>
                    <td>${formatDate(product.updated_at)}</td>
                    <td class="action-buttons">
                        <button onclick="viewLogs('${product.id}')" class="btn btn-warning btn-sm">История</button>
                        ${canEdit ? `<button onclick="editProduct('${product.id}', '${escapeHtml(product.name)}', '${escapeHtml(product.description || '')}', ${product.quantity})" class="btn btn-primary btn-sm">Редактировать</button>` : ''}
                        ${canDelete ? `<button onclick="deleteProduct('${escapeHtml(product.name)}')" class="btn btn-danger btn-sm">Удалить</button>` : ''}
                    </td>
                </tr>
            `;
        }).join('');
    } catch (error) {
        tbody.innerHTML = `<tr><td colspan="6" class="error-message show">Ошибка: ${error.message}</td></tr>`;
    }
}

// Показать форму добавления товара
function showAddProductForm() {
    editingProductId = null;
    document.getElementById('formTitle').textContent = 'Добавить товар';
    document.getElementById('productName').value = '';
    document.getElementById('productDescription').value = '';
    document.getElementById('productQuantity').value = '';
    document.getElementById('productForm').style.display = 'block';
    document.getElementById('productName').disabled = false;
}

// Редактировать товар
function editProduct(id, name, description, quantity) {
    editingProductId = id;
    document.getElementById('formTitle').textContent = 'Редактировать товар';
    document.getElementById('productName').value = name;
    document.getElementById('productDescription').value = description;
    document.getElementById('productQuantity').value = quantity;
    document.getElementById('productForm').style.display = 'block';
    document.getElementById('productName').disabled = true; // Нельзя менять имя при редактировании
}

// Сохранить товар
async function saveProduct() {
    const name = document.getElementById('productName').value.trim();
    const description = document.getElementById('productDescription').value.trim();
    const quantity = parseInt(document.getElementById('productQuantity').value);

    if (!name) {
        alert('Введите название товара');
        return;
    }

    if (isNaN(quantity) || quantity < 0) {
        alert('Введите корректное количество');
        return;
    }

    try {
        const productData = {
            name: name,
            description: description,
            quantity: quantity
        };

        if (editingProductId) {
            productData.id = editingProductId;
        }

        const url = `${API_BASE_URL}/api/v1/items`;
        const method = editingProductId ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify(productData)
        });

        if (!response.ok) {
            const data = await response.json();
            if (response.status === 401) {
                alert('Сессия истекла. Пожалуйста, войдите снова.');
                logout();
                return;
            }
            if (response.status === 403) {
                alert('У вас нет прав для выполнения этого действия');
                return;
            }
            throw new Error(data.error || 'Ошибка сохранения товара');
        }

        cancelProductForm();
        loadProducts();
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
}

// Отменить форму
function cancelProductForm() {
    document.getElementById('productForm').style.display = 'none';
    editingProductId = null;
}

// Удалить товар
async function deleteProduct(productName) {
    if (!confirm(`Вы уверены, что хотите удалить товар "${productName}"?`)) {
        return;
    }

    try {
        const encodedName = encodeURIComponent(productName);
        const response = await fetch(`${API_BASE_URL}/api/v1/items/${encodedName}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (!response.ok) {
            const data = await response.json();
            if (response.status === 401) {
                alert('Сессия истекла. Пожалуйста, войдите снова.');
                logout();
                return;
            }
            if (response.status === 403) {
                alert('У вас нет прав для удаления товаров');
                return;
            }
            throw new Error(data.error || 'Ошибка удаления товара');
        }

        loadProducts();
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
}

// Просмотр истории изменений
async function viewLogs(productId) {
    currentProductIdForLogs = productId; // Сохраняем ID для сохранения CSV
    const modal = document.getElementById('logsModal');
    const tbody = document.getElementById('logsTableBody');
    tbody.innerHTML = '<tr><td colspan="4" class="loading">Загрузка...</td></tr>';
    modal.style.display = 'block';

    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/product/logs/${productId}`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                alert('Сессия истекла. Пожалуйста, войдите снова.');
                logout();
                closeLogsModal();
                return;
            }
            throw new Error('Ошибка загрузки истории');
        }

        const data = await response.json();
        const logs = data.logs || [];

        if (logs.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="empty-state">История изменений отсутствует</td></tr>';
            return;
        }

        tbody.innerHTML = logs.map(log => {
            return `
                <tr>
                    <td>${formatDate(log.changed_at)}</td>
                    <td>${escapeHtml(log.old_name || '')} → ${escapeHtml(log.new_name || '')}</td>
                    <td>${escapeHtml(log.old_description || '')} → ${escapeHtml(log.new_description || '')}</td>
                    <td>${log.old_quantity} → ${log.new_quantity}</td>
                </tr>
            `;
        }).join('');
    } catch (error) {
        tbody.innerHTML = `<tr><td colspan="4" class="error-message show">Ошибка: ${error.message}</td></tr>`;
    }
}

// Сохранить логи в CSV
async function saveLogsToCSV() {
    if (!currentProductIdForLogs) {
        alert('Нет данных для сохранения');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/product/save/${currentProductIdForLogs}`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                alert('Сессия истекла. Пожалуйста, войдите снова.');
                logout();
                return;
            }
            const errorText = await response.text();
            let errorMessage = 'Ошибка сохранения CSV';
            try {
                const errorData = JSON.parse(errorText);
                errorMessage = errorData.error || errorMessage;
            } catch (e) {
                errorMessage = errorText || errorMessage;
            }
            throw new Error(errorMessage);
        }

        // Получаем имя файла из заголовка Content-Disposition или создаем свое
        const contentDisposition = response.headers.get('Content-Disposition');
        let fileName = `${currentProductIdForLogs}_history.csv`;
        
        if (contentDisposition) {
            const fileNameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
            if (fileNameMatch && fileNameMatch[1]) {
                fileName = fileNameMatch[1].replace(/['"]/g, '');
            }
        }

        // Получаем данные как blob
        const blob = await response.blob();
        
        // Создаем ссылку для скачивания
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        
        // Очищаем
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
}

// Закрыть модальное окно истории
function closeLogsModal() {
    document.getElementById('logsModal').style.display = 'none';
    currentProductIdForLogs = null;
}

// Закрыть модальное окно при клике вне его
window.onclick = function(event) {
    const modal = document.getElementById('logsModal');
    if (event.target === modal) {
        closeLogsModal();
    }
}

// Вспомогательные функции
function showError(element, message) {
    element.textContent = message;
    element.classList.add('show');
    setTimeout(() => {
        element.classList.remove('show');
    }, 5000);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}