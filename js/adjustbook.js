document.addEventListener('DOMContentLoaded', (event) => {  
    const dateInput = document.getElementById('adjust_date');  
    const today = new Date().toISOString().split('T')[0]; // 获取当前日期，格式为 YYYY-MM-DD  
    dateInput.value = today;  
});

document.getElementById('adjustbook').addEventListener('submit', function(event) {
    event.preventDefault(); // 阻止默认表单提交

    if (confirm('确认调整信息无误？')) {
        const formData = new FormData(this);

        // 使用 fetch 提交
        fetch('/adjust/book', {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (response.ok) {
                alert('调整信息发布成功');
                window.location.href = '/adjust/book';
            }else if (response.status == 401){
                alert("当日调整当日发布，不可提前或延后发布")
            }else{
                alert("服务器错误")
                document.getElementById('adjustbook').reset();
            }
        })
        .catch(error => {
            console.error('错误:', error); // 用于调试
        });
    }
});