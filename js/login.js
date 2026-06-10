// ---- 表单切换 ----
function showLogin() {
    document.getElementById('loginForm').classList.remove('hidden');
    document.getElementById('registerForm').classList.add('hidden');
    document.getElementById('resetForm').classList.add('hidden');
}

function showRegister() {
    document.getElementById('loginForm').classList.add('hidden');
    document.getElementById('registerForm').classList.remove('hidden');
    document.getElementById('resetForm').classList.add('hidden');
}

function showReset() {
    document.getElementById('loginForm').classList.add('hidden');
    document.getElementById('registerForm').classList.add('hidden');
    document.getElementById('resetForm').classList.remove('hidden');
}

// ---- 登录 ----
async function login(event) {
    event.preventDefault();
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;
    if (!username || !password) { alert('请输入账号和密码'); return false; }

    const r = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ username, password })
    });

    if (r.ok) {
        const d = await r.json();
        localStorage.setItem('access_token', d.access_token);
        localStorage.setItem('refresh_token', d.refresh_token);
        localStorage.setItem('isLoggedIn', 'true');
        const redirect = new URLSearchParams(window.location.search).get('redirect') || '/index';
        window.location.href = redirect;
    } else {
        alert('登录失败，请检查用户名和密码');
    }
    return false;
}

// ---- 注册 ----
let codeTimer = null;

function startCountdown(btn) {
    let sec = 60;
    if (codeTimer) clearInterval(codeTimer);
    btn.disabled = true;
    codeTimer = setInterval(() => {
        btn.textContent = sec + '秒';
        sec--;
        if (sec < 0) {
            clearInterval(codeTimer);
            codeTimer = null;
            btn.disabled = false;
            btn.textContent = '获取验证码';
        }
    }, 1000);
}

async function sendVerifyCode() {
    const email = document.getElementById('reg_email').value.trim();
    const btn = document.getElementById('sendCodeBtn');
    if (!email) { alert('请先输入邮箱'); return; }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) { alert('邮箱格式不正确'); return; }

    btn.disabled = true;
    btn.textContent = '发送中...';

    const r = await fetch('/send/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ email })
    });
    const d = await r.json();
    if (r.ok) { alert(d.message); startCountdown(btn); }
    else { alert(d.error || '发送失败'); btn.disabled = false; btn.textContent = '获取验证码'; }
}

async function register(event) {
    event.preventDefault();
    const email = document.getElementById('reg_email').value.trim();
    const code = document.getElementById('reg_code').value.trim();
    const password = document.getElementById('reg_password').value;
    const password2 = document.getElementById('reg_password2').value;

    if (!email || !code || !password) { alert('请填写完整信息'); return false; }
    if (password !== password2) { alert('两次密码输入不一致'); return false; }
    if (password.length < 6) { alert('密码长度不能少于6位'); return false; }

    const r = await fetch('/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ password, email, code })
    });
    const d = await r.json();
    if (r.ok) {
        // 显示成功弹窗
        document.getElementById('modalUsername').textContent = d.username;
        document.getElementById('modalPassword').textContent = password;
        document.getElementById('successModal').classList.remove('hidden');
    } else {
        alert(d.error || '注册失败');
    }
    return false;
}

// 注册成功确认 - 跳转到登录并填充账号
function confirmRegister() {
    const username = document.getElementById('modalUsername').textContent;
    const password = document.getElementById('modalPassword').textContent;
    document.getElementById('successModal').classList.add('hidden');
    showLogin();
    document.getElementById('username').value = username;
    document.getElementById('password').value = password;
}

// ---- 重置密码 ----
async function sendResetCode() {
    const email = document.getElementById('rst_email').value.trim();
    const btn = document.getElementById('sendResetCodeBtn');
    if (!email) { alert('请先输入邮箱'); return; }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) { alert('邮箱格式不正确'); return; }

    btn.disabled = true;
    btn.textContent = '发送中...';

    const r = await fetch('/send/reset-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ email })
    });
    const d = await r.json();
    if (r.ok) { alert(d.message); startCountdown(btn); }
    else { alert(d.error || '发送失败'); btn.disabled = false; btn.textContent = '获取验证码'; }
}

async function resetPassword(event) {
    event.preventDefault();
    const email = document.getElementById('rst_email').value.trim();
    const code = document.getElementById('rst_code').value.trim();
    const password = document.getElementById('rst_password').value;
    const password2 = document.getElementById('rst_password2').value;
    if (!email || !code || !password) { alert('请填写完整信息'); return false; }
    if (password !== password2) { alert('两次密码输入不一致'); return false; }
    if (password.length < 6) { alert('密码长度不能少于6位'); return false; }

    const r = await fetch('/reset', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ email, code, password })
    });
    const d = await r.json();
    if (r.ok) {
        alert('密码重置成功，请重新登录');
        showLogin();
        document.getElementById('username').value = email;
    } else {
        alert(d.error || '重置失败');
    }
    return false;
}
