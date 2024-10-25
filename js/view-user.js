document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchInput');
    const searchButton = document.getElementById('searchButton');
    const userTableBody = document.getElementById('userTableBody');

    // 初始加载所有用户
    fetchUsers('');

    // 搜索按钮点击事件
    searchButton.addEventListener('click', function() {
        fetchUsers(searchInput.value);
    });

    // 回车键搜索
    searchInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            fetchUsers(searchInput.value);
        }
    });

    function fetchUsers(username) {
        fetch('/view/user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `username=${username}`
        })
        .then(response => response.json())
        .then(data => {
            displayUsers(data);
        })
        .catch(error => {
            console.error('Error:', error);
            alert('获取用户数据失败，请重试。');
        });
    }

    function displayUsers(users) {
        userTableBody.innerHTML = '';
        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${user.name}</td>
                <td>${user.username}</td>
                <td>${user.user_cur_lend_amount}</td>
                <td>${user.user_his_lend_amount}</td>
                <td>${user.birthday}</td>
                <td>${user.age}</td>
                <td><button class="reset-password" data-username="${user.username}">重置密码</button></td>
            `;
            userTableBody.appendChild(row);
        });

        // 为所有重置密码按钮添加事件监听器
        document.querySelectorAll('.reset-password').forEach(button => {
            button.addEventListener('click', function() {
                resetPassword(this.getAttribute('data-username'));
            });
        });
    }

    function resetPassword(username) {
        if (confirm(`确定要重置用户 ${username} 的密码吗？`)) {
            fetch('/reset', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `username=${username}`
            })
            .then(response => {
                if (response.ok) {
                    alert('密码已重置为123456');
                } else {
                    throw new Error('密码重置失败');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('密码重置失败，请重试。');
            });
        }
    }
});