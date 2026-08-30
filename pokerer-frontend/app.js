const API = "http://localhost:8080";

let currentUser = null;
let wallet = null;
let currentTable = null;

async function api(path, options = {}) {
  const token = localStorage.getItem("pokerer_token");

  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {})
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(API + path, {
    ...options,
    headers
  });

  if (!response.ok) {
    let message = `Request failed: ${response.status}`;

    try {
      const body = await response.json();
      if (body.error) message = body.error;
    } catch {}

    throw new Error(message);
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
}

function money(value) {
  return "$" + Number(value).toLocaleString();
}

function renderNav() {
  const nav = document.getElementById("nav");

  if (currentUser) {
    nav.innerHTML = `
      <span class="balance">
        Balance <strong>${money(wallet?.balance ?? 0)}</strong>
      </span>
      <span>${currentUser.username || currentUser.email}</span>
      <button class="secondary" onclick="logout()">Logout</button>
    `;
  } else {
    nav.innerHTML = `
      <button class="secondary" onclick="showLogin()">Login</button>
      <button onclick="showRegister()">Register</button>
    `;
  }
}

async function loadSession() {
  if (!localStorage.getItem("pokerer_token")) {
    currentUser = null;
    wallet = null;
    renderNav();
    return;
  }

  try {
    currentUser = await api("/api/v1/me");
    wallet = await api("/api/v1/wallet");
  } catch {
    localStorage.removeItem("pokerer_token");
    currentUser = null;
    wallet = null;
  }

  renderNav();
}

async function showLobby() {
  currentTable = null;

  try {
    const tables = await api("/api/v1/tables");

    document.getElementById("app").innerHTML = `
      <section class="hero">
        <h1>Poker tables</h1>
        <p>Choose a table and take your seat.</p>
      </section>

      <section class="grid">
        ${tables.length
          ? tables.map(tableCard).join("")
          : `<p>No tables available.</p>`
        }
      </section>
    `;
  } catch (error) {
    showError(error);
  }
}

function tableCard(table) {
  return `
    <article class="card">
      <h2>${escapeHtml(table.name)}</h2>

      <div class="stakes">
        ${money(table.small_blind)}
        <span>/</span>
        ${money(table.big_blind)}
      </div>

      <div class="info">
        <span>Buy-in</span>
        <strong>${money(table.min_buy_in)} - ${money(table.max_buy_in)}</strong>
      </div>

      <div class="info">
        <span>Players</span>
        <strong>Up to ${table.max_players}</strong>
      </div>

      <button class="full" onclick="openTable('${table.id}')">
        View table
      </button>
    </article>
  `;
}

async function openTable(id) {
  try {
    currentTable = await api(`/api/v1/tables/${id}`);
    renderTable();
  } catch (error) {
    showError(error);
  }
}

function renderTable() {
  const table = currentTable;

  document.getElementById("app").innerHTML = `
    <div class="table-header">
      <div>
        <button class="secondary" onclick="showLobby()">← Lobby</button>
        <h1>${escapeHtml(table.name)}</h1>
        <p>
          ${money(table.small_blind)} / ${money(table.big_blind)}
          · ${money(table.min_buy_in)} - ${money(table.max_buy_in)}
        </p>
      </div>

      ${
        currentUser &&
        table.players.some(p => p.user_id === currentUser.id)
          ? `<button class="danger" onclick="leaveTable('${table.id}')">
               Leave table
             </button>`
          : ""
      }
    </div>

    <div class="table-layout">
      <section class="poker-table">
        <div class="community">
          <div class="community-card">?</div>
          <div class="community-card">?</div>
          <div class="community-card">?</div>
          <div class="community-card">?</div>
          <div class="community-card">?</div>
        </div>

        <div class="pot">
          <small>Pot</small>
          <strong>—</strong>
        </div>

        ${table.players.map(player => `
          <div class="seat seat-${player.seat_number}">
            <div class="player">${escapeHtml(player.user_id)}</div>
            <div class="chips">${money(player.chips)}</div>
          </div>
        `).join("")}
      </section>

      <aside class="side-panel">
        ${renderJoinPanel(table)}

        <hr>

        <h3>Players</h3>

        ${table.players.length
          ? table.players.map(player => `
              <div class="player-row">
                <span>Seat ${player.seat_number}</span>
                <strong>${escapeHtml(player.user_id)}</strong>
                <span>${money(player.chips)}</span>
              </div>
            `).join("")
          : "<p>No players yet.</p>"
        }
      </aside>
    </div>
  `;
}

function renderJoinPanel(table) {
  if (!currentUser) {
    return `
      <h2>Join table</h2>
      <p>You need to log in first.</p>
      <button class="full" onclick="showLogin()">Login</button>
    `;
  }

  const seated = table.players.some(p => p.user_id === currentUser.id);

  if (seated) {
    return `
      <h2>You are seated</h2>
      <p>Your chips are currently at this table.</p>
    `;
  }

  return `
    <h2>Join table</h2>

    <label>
      Buy-in
      <input
        id="buyIn"
        type="number"
        min="${table.min_buy_in}"
        max="${table.max_buy_in}"
        value="${table.min_buy_in}"
      >
    </label>

    <button class="full" onclick="joinTable('${table.id}')">
      Sit down
    </button>
  `;
}

async function joinTable(id) {
  const buyIn = Number(document.getElementById("buyIn").value);

  try {
    await api(`/api/v1/tables/${id}/join`, {
      method: "POST",
      body: JSON.stringify({
        buy_in: buyIn
      })
    });

    wallet = await api("/api/v1/wallet");
    currentTable = await api(`/api/v1/tables/${id}`);

    renderNav();
    renderTable();
  } catch (error) {
    showError(error);
  }
}

async function leaveTable(id) {
  try {
    await api(`/api/v1/tables/${id}/leave`, {
      method: "POST"
    });

    wallet = await api("/api/v1/wallet");
    renderNav();

    showLobby();
  } catch (error) {
    showError(error);
  }
}

function showLogin() {
  document.getElementById("app").innerHTML = `
    <form class="form-card" onsubmit="login(event)">
      <h1>Login</h1>

      <label>
        Email
        <input id="email" type="email" required>
      </label>

      <label>
        Password
        <input id="password" type="password" required>
      </label>

      <button class="full">Login</button>

      <p>
        Don't have an account?
        <a href="#" onclick="showRegister()">Register</a>
      </p>
    </form>
  `;
}

async function login(event) {
  event.preventDefault();

  try {
    const response = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({
        email: document.getElementById("email").value,
        password: document.getElementById("password").value
      })
    });

    localStorage.setItem("pokerer_token", response.access_token);

    await loadSession();
    showLobby();
  } catch (error) {
    showError(error);
  }
}

function showRegister() {
  document.getElementById("app").innerHTML = `
    <form class="form-card" onsubmit="register(event)">
      <h1>Register</h1>

      <label>
        Username
        <input id="username" required>
      </label>

      <label>
        Email
        <input id="email" type="email" required>
      </label>

      <label>
        Password
        <input id="password" type="password" required>
      </label>

      <button class="full">Create account</button>

      <p>
        Already have an account?
        <a href="#" onclick="showLogin()">Login</a>
      </p>
    </form>
  `;
}

async function register(event) {
  event.preventDefault();

  try {
    await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({
        username: document.getElementById("username").value,
        email: document.getElementById("email").value,
        password: document.getElementById("password").value
      })
    });

    showLogin();
  } catch (error) {
    showError(error);
  }
}

function logout() {
  localStorage.removeItem("pokerer_token");
  currentUser = null;
  wallet = null;
  renderNav();
  showLobby();
}

function showError(error) {
  const message = error instanceof Error ? error.message : String(error);

  const old = document.querySelector(".error");
  if (old) old.remove();

  const element = document.createElement("div");
  element.className = "error";
  element.textContent = message;

  document.body.appendChild(element);

  setTimeout(() => element.remove(), 4000);
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

(async function init() {
  await loadSession();
  await showLobby();
})();
