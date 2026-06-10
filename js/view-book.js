document.addEventListener('DOMContentLoaded', function() {
    const bookList = document.getElementById('bookList');
    const searchForm = document.getElementById('searchForm');
    const searchCategory = document.getElementById('searchCategory');
    const searchInput = document.getElementById('searchInput');
    const editForm = document.getElementById('editForm');
    const saveChangesBtn = document.getElementById('saveChanges');
    const deleteBookBtn = document.getElementById('deleteBook');
    const recState = document.getElementById('rec_state');
    const recTypeContainer = document.getElementById('recTypeContainer');
    const detailModal = document.getElementById('detailModal');
    const editModal = document.getElementById('editModal');
    const closeBtns = document.getElementsByClassName('close');

    let currentBooks = [];

    function loadBooks(books) {
        bookList.innerHTML = '';
        currentBooks = books;
        books.forEach((book, index) => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td><img src="../${book.cover}" style="max-height: 50px;"></td>
                <td>${book.title}</td>
                <td>${book.author}</td>
                <td>${book.isbn}</td>
                <td>${book.press}</td>
                <td>${book.cur_lend_amount}</td>
                <td>${book.rec_state == "1" ? '是' : '否'}</td>
                <td>
                    <button class="btn btn-info view-details" data-index="${index}">详细信息</button>
                    <button class="btn btn-primary edit-book" data-index="${index}">编辑</button>
                </td>
            `;
            bookList.appendChild(row);
        });
    }

    searchForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const category = searchCategory.value;
        const searchTerm = searchInput.value;

        fetch('/view/book', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `category=${category}&title=${encodeURIComponent(searchTerm)}`
        })
        .then(response => response.json())
        .then(data => {
            loadBooks(data);
        })
        .catch(error => {
            console.error('搜索图书失败:', error);
            alert('搜索图书失败，请重试。');
        });
    });

    bookList.addEventListener('click', function(e) {
        if (e.target.classList.contains('view-details')) {
            const index = e.target.dataset.index;
            const book = currentBooks[index];
            const detailContent = document.getElementById('detailContent');
            detailContent.innerHTML = `
                <p><strong>书名：</strong>${book.title}</p>
                <p><strong>作者：</strong>${book.author}</p>
                <p><strong>图书类型：</strong>${book.book_type}</p>
                <p><strong>出版社：</strong>${book.press}</p>
                <p><strong>出版日期：</strong>${book.press_date}</p>
                <p><strong>ISBN：</strong>${book.isbn}</p>
                <p><strong>封面：</strong><img src="../${book.cover}" alt="封面" style="max-width: 200px;"></p>
                <p><strong>简介：</strong>${book.intro}</p>
                <p><strong>价格：</strong>${book.price}</p>
                <p><strong>数量：</strong>${book.amount}</p>
                <p><strong>可借数量：</strong>${book.lend_amount}</p>
                <p><strong>当前可借数量：</strong>${book.cur_lend_amount}</p>
                <p><strong>是否推荐：</strong>${book.rec_state == "1" ? '是' : '否'}</p>
                ${book.rec_state == "1" ? `<p><strong>推荐类型：</strong>${book.rec_type}</p>` : ''}
            `;
            detailModal.classList.add('active');
        } else if (e.target.classList.contains('edit-book')) {
            const index = e.target.dataset.index;
            const book = currentBooks[index];
            for (const [key, value] of Object.entries(book)) {
                const input = document.getElementById(key);
                if (input) {
                    if (input.type === 'checkbox') {
                        input.checked = value == "1";
                    } else if (input.type !== 'file') {
                        input.value = value;
                    }
                }
            }
            document.getElementById('yisbn').value = book.isbn;
            toggleRecType();
            editModal.classList.add('active');
        }
    });

    recState.addEventListener('change', toggleRecType);

    function closeModals() {
        detailModal.classList.remove('active');
        editModal.classList.remove('active');
    }

    function toggleRecType() {
        recTypeContainer.style.display = recState.checked ? 'block' : 'none';
    }

    saveChangesBtn.addEventListener('click', function() {
        const formData = new FormData(editForm);
        console.log(formData)
        formData.set('rec_state', recState.checked ? "1" : "0");

        fetch('/update/book', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (response.ok) {
                alert("更新成功")
                closeModals();
                searchForm.dispatchEvent(new Event('submit'));
            } else if (response.status === 409){
                alert("此isbn被占用")
            } else if (response.status === 422){
                alert("数量设置错误")
            } else{
                alert("服务器错误")
            }
            return response.text();

        })
        .catch(error => {
            console.error('更新图书信息失败:', error);
        });
    });

    deleteBookBtn.addEventListener('click', function() {
        const yisbn = document.getElementById('yisbn').value;
        if (confirm('确定要删除这本书吗？')) {
            fetch('/delete/book', {
                method: 'POST',
                headers:  {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `yisbn=${yisbn}`
            })
            .then(response => {
                if (response.ok) {
                    alert("删除成功");
                    closeModals();
                    searchForm.dispatchEvent(new Event('submit'));
                }
            })
            .catch(error => {
                console.error('删除图书失败:', error);
                alert('删除图书失败，请重试。');
            });
        }
    });

    for (let i = 0; i < closeBtns.length; i++) {
        closeBtns[i].onclick = closeModals;
    }

    window.onclick = function(event) {
        if (event.target == detailModal || event.target == editModal) {
            closeModals();
        }
    }

    // Initial load of all books
    fetch('/view/book', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'category=1&title='
    })
    .then(response => response.json())
    .then(data => {
        loadBooks(data);
    })
    .catch(error => {
        console.error('加载图书失败:', error);
    });
});
