const API_BASE = ""; // тот же хост, где крутится backend

const selectors = {
    createForm: document.getElementById("create-form"),
    createMessage: document.getElementById("create-message"),
    itemsRefreshBtn: document.getElementById("items-refresh"),
    itemsMessage: document.getElementById("items-message"),
    itemsTableBody: document.querySelector("#items-table tbody"),
    filterFrom: document.getElementById("filter-from"),
    filterTo: document.getElementById("filter-to"),
    updateForm: document.getElementById("update-form"),
    updateMessage: document.getElementById("update-message"),
    updateCancel: document.getElementById("update-cancel"),
    analyticsRefresh: document.getElementById("analytics-refresh"),
    analyticsMessage: document.getElementById("analytics-message"),
    analyticsFrom: document.getElementById("analytics-from"),
    analyticsTo: document.getElementById("analytics-to"),
    analyticsCategory: document.getElementById("analytics-category"),
    metricCount: document.getElementById("metric-count"),
    metricSum: document.getElementById("metric-sum"),
    metricAvg: document.getElementById("metric-avg"),
    metricMedian: document.getElementById("metric-median"),
    metricP90: document.getElementById("metric-p90"),
    chartCanvas: document.getElementById("analytics-chart"),
};

const itemsCache = new Map();
let analyticsChart;

const DEFAULT_RANGE = {
    from: "1970-01-01T00:00:00.000Z",
    to: "2099-12-31T23:59:59.000Z",
};

function safeReset(form) {
    if (form && typeof form.reset === "function") {
        form.reset();
    }
}

const currencyFormatter = new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 2,
});

const numberFormatter = new Intl.NumberFormat("ru-RU", {
    maximumFractionDigits: 2,
});

const dateTimeFormatter = new Intl.DateTimeFormat("ru-RU", {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
});

function setMessage(target, text, type = "info") {
    if (!target) return;
    target.textContent = text ?? "";
    target.classList.remove("error", "success");
    if (type === "error") {
        target.classList.add("error");
    } else if (type === "success") {
        target.classList.add("success");
    }
}

function toISOStringOrNull(value) {
    if (!value) return null;
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return null;
    return date.toISOString();
}

function isoToDatetimeLocal(value) {
    if (!value) return "";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "";
    const offset = date.getTimezoneOffset();
    const local = new Date(date.getTime() - offset * 60000);
    return local.toISOString().slice(0, 16);
}

function normalizeDate(value, { endOfDay = false } = {}) {
    if (!value) return null;

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return null;
    }

    if (endOfDay && value.length <= 10) {
        date.setHours(23, 59, 59, 999);
    }

    return date.toISOString();
}

function buildRangeSearchParams(fromValue, toValue, { fallback = true } = {}) {
    const params = new URLSearchParams();
    const fromISO = normalizeDate(fromValue, { endOfDay: false });
    const toISO = normalizeDate(toValue, { endOfDay: true });

    if (fromISO) {
        params.set("from", fromISO);
    } else if (fallback) {
        params.set("from", DEFAULT_RANGE.from);
    }

    if (toISO) {
        params.set("to", toISO);
    } else if (fallback) {
        params.set("to", DEFAULT_RANGE.to);
    }

    return params;
}

async function request(url, options = {}) {
    try {
        const response = await fetch(`${API_BASE}${url}`, {
            headers: {
                "Content-Type": "application/json",
                ...(options.headers ?? {}),
            },
            ...options,
        });

        if (!response.ok) {
            const errorBody = await response.json().catch(() => ({}));
            const message = errorBody.error || response.statusText || "Неизвестная ошибка";
            throw new Error(message);
        }

        if (response.status === 204) {
            return null;
        }

        const contentType = response.headers.get("content-type") ?? "";
        if (contentType.includes("application/json")) {
            return response.json();
        }

        return response.text();
    } catch (error) {
        throw new Error(error.message || "Не удалось выполнить запрос");
    }
}

async function fetchItems() {
    const params = buildRangeSearchParams(
        selectors.filterFrom.value,
        selectors.filterTo.value,
        { fallback: true }
    );

    const url = `/items?${params.toString()}`;

    setMessage(selectors.itemsMessage, "Загружаем данные...");
    try {
        const response = await request(url, { method: "GET" });
        const list = Array.isArray(response) ? response : [];
        renderItems(list);
        setMessage(selectors.itemsMessage, `Найдено записей: ${list.length}`, "success");
    } catch (error) {
        renderItems([]);
        setMessage(selectors.itemsMessage, error.message, "error");
    }
}

function renderItems(items) {
    selectors.itemsTableBody.innerHTML = "";
    itemsCache.clear();

    if (!items.length) {
        const emptyRow = document.createElement("tr");
        emptyRow.innerHTML = `<td colspan="5">Записей не найдено</td>`;
        selectors.itemsTableBody.append(emptyRow);
        return;
    }

    items.forEach((item) => {
        itemsCache.set(item.id, item);
        const row = document.createElement("tr");
        row.dataset.itemId = item.id;
        row.innerHTML = `
            <td>${item.title}</td>
            <td>${item.category}</td>
            <td class="text-right">${currencyFormatter.format(item.amount)}</td>
            <td>${dateTimeFormatter.format(new Date(item.date))}</td>
            <td class="actions">
                <div class="action-group">
                    <button type="button" data-action="edit">Редактировать</button>
                    <button type="button" data-action="delete">Удалить</button>
                </div>
            </td>
        `;

        selectors.itemsTableBody.append(row);
    });
}

async function handleCreateSubmit(event) {
    event.preventDefault();

    const form = event.currentTarget;
    const formData = new FormData(form);
    const payload = {
        title: formData.get("title")?.trim(),
        category: formData.get("category")?.trim(),
        amount: Number(formData.get("amount")),
        date: toISOStringOrNull(formData.get("date")),
    };

    if (!payload.date) {
        setMessage(selectors.createMessage, "Заполните корректную дату", "error");
        return;
    }

    setMessage(selectors.createMessage, "Отправляем...");
    try {
        await request("/items", {
            method: "POST",
            body: JSON.stringify(payload),
        });
        setMessage(selectors.createMessage, "Запись добавлена", "success");
        safeReset(form);
        await fetchItems();
    } catch (error) {
        setMessage(selectors.createMessage, error.message, "error");
    }
}

function handleTableClick(event) {
    const button = event.target.closest("button[data-action]");
    if (!button) return;

    const row = button.closest("tr[data-item-id]");
    if (!row) return;

    const id = row.dataset.itemId;

    if (button.dataset.action === "edit") {
        openUpdateForm(id);
        return;
    }

    if (button.dataset.action === "delete") {
        deleteItem(id);
    }
}

function openUpdateForm(id) {
    const item = itemsCache.get(id);
    if (!item) {
        setMessage(selectors.itemsMessage, "Не удалось найти запись для редактирования", "error");
        return;
    }

    selectors.updateForm.classList.remove("hidden");
    selectors.updateForm.elements.id.value = item.id;
    selectors.updateForm.elements.title.value = item.title ?? "";
    selectors.updateForm.elements.category.value = item.category ?? "";
    selectors.updateForm.elements.amount.value = item.amount ?? 0;
    selectors.updateForm.elements.date.value = isoToDatetimeLocal(item.date);
    setMessage(selectors.updateMessage, "Вы редактируете запись", "info");
}

async function deleteItem(id) {
    if (!window.confirm("Удалить запись?")) return;
    setMessage(selectors.itemsMessage, "Удаляем...");

    try {
        await request(`/items/${id}`, { method: "DELETE" });
        setMessage(selectors.itemsMessage, "Запись удалена", "success");
        await fetchItems();
    } catch (error) {
        setMessage(selectors.itemsMessage, error.message, "error");
    }
}

async function handleUpdateSubmit(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const formData = new FormData(form);

    const payload = {
        title: formData.get("title")?.trim(),
        category: formData.get("category")?.trim(),
        amount: Number(formData.get("amount")),
        date: toISOStringOrNull(formData.get("date")),
    };

    const id = formData.get("id");
    if (!id) {
        setMessage(selectors.updateMessage, "Не удалось определить запись", "error");
        return;
    }

    if (!payload.date) {
        setMessage(selectors.updateMessage, "Проверьте дату", "error");
        return;
    }

    setMessage(selectors.updateMessage, "Сохраняем изменения...");
    try {
        await request(`/items/${id}`, {
            method: "PUT",
            body: JSON.stringify(payload),
        });
        setMessage(selectors.updateMessage, "Запись обновлена", "success");
        selectors.updateForm.classList.add("hidden");
        safeReset(form);
        await fetchItems();
    } catch (error) {
        setMessage(selectors.updateMessage, error.message, "error");
    }
}

function handleCancelUpdate() {
    selectors.updateForm.classList.add("hidden");
    safeReset(selectors.updateForm);
    setMessage(selectors.updateMessage, "");
}

async function fetchAnalytics() {
    const params = buildRangeSearchParams(
        selectors.analyticsFrom.value,
        selectors.analyticsTo.value,
        { fallback: true }
    );
    if (selectors.analyticsCategory.value.trim()) {
        params.set("category", selectors.analyticsCategory.value.trim());
    }

    const url = `/analytics?${params.toString()}`;

    setMessage(selectors.analyticsMessage, "Считаем аналитику...");
    try {
        const response = await request(url, { method: "GET" });
        const analytics = response?.analytics ?? null;
        updateAnalyticsView(analytics);
        setMessage(selectors.analyticsMessage, "Готово", "success");
    } catch (error) {
        updateAnalyticsView(null);
        setMessage(selectors.analyticsMessage, error.message, "error");
    }
}

function updateAnalyticsView(data) {
    if (!data) {
        selectors.metricCount.textContent = "—";
        selectors.metricSum.textContent = "—";
        selectors.metricAvg.textContent = "—";
        selectors.metricMedian.textContent = "—";
        selectors.metricP90.textContent = "—";
        renderAnalyticsChart(null);
        return;
    }

    const totalCount = data.totalCount ?? data.TotalCount ?? 0;
    const totalSum = data.totalSum ?? data.TotalSum ?? 0;
    const avgAmount = data.avgAmount ?? data.AvgAmount ?? 0;
    const median = data.median ?? data.Median ?? 0;
    const p90 = data.p90 ?? data.P90 ?? 0;

    selectors.metricCount.textContent = numberFormatter.format(totalCount);
    selectors.metricSum.textContent = currencyFormatter.format(totalSum);
    selectors.metricAvg.textContent = currencyFormatter.format(avgAmount);
    selectors.metricMedian.textContent = currencyFormatter.format(median);
    selectors.metricP90.textContent = currencyFormatter.format(p90);

    renderAnalyticsChart({
        totalSum,
        avgAmount,
        median,
        p90,
    });
}

function renderAnalyticsChart(data) {
    if (!selectors.chartCanvas) return;

    if (!data) {
        if (analyticsChart) {
            analyticsChart.destroy();
            analyticsChart = null;
        }
        return;
    }

    const datasets = [
        data.totalSum ?? 0,
        data.avgAmount ?? 0,
        data.median ?? 0,
        data.p90 ?? 0,
    ];

    const labels = ["Сумма", "Среднее", "Медиана", "P90"];

    if (analyticsChart) {
        analyticsChart.data.labels = labels;
        analyticsChart.data.datasets[0].data = datasets;
        analyticsChart.update();
        return;
    }

    analyticsChart = new Chart(selectors.chartCanvas, {
        type: "bar",
        data: {
            labels,
            datasets: [
                {
                    label: "Рубли",
                    data: datasets,
                    backgroundColor: "rgba(58, 111, 247, 0.6)",
                    borderColor: "rgba(58, 111, 247, 0.9)",
                    borderWidth: 1,
                },
            ],
        },
        options: {
            responsive: true,
            plugins: {
                legend: { display: false },
                tooltip: {
                    callbacks: {
                        label(context) {
                            return currencyFormatter.format(context.parsed.y ?? 0);
                        },
                    },
                },
            },
            scales: {
                y: {
                    ticks: {
                        callback(value) {
                            return currencyFormatter.format(value ?? 0);
                        },
                    },
                },
            },
        },
    });
}

function initEventListeners() {
    selectors.createForm?.addEventListener("submit", handleCreateSubmit);
    selectors.itemsRefreshBtn?.addEventListener("click", fetchItems);
    selectors.itemsTableBody?.addEventListener("click", handleTableClick);
    selectors.updateForm?.addEventListener("submit", handleUpdateSubmit);
    selectors.updateCancel?.addEventListener("click", handleCancelUpdate);
    selectors.analyticsRefresh?.addEventListener("click", fetchAnalytics);
}

async function bootstrap() {
    initEventListeners();
    await fetchItems();
    await fetchAnalytics();
}

document.addEventListener("DOMContentLoaded", bootstrap);

