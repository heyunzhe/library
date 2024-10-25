let currentPage = 1;
const booksPerPage = 5;
let allBooks = [];


document.addEventListener('DOMContentLoaded', function() {
    // 为所有分类链接添加点击事件监听器
    document.querySelectorAll('.box7 a, .box8 a, .box9 a, .box10 a').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault(); // 阻止默认的链接行为
            
            let category, value;
            if (this.closest('.box7')) {
                category = '作者';
            } else if (this.closest('.box8')) {
                category = '出版社';
            }else if (this.closest('.box9')) {
                category = '类型';
            }else if (this.closest('.box10')) {
                category = '出版日期';
            }
            value = this.textContent;

            const select = document.getElementById('select');
            select.innerHTML = `${value}`

            // 创建 FormData 对象
            const formData = new FormData();
            formData.append('category', category);
            formData.append('value', value);


            // 发送请求到服务器
            fetch('/lend/book', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                allBooks = data;
                currentPage = 1; // 重置到第一页
                displayBooks();
                setupPagination();
            })
            .catch(error => console.error('Error:', error));
        });
    });

    // 初始加载所有书籍
    fetchBooks();
});

function fetchBooks() {
    fetch('/lend/book', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => response.json())
    .then(data => {
        allBooks = data;
        displayBooks();
        setupPagination();
    })
    .catch(error => console.error('Error:', error));
}

const searchForm = document.querySelector('form[action="/search/book"]');
searchForm.addEventListener('submit', function(e) {
    e.preventDefault();
    const formData = new FormData(this);
    searchBooks(formData);
});

function searchBooks(formData) {
    fetch('/search/book', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        allBooks = data;
        currentPage = 1; // 重置到第一页
        displayBooks();
        setupPagination();
    })
    .catch(error => console.error('Error:', error));
}


function displayBooks() {
    const bookContainer = document.querySelector('.box13');
    bookContainer.innerHTML = '';

    const start = (currentPage - 1) * booksPerPage;
    const end = start + booksPerPage;
    const booksToDisplay = allBooks.slice(start, end);

    booksToDisplay.forEach((book, index) => {
        const bookElement = document.createElement('div');
        bookElement.className = `box13-${index + 1}`;
        bookElement.innerHTML = `
            <img src="../${book.cover}" alt="${book.title}" class="book">
            <div class="box13-${index + 1}-1">
                <h3><a href="">${book.title}</a></h3>
                <ul>
                    <li>作者：${book.author}</li>
                    <label for="">ISBN:</label><li class="ISBN">${book.isbn}</li>
                </ul>
                <ul>
                    <li>出版年份：${book.press_date}</li>
                    <li>出版社：${book.press}</li>
                </ul>
                <ul>
                    <li>价格：${book.price}</li>
                    <li>可借数：${book.cur_lend_amount}</li>
                </ul>
                <ul>
                    <li>简介：</li>
                </ul>
                <div class="intro_style">
                    <p class="intro">${book.intro}</p>
                </div>
                <div class="lend_book">
                    <a href="#" id="lend_book">加入借书架</a>
                </div>
            </div>
        `;
        bookContainer.appendChild(bookElement);
    });
}

function setupPagination() {
    const totalPages = Math.ceil(allBooks.length / booksPerPage);
    const paginationContainer = document.getElementById('pagination');
    paginationContainer.innerHTML = '';

    // Previous button
    const prevButton = document.createElement('button');
    prevButton.textContent = '上一页';
    prevButton.addEventListener('click', () => {
        if (currentPage > 1) {
            currentPage--;
            displayBooks();
        }
    });
    paginationContainer.appendChild(prevButton);

    // Page numbers
    for (let i = 1; i <= totalPages; i++) {
        const pageButton = document.createElement('button');
        pageButton.textContent = i;
        pageButton.addEventListener('click', () => {
            currentPage = i;
            displayBooks();
        });
        paginationContainer.appendChild(pageButton);
    }

    // Next button
    const nextButton = document.createElement('button');
    nextButton.textContent = '下一页';
    nextButton.addEventListener('click', () => {
        if (currentPage < totalPages) {
            currentPage++;
            displayBooks();
        }
    });
    paginationContainer.appendChild(nextButton);
}

document.addEventListener('DOMContentLoaded', fetchBooks);


