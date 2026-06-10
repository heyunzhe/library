document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('returnRecordsForm');
    const queryType = document.getElementById('queryType');
    const queryInput = document.getElementById('queryInput');
    const table = document.getElementById('returnRecordsTable');

    queryType.addEventListener('change', function() {
        if (this.value === '2') {
            queryInput.type = 'date';
            queryInput.placeholder = '请选择借阅日期';
        } else {
            queryInput.type = 'text';
            queryInput.placeholder = this.value === '1' ? '' : '请输入用户账号';
        }
    });

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        
        fetch('/return/records', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            updateTable(data);
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
            // alert('获取数据失败，请重试');
        });
    });

    function updateTable(data) {
        const tbody = table.querySelector('tbody');
        tbody.innerHTML = ''; // Clear existing rows

        if (data.length === 0) {
            const row = tbody.insertRow();
            const cell = row.insertCell();
            cell.colSpan = 8;
            cell.textContent = '没有找到相关记录';
            return;
        }

        data.forEach(record => {
            const row = tbody.insertRow();
            row.insertCell().textContent = record.return_id;
            row.insertCell().textContent = record.username;
            row.insertCell().textContent = record.title;
            row.insertCell().textContent = record.isbn;
            row.insertCell().textContent = record.lend_date;
            row.insertCell().textContent = record.exp_return_date;
            row.insertCell().textContent = record.return_date;
            row.insertCell().textContent = record.late_fee;
        });
    }
    form.dispatchEvent(new Event('submit'));
});