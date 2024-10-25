document.addEventListener('DOMContentLoaded', (event) => {  
    const dateInput = document.getElementById('notice_date');  
    const today = new Date().toISOString().split('T')[0]; // 获取当前日期，格式为 YYYY-MM-DD  
    dateInput.value = today;  
});

document.getElementById('addnotice').addEventListener('submit', function(event) {
    event.preventDefault(); // 阻止默认表单提交

    if (confirm('确认发布这则公告？')) {
        const formData = new FormData(this);

        // 使用 fetch 提交
        fetch('/add/notice', {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (response.ok) {
                alert('公告发布成功');
                window.location.href = '/add/notice';
            }else if (response.status == 400){
                alert("不能在当前日期之前发布公告")
            }else{
                alert("服务器错误")
                document.getElementById('addnotice').reset();
            }
            return response.json(); // 假设返回 JSON 格式的响应
        })
        .catch(error => {
            console.error('错误:', error); // 用于调试
        });
    }
});