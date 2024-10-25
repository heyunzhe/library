document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('searchForm');
    const noticesList = document.getElementById('noticesList');
    const editModal = document.getElementById('editModal');
    const editForm = document.getElementById('editForm');
    const closeBtn = document.getElementsByClassName('close')[0];

    searchForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const formData = new FormData(searchForm);

        fetch('/view/notice', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('HTTP error ' + response.status);
            }
            return response.json();
        })
        .then(data => {
            displayNotices(data);
        })
        .catch(error => {
            console.error('Error:', error);
            alert('查询失败，请重试。');
        });
    });

    function displayNotices(notices) {
        noticesList.innerHTML = '';
        notices.forEach(notice => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${notice.notice_id}</td>
                <td>${notice.notice_date}</td>
                <td>${notice.notice_title}</td>
                <td>${notice.notice}</td>
                <td>
                    <button class="edit-btn" data-id="${notice.notice_id}">编辑</button>
                    <button class="delete-btn" data-id="${notice.notice_id}">删除</button>
                </td>
            `;
            noticesList.appendChild(row);
        });
    }

    noticesList.addEventListener('click', function(e) {
        if (e.target.classList.contains('edit-btn')) {
            const noticeId = e.target.getAttribute('data-id');
            const notice = findNoticeById(noticeId);
            if (notice) {
                openEditModal(notice);
            }
        } else if (e.target.classList.contains('delete-btn')) {
            const noticeId = e.target.getAttribute('data-id');
            if (confirm('确定要删除这条公告吗？')) {
                deleteNotice(noticeId);
            }
        }
    });

    function findNoticeById(id) {
        const row = document.querySelector(`button[data-id="${id}"]`).closest('tr');
        return {
            Notice_id: row.cells[0].textContent,
            Notice_date: row.cells[1].textContent,
            Notice_title: row.cells[2].textContent,
            Notice: row.cells[3].textContent
        };
    }

    function openEditModal(notice) {
        document.getElementById('edit_notice_id').value = notice.Notice_id;
        document.getElementById('edit_notice_date').value = notice.Notice_date;
        document.getElementById('edit_notice_title').value = notice.Notice_title;
        document.getElementById('edit_notice').value = notice.Notice;
        editModal.style.display = 'block';
    }

    closeBtn.onclick = function() {
        editModal.style.display = 'none';
    }

    window.onclick = function(event) {
        if (event.target == editModal) {
            editModal.style.display = 'none';
        }
    }

    editForm.addEventListener('submit', function(e) {
        e.preventDefault();
        if (confirm('确定要修改这条公告吗？')) {
            const formData = new FormData(editForm);

            fetch('/update/notice', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    alert('公告更新成功');
                    editModal.style.display = 'none';
                    searchForm.dispatchEvent(new Event('submit'));
                } else if (response.status === 400) {
                    alert('无法修改先前的公告');
                } else if (response.status === 500) {
                    alert('服务器错误，请稍后重试');
                } else {
                    alert('更新失败，请重试');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('公告更新失败，请重试。');
            });
        }
    });

    function deleteNotice(noticeId) {
        fetch('/delete/notice', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `notice_id=${noticeId}`
        })
        .then(response => {
            if (response.ok) {
                alert('公告删除成功');
                searchForm.dispatchEvent(new Event('submit'));
            }else if (response.status === 500) {
                alert('服务器错误，请稍后重试');
            } else {
                alert('删除失败，请重试');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('公告删除失败，请重试。');
        });
    }

    // 初始加载所有公告
    searchForm.dispatchEvent(new Event('submit'));
});