function apiBase() {
  const v = document.getElementById('apiBase').value.trim();
  return v || 'http://order-server:8080';
}


document.getElementById('saveBase').addEventListener('click', () => {
  localStorage.setItem('orders_api_base', apiBase());
});

document.getElementById('getById').addEventListener('click', async () => {
  const id = document.getElementById('orderId').value.trim();
  const out = document.getElementById('orderResult');
  out.textContent = '';
  if (!id) {
    out.textContent = 'Введите ID заказа';
    return;
  }
  try {
    const res = await fetch(`${apiBase()}/api/order/${encodeURIComponent(id)}`);
    const text = await res.text();
    try {
      out.textContent = JSON.stringify(JSON.parse(text), null, 2);
    } catch {
      out.textContent = text;
    }
  } catch (e) {
    out.textContent = 'Ошибка запроса: ' + e;
  }
});

document.getElementById('getAll').addEventListener('click', async () => {
  const list = document.getElementById('ordersList');
  list.innerHTML = '';
  try {
    const res = await fetch(`${apiBase()}/api/orders`);
    const data = await res.json();
    if (!Array.isArray(data)) {
      list.textContent = 'Неверный формат ответа';
      return;
    }
    for (const o of data) {
      const div = document.createElement('div');
      div.className = 'item';
      div.innerHTML = `
        <div><b>UID:</b> ${escapeHtml(o.order_uid)}</div>
        <div><b>Track:</b> ${escapeHtml(o.track_number)}</div>
        <div><b>Customer:</b> ${escapeHtml(o.customer_id)}</div>
        <div><b>Date:</b> ${o.date_created ? new Date(o.date_created).toLocaleString() : ''}</div>
      `;
      list.appendChild(div);
    }
  } catch (e) {
    list.textContent = 'Ошибка запроса: ' + e;
  }
});

function escapeHtml(s) {
  return String(s ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}


