document.addEventListener('DOMContentLoaded', function() {
    const recStateCheckbox = document.getElementById('recStateCheckbox');
    const recTypeInput = document.getElementById('recTypeInput');

    recStateCheckbox.addEventListener('change', function() {
        if (this.checked) {
            recTypeInput.classList.remove('hidden');
        } else {
            recTypeInput.classList.add('hidden');
            recTypeInput.value = ''; // Clear the input when hidden
        }
    });
});

document.getElementById('addbook').addEventListener('submit', function(event) {
    event.preventDefault(); // 阻止默认表单提交

    if (confirm('确定要添加这本图书吗？')) {
        const formData = new FormData(this);

        // 使用 fetch 提交
        fetch('/add/book', {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (!response.ok) {
                // 根据状态码处理错误
                if (response.status === 409) {
                    alert('已存在相同ISBN的图书');
                } else if (response.status === 422) {
                    alert('数据错误，请检查输入');
                } else {
                    alert('服务器错误，请稍后再试');
                }
            }else{
                alert('图书添加成功');
                document.getElementById('addbook').reset();
            }
            return response.json(); // 假设返回 JSON 格式的响应
        })
        .catch(error => {
            console.error('错误:', error); // 用于调试
        });
    }
});
