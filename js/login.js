const loginForm = document.getElementById('loginForm');
const registerForm = document.getElementById('registerForm');
const toggleBtn = document.getElementById('toggleForm');


        toggleBtn.addEventListener('click', function() {
            if (loginForm.classList.contains('hidden')) {
                loginForm.classList.remove('hidden');
                registerForm.classList.add('hidden');
                this.textContent = '切换到注册';
            } else {
                loginForm.classList.add('hidden');
                registerForm.classList.remove('hidden');
                this.textContent = '切换到登录';
            }
        });


async function login(event) {
    event.preventDefault(); // 阻止默认的表单提交

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const response = await fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            username: username,
            password: password
        })
    });

    if (response.ok) {
        localStorage.setItem('isLoggedIn', 'true');
        window.location.href = '/index'; // 登录成功后重定向
    } else {
        alert('登录失败，请检查用户名和密码');
    }
}




