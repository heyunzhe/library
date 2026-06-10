document.addEventListener('DOMContentLoaded', function() {
    var modal = document.getElementById('adminModal');
    var adminLink = document.getElementById('adminLink');
    var span = document.getElementsByClassName('close')[0];
    var adminForm = document.getElementById('adminForm');

    // 当用户点击链接时，打开模态窗口
    adminLink.onclick = function(event) {
        event.preventDefault();
        modal.style.display = 'block';
    };

    // 当用户点击 (x) 时，关闭模态窗口
    span.onclick = function() {
        modal.style.display = 'none';
    };

    // 当用户点击模态窗口外部时，关闭它
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = 'none';
        }
    };

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
                document.getElementById('adminModal').style.display = 'none';
                window.location.href = '/admin';
            } else if (response.status === 401) {
                alert('帐号或密码错误');
                document.getElementById('adminPassword').value = '';
            } else {
                throw new Error('服务器错误');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('登录失败，请稍后重试。');
        });
    };
});