// ==== 借书架弹窗 ====
const modal = document.getElementById("bookshelfModal");
const bookCountEl = document.getElementById("bookCount");
let bookCount = 0;

document.getElementById("myBookshelf").onclick = function() {
    if (localStorage.getItem('isLoggedIn') === 'true') {
        modal.style.display = "block";
    } else {
        window.location.href = '/login';
    }
};

document.querySelector(".close").onclick = () => modal.style.display = "none";
window.onclick = (e) => { if (e.target === modal) modal.style.display = "none"; };

// ==== 登录态 ====
document.addEventListener('DOMContentLoaded', () => {
    const authBtn = document.getElementById('authButton');
    if (localStorage.getItem('isLoggedIn') === 'true') {
        authBtn.textContent = '👤 个人中心';
        authBtn.href = '/user/library';
    }
});

// ==== 添加到借书架 ====
function addToBookshelf(title, isbn, available) {
    if (localStorage.getItem('isLoggedIn') !== 'true') {
        window.location.href = '/login';
        return;
    }
    if (available <= 0) { alert('此书已无存余'); return; }

    const container = document.getElementById('borrowedBooks');
    const el = document.createElement('div');
    el.className = 'borrowed-item';
    el.innerHTML = `
        <p><strong>${title}</strong> <span style="color:#94a3b8;font-size:12px">ISBN: ${isbn}</span></p>
        <form class="lendform">
            <label style="font-size:13px;color:#64748b">预计归还日期：</label>
            <input type="date" name="exp_return_date" required>
            <input type="hidden" name="isbn" value="${isbn}">
            <div class="borrow-actions">
                <button type="submit" class="confirm-borrow">确认借阅</button>
                <button type="button" class="remove-book">取消</button>
            </div>
        </form>
    `;
    container.appendChild(el);
    bookCount++;
    bookCountEl.textContent = `（${bookCount}）`;

    el.querySelector('.remove-book').onclick = () => removeBook(el);

    el.querySelector('.lendform').onsubmit = function(e) {
        e.preventDefault();
        if (!confirm("确认借阅《" + title + "》？")) return;
        const fd = new FormData(this);
        const token = localStorage.getItem('access_token');
        const headers = token ? { 'Authorization': 'Bearer ' + token } : {};
        const options = token ? { method: 'POST', body: fd, headers: headers } : { method: 'POST', body: fd };
        fetch('/lend/book', options)
            .then(r => {
                if (r.ok) {
                    alert("借阅成功！");
                    removeBook(el);
                    // 立即减1，不等刷新
                    updateBookCount(isbn);
                } else if (r.status === 400) alert("日期设置有误");
                else if (r.status === 403) alert("已达借阅上限或已借阅此书");
                else if (r.status === 404) alert("此书已被借完");
                else alert("借阅失败");
            })
            .catch(() => alert("网络错误"));
    };
}

function removeBook(el) {
    const container = document.getElementById('borrowedBooks');
    if (container.contains(el)) {
        container.removeChild(el);
        bookCount = Math.max(0, bookCount - 1);
        bookCountEl.textContent = `（${bookCount}）`;
    }
}

// 借书成功后立即更新页面上的可借数量
function updateBookCount(isbn) {
    document.querySelectorAll('.book-card').forEach(card => {
        const isbnEl = card.querySelector('.book-isbn');
        if (isbnEl && isbnEl.textContent.trim() === isbn) {
            const countEl = card.querySelector('.book-count');
            const btn = card.querySelector('.btn-borrow');
            let cur = parseInt(countEl.textContent);
            if (cur > 0) {
                cur--;
                countEl.textContent = cur;
                countEl.className = 'book-count ' + (cur <= 0 ? 'none' : cur <= 1 ? 'low' : 'available');
                if (cur <= 0) btn.disabled = true;
            }
        }
    });
}

// ==== 图书详情弹窗 ====
const detailModal = document.getElementById('bookDetailModal');
const detailBody = document.getElementById('bookDetailBody');

function showBookDetail(book) {
    const avail = parseInt(book.cur_lend_amount) || 0;
    const countClass = avail <= 0 ? 'none' : avail <= 1 ? 'low' : 'available';
    const coverUrl = book.cover ? '/' + book.cover : '/images/default-avatar.svg';

    detailBody.innerHTML = `
        <div style="display:flex;gap:24px">
            <div style="flex-shrink:0">
                <img src="${coverUrl}" alt="${book.title}" style="width:160px;height:230px;object-fit:cover;border-radius:8px;background:#f1f5f9">
            </div>
            <div style="flex:1;min-width:0">
                <h2 style="font-size:22px;margin-bottom:8px">${book.title}</h2>
                <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px 20px;font-size:14px;color:#475569;margin-bottom:12px">
                    <div>✍️ 作者：${book.author || '--'}</div>
                    <div>📖 类型：${book.book_type || '--'}</div>
                    <div>🏷️ 出版社：${book.press || '--'}</div>
                    <div>📅 出版日期：${book.press_date || '--'}</div>
                    <div>🔢 ISBN：${book.isbn}</div>
                    <div>💰 价格：¥${book.price || '0'}</div>
                    <div>📚 总藏书：${book.amount || '--'} 本</div>
                    <div>📊 总借出：${book.lend_amount || '0'} 本</div>
                </div>
                <div style="margin-bottom:12px">
                    <span style="font-size:28px;font-weight:700;color:${avail <= 0 ? '#ef4444' : avail <= 1 ? '#f59e0b' : '#22c55e'}">${avail}</span>
                    <span style="font-size:14px;color:#94a3b8"> 本可借</span>
                </div>
                ${book.rec_type ? '<div style="margin-bottom:8px"><span style="background:#fef3c7;color:#92400e;font-size:12px;padding:2px 10px;border-radius:4px">⭐ ' + book.rec_type + '</span></div>' : ''}
                <button class="btn-borrow" ${avail <= 0 ? 'disabled' : ''} onclick="addToBookshelf('${book.title.replace(/'/g, "\\'")}', '${book.isbn}', ${avail}); detailModal.style.display='none'">加入借书架</button>
            </div>
        </div>
        <div style="margin-top:16px;padding-top:16px;border-top:1px solid #e2e8f0">
            <div style="font-size:15px;font-weight:600;margin-bottom:8px">📝 内容简介</div>
            <div style="font-size:14px;line-height:1.8;color:#475569;white-space:pre-wrap">${book.intro || '暂无简介'}</div>
        </div>
    `;
    detailModal.style.display = 'block';
}

document.getElementById('closeDetail').onclick = () => detailModal.style.display = 'none';
detailModal.onclick = (e) => { if (e.target === detailModal) detailModal.style.display = 'none'; };

// ==== 渲染书籍列表 ====
function renderBooks(books) {
    const container = document.getElementById('bookContainer');
    container.innerHTML = '';
    document.getElementById('resultCount').textContent = `共 ${books.length} 本`;

    if (!books || books.length === 0) {
        container.innerHTML = '<div style="text-align:center;padding:60px 0;color:#94a3b8">没有找到匹配的书籍</div>';
        return;
    }

    const start = (currentPage - 1) * booksPerPage;
    const end = start + booksPerPage;
    const page = books.slice(start, end);

    page.forEach(book => {
        const card = document.createElement('div');
        card.className = 'book-card';
        const avail = parseInt(book.cur_lend_amount) || 0;
        const countClass = avail <= 0 ? 'none' : avail <= 1 ? 'low' : 'available';
        const coverUrl = book.cover ? '/' + book.cover : '/images/default-avatar.svg';

        card.innerHTML = `
            <div class="book-cover" style="cursor:pointer"><img src="${coverUrl}" alt="${book.title}" loading="lazy"></div>
            <div class="book-info">
                <div class="book-title" style="cursor:pointer;color:#3b82f6">${book.title}</div>
                <div class="book-meta">
                    <span>✍️ ${book.author || '--'}</span>
                    <span>📖 ${book.book_type || '--'}</span>
                    <span>🏷️ ${book.press || '--'}</span>
                    <span>📅 ${book.press_date || '--'}</span>
                    <span class="book-isbn" style="display:none">${book.isbn}</span>
                </div>
                <div class="book-meta" style="font-size:12px;color:#94a3b8">
                    <span>ISBN: ${book.isbn}</span>
                    <span>💰 ¥${book.price || '0'}</span>
                    <span>${book.rec_type ? '⭐ ' + book.rec_type : ''}</span>
                </div>
                <div class="book-intro">${(book.intro || '暂无简介').substring(0, 120)}${(book.intro || '').length > 120 ? '...' : ''}</div>
                <div class="book-actions">
                    <span class="book-count ${countClass}">${avail}</span>
                    <span style="font-size:12px;color:#94a3b8">本可借</span>
                    <button class="btn-borrow" ${avail <= 0 ? 'disabled' : ''}>加入借书架</button>
                    <button class="btn-detail" data-index="${allBooks.indexOf(book)}">查看详情</button>
                </div>
            </div>
        `;

        card.querySelector('.btn-borrow').onclick = () => {
            addToBookshelf(book.title, book.isbn, avail);
        };
        card.querySelector('.book-cover').onclick = () => showBookDetail(book);
        card.querySelector('.book-title').onclick = () => showBookDetail(book);
        card.querySelector('.btn-detail').onclick = () => showBookDetail(book);

        container.appendChild(card);
    });
}

// 覆盖 change-page.js 中的 displayBooks
function displayBooks() {
    renderBooks(allBooks);
    setupPagination();
}

// ==== 实时刷新（每30秒）====
let refreshTimer = null;

function startAutoRefresh() {
    if (refreshTimer) clearInterval(refreshTimer);
    refreshTimer = setInterval(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const selsearch = urlParams.get('selsearch');
        const inpsearch = urlParams.get('inpsearch');
        if (selsearch && inpsearch) {
            fetchBooksWithParams(selsearch, inpsearch);
        } else if (Object.values(selectedValues).some(v => v)) {
            // 有筛选条件时通过筛选API刷新
            const fd = new FormData();
            for (let key in selectedValues) {
                if (selectedValues[key]) fd.append(key, selectedValues[key]);
            }
            fetch('/class/search', { method: 'POST', body: fd })
                .then(r => r.json())
                .then(data => { if (Array.isArray(data)) { allBooks = data; displayBooks(); } })
                .catch(() => {});
        } else {
            fetchBooksSilent();
        }
    }, 30000);
}

function fetchBooksSilent() {
    fetch('/search/book', { method: 'POST', headers: { 'Accept': 'application/json' } })
        .then(r => r.json())
        .then(data => { if (Array.isArray(data)) { allBooks = data; displayBooks(); }})
        .catch(() => {});
}

function fetchBooksWithParams(sel, inp) {
    fetch(`/search/book?selsearch=${encodeURIComponent(sel)}&inpsearch=${encodeURIComponent(inp)}`, {
        headers: { 'Accept': 'application/json' }
    })
    .then(r => r.json())
    .then(data => { if (Array.isArray(data)) { allBooks = data; displayBooks(); } })
    .catch(() => {});
}

// 覆盖 change-page.js 中的 fetchBooks 和 searchBooks
const origFetchBooks = window.fetchBooks || function(){};
const origSearchBooks = window.searchBooks || function(){};

window.fetchBooks = function() {
    origFetchBooks();
    startAutoRefresh();
};

window.searchBooks = function(sel, inp) {
    fetchBooksWithParams(sel, inp);
    startAutoRefresh();
};

// 在 change-page.js 加载完成后启动
setTimeout(startAutoRefresh, 1000);
