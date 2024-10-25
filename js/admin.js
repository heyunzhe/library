var modal = document.getElementById('adminModal');

        // 获取打开模态窗口的链接
        var adminLink = document.getElementById('adminLink');

        // 获取关闭按钮
        var span = document.getElementsByClassName('close')[0];

        // 获取表单和密码输入框
        var adminForm = document.getElementById('adminForm');
        var adminPassword = document.getElementById('adminPassword');

        // 当用户点击链接时，打开模态窗口
        adminLink.onclick = function(event) {
            event.preventDefault();
            modal.style.display = 'block';
        }

        // 当用户点击 (x) 时，关闭模态窗口
        span.onclick = function() {
            modal.style.display = 'none';
        }

        // 当用户点击模态窗口外部时，关闭它
        window.onclick = function(event) {
            if (event.target == modal) {
                modal.style.display = 'none';
            }
        }


        adminForm.onsubmit = function(event) {
            event.preventDefault();
            var formData = new FormData(adminForm);

            fetch('/admin', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    document.getElementById('adminPassword').value = '';
                    window.location.href = '/admin'; // 重定向到管理页面
                } else if (response.status === 401) {
                    alert('密码错误或已有管理员登录');
                    document.getElementById('adminPassword').value = ''; // 清空密码输入框
                } else {
                    throw new Error('服务器错误');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('登录失败，请稍后重试。');
            });
        }