let currentPage = 1;
const booksPerPage = 5;
let allBooks = [];

let selectedValues = {
    value: null,
    value1: null,
    value2: null,
    value3: null
};

document.addEventListener('DOMContentLoaded', function() {
    // 为所有分类链接添加点击事件监听器
    document.querySelectorAll('.box7 a, .box8 a, .box9 a, .box10 a').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            
            let category;
            if (this.closest('.box7')) {
                category = 'value';
            } else if (this.closest('.box8')) {
                category = 'value1';
            } else if (this.closest('.box9')) {
                category = 'value2';
            } else if (this.closest('.box10')) {
                category = 'value3';
            }
            
            selectedValues[category] = this.textContent;

            updateSelectDisplay();

            const formData = new FormData();
            for (let key in selectedValues) {
                if (selectedValues[key]) {
                    formData.append(key, selectedValues[key]);
                }
            }

            fetch('/class/search', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (!Array.isArray(data)) {
                    throw new Error('Server did not return an array');
                }
                console.log('Received data:', data);
                allBooks = data;
                currentPage = 1;
                displayBooks();
                setupPagination();
            })
            .catch(error => {
                console.error('Error:', error);
                allBooks = [];
                displayBooks();
                setupPagination();
            });
        });
    });
    
    // 检查URL参数并执行搜索（如果有参数）
    const urlParams = new URLSearchParams(window.location.search);
    const selsearch = urlParams.get('selsearch');
    const inpsearch = urlParams.get('inpsearch');
    if (selsearch && inpsearch) {
        searchBooks(selsearch, inpsearch);
    } else {
        // 如果没有搜索参数，则加载所有书籍
        fetchBooks();
    }

    // 添加搜索表单的事件监听器
    const searchForm = document.querySelector('form[action="/search/book"]');
    searchForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const formData = new FormData(this);
        const selsearch = formData.get('selsearch');
        const inpsearch = formData.get('inpsearch');
        searchBooks(selsearch, inpsearch);
    });
});

function updateSelectDisplay() {
    const select = document.getElementById('select');
    select.innerHTML = Object.values(selectedValues).filter(Boolean).join(' ');
}

function fetchBooks() {
    fetch('/search/book', {
        headers: {
            'Accept': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log('Fetched all books:', data);
        allBooks = data;
        displayBooks();
        setupPagination();
    })
    .catch(error => console.error('Error:', error));
}

function searchBooks(selsearch, inpsearch) {
    const url = `/search/book?selsearch=${encodeURIComponent(selsearch)}&inpsearch=${encodeURIComponent(inpsearch)}`;
    fetch(url, {
        headers: {
            'Accept': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log('Search results:', data);
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

    if (!allBooks || allBooks.length === 0) {
        bookContainer.innerHTML = '<p>没有找到匹配的书籍。</p>';
        return;
    }

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