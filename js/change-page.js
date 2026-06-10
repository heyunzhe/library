let currentPage = 1;
const booksPerPage = 5;
let allBooks = [];

let selectedValues = {
    value: null,
    value1: null,
    value2: null,
    value3: null,
    value4: null
};

document.addEventListener('DOMContentLoaded', function() {
    // 分类筛选链接
    document.querySelectorAll('.filter-links a').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const section = this.closest('.filter-section');
            if (!section) return;
            const catLink = section.querySelector('.filter-links');
            if (!catLink) return;
            const category = catLink.dataset.category;
            if (!category) return;

            if (selectedValues[category] === this.textContent) {
                selectedValues[category] = null;
            } else {
                selectedValues[category] = this.textContent;
            }

            updateSelectDisplay();

            const fd = new FormData();
            for (let key in selectedValues) {
                if (selectedValues[key]) fd.append(key, selectedValues[key]);
            }

            fetch('/class/search', { method: 'POST', body: fd })
                .then(r => {
                    if (!r.ok) throw new Error('HTTP ' + r.status);
                    return r.json();
                })
                .then(data => {
                    allBooks = Array.isArray(data) ? data : [];
                    currentPage = 1;
                    displayBooks();
                    setupPagination();
                })
                .catch(() => {
                    allBooks = [];
                    displayBooks();
                    setupPagination();
                });
        });
    });

    // URL搜索参数
    const params = new URLSearchParams(window.location.search);
    const selsearch = params.get('selsearch');
    const inpsearch = params.get('inpsearch');
    if (selsearch && inpsearch) {
        searchBooks(selsearch, inpsearch);
    } else {
        fetchBooks();
    }

    // 搜索表单
    document.getElementById('searchForm').addEventListener('submit', function(e) {
        e.preventDefault();
        const fd = new FormData(this);
        searchBooks(fd.get('selsearch'), fd.get('inpsearch'));
    });
});

function updateSelectDisplay() {
    const el = document.getElementById('select');
    el.textContent = Object.values(selectedValues).filter(Boolean).join(' ');
}

function fetchBooks() {
    fetch('/search/book', { method: 'POST', headers: { 'Accept': 'application/json' } })
        .then(r => r.json())
        .then(data => {
            allBooks = Array.isArray(data) ? data : [];
            displayBooks();
            setupPagination();
        })
        .catch(() => { allBooks = []; displayBooks(); setupPagination(); });
}

function searchBooks(selsearch, inpsearch) {
    fetch(`/search/book?selsearch=${encodeURIComponent(selsearch)}&inpsearch=${encodeURIComponent(inpsearch)}`, {
        headers: { 'Accept': 'application/json' }
    })
    .then(r => r.json())
    .then(data => {
        allBooks = Array.isArray(data) ? data : [];
        currentPage = 1;
        displayBooks();
        setupPagination();
    })
    .catch(() => { allBooks = []; displayBooks(); setupPagination(); });
}

function displayBooks() {
    // 由 lend-book-list.js 中的 renderBooks 处理
    if (typeof renderBooks === 'function') {
        renderBooks(allBooks);
    } else {
        // 降级：如果没有 renderBooks，用基础渲染
        const container = document.getElementById('bookContainer');
        if (!container) return;
        container.innerHTML = '';
        if (!allBooks || allBooks.length === 0) {
            container.innerHTML = '<div style="text-align:center;padding:40px;color:#94a3b8">暂无书籍</div>';
            return;
        }
        const start = (currentPage - 1) * booksPerPage;
        const end = start + booksPerPage;
        allBooks.slice(start, end).forEach(book => {
            const div = document.createElement('div');
            div.textContent = book.title;
            container.appendChild(div);
        });
    }
}

function setupPagination() {
    const total = Math.ceil(allBooks.length / booksPerPage);
    const el = document.getElementById('pagination');
    if (!el) return;
    el.innerHTML = '';

    const prev = document.createElement('button');
    prev.textContent = '‹ 上一页';
    prev.disabled = currentPage <= 1;
    prev.onclick = () => { if (currentPage > 1) { currentPage--; displayBooks(); setupPagination(); } };
    el.appendChild(prev);

    for (let i = 1; i <= total; i++) {
        const btn = document.createElement('button');
        btn.textContent = i;
        if (i === currentPage) btn.className = 'active';
        btn.onclick = () => { currentPage = i; displayBooks(); setupPagination(); };
        el.appendChild(btn);
    }

    const next = document.createElement('button');
    next.textContent = '下一页 ›';
    next.disabled = currentPage >= total;
    next.onclick = () => { if (currentPage < total) { currentPage++; displayBooks(); setupPagination(); } };
    el.appendChild(next);
}
